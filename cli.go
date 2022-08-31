package gh2changelog

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"strings"
)

const cmdName = "gh2changelog"

// Run the gh2changelog
func Run(ctx context.Context, argv []string, outStream, errStream io.Writer) error {
	log.SetOutput(errStream)
	fs := flag.NewFlagSet(
		fmt.Sprintf("%s (v%s rev:%s)", cmdName, version, revision), flag.ContinueOnError)
	fs.SetOutput(errStream)
	var (
		ver = fs.Bool("version", false, "display version")

		git  = fs.String("git", "git", "git path")
		repo = fs.String("repo", ".", "local repository path")

		verbose = fs.Bool("verbose", false, "verbose")

		tag        = fs.String("tag", "", "specify existing tag")
		next       = fs.String("next", "", "tag to be released next")
		unreleased = fs.Bool("unreleased", false, "output unreleased")
		latest     = fs.Bool("latest", false, "get latest changelog section")
		limit      = fs.Int("limit", 0, "outputs the specified number of most recent changelogs")
		all        = fs.Bool("all", false, "outputs all changelogs")

		alone = fs.Bool("alone", false, "only outputs the specified changelog without merging with CHANGELOG.md.")
		write = fs.Bool("w", false, "write result to CHANGELOG.md")
	)

	if err := fs.Parse(argv); err != nil {
		return err
	}
	if *ver {
		return printVersion(outStream)
	}

	opts := []Option{GitPath(*git), RepoPath(*repo)}
	if *verbose {
		opts = append(opts, SetOutputs(outStream, errStream))
	}
	gch, err := New(ctx, opts...)
	if err != nil {
		return err
	}
	if *all {
		if *limit != 0 {
			log.Println("Both the limit and all options are specified, but the all option takes precedence.")
		}
		*limit = -1
	}

	if *limit != 0 {
		logs, _, err := gch.Changelogs(ctx, *limit)
		if err != nil {
			return err
		}
		if *unreleased {
			log, _, err := gch.Unreleased(ctx)
			if err != nil {
				return err
			}
			logs = append([]string{log}, logs...)
		}
		return gch.output(outStream, strings.Join(logs, "\n"), Trunc, *write)
	}

	var outs []string

	if *next != "" {
		if *unreleased {
			log.Println("Both unreleased and next options are specified, but next takes precedence.")
		}
		out, _, err := gch.Draft(ctx, *next)
		if err != nil {
			return err
		}
		outs = append(outs, out)
	} else if *unreleased {
		out, _, err := gch.Unreleased(ctx)
		if err != nil {
			return err
		}
		outs = append(outs, out)
	}

	if *latest || (len(outs) < 1 && *tag == "") {
		out, _, err := gch.Latest(ctx)
		if err != nil {
			return err
		}
		outs = append(outs, out)
	}

	if *tag != "" {
		out, _, err := gch.Changelog(ctx, *tag)
		if err != nil {
			return err
		}
		outs = append(outs, out)
	}

	out := strings.Join(outs, "\n")

	if *alone {
		if *write {
			log.Println("Both alone and w options are specified, but alone takes precedence.")
		}
		_, err = fmt.Fprint(outStream, out)
		return err
	}
	return gch.output(outStream, out, 0, *write)
}

func (gch *GH2Changelog) output(outStream io.Writer, out string, mode int, write bool) error {
	if !write {
		mode |= DryRun
	}
	out, err := gch.Update(out, mode)
	if err != nil {
		return err
	}
	if !write {
		_, err = fmt.Fprint(outStream, out)
		return err
	}
	return nil
}

func printVersion(out io.Writer) error {
	_, err := fmt.Fprintf(out, "%s v%s (rev:%s)\n", cmdName, version, revision)
	return err
}
