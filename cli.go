package gh2changelog

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	cmdName     = "gh2changelog"
	changelogMd = "CHANGELOG.md"
	Heading     = "# Changelog\n"
)

// Run the gh2changelog
func Run(ctx context.Context, argv []string, outStream, errStream io.Writer) error {
	log.SetOutput(errStream)
	fs := flag.NewFlagSet(
		fmt.Sprintf("%s (v%s rev:%s)", cmdName, version, revision), flag.ContinueOnError)
	fs.SetOutput(errStream)
	var (
		ver     = fs.Bool("version", false, "display version")
		verbose = fs.Bool("verbose", false, "verbose")
		git     = fs.String("git", "git", "git path")
		repo    = fs.String("repo", ".", "local repository path")

		tag        = fs.String("tag", "", "specify the tag")
		next       = fs.String("next", "", "tag to be released next")
		unreleased = fs.Bool("unreleased", false, "output unreleased")
		latest     = fs.Bool("latest", false, "get latest changelog section")
		limit      = fs.Int("limit", 0, "limit")
		all        = fs.Bool("all", false, "output all changelogs")
		write      = fs.Bool("w", false, "write result to file")
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

	chMdPath := filepath.Join(*repo, changelogMd)

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

		out := strings.Join(append([]string{Heading}, logs...), "\n")

		if !*write {
			_, err := fmt.Fprint(outStream, out)
			return err
		}
		return os.WriteFile(chMdPath, []byte(out), 0666)
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

	b, err := os.ReadFile(chMdPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	orig := string(b)
	if orig == "" {
		out = Heading + "\n" + out
	} else {
		out = InsertNewChangelog(orig, out)
	}
	if *write {
		return os.WriteFile(chMdPath, []byte(out), 0666)
	}
	_, err = fmt.Fprint(outStream, out)
	return err
}

func printVersion(out io.Writer) error {
	_, err := fmt.Fprintf(out, "%s v%s (rev:%s)\n", cmdName, version, revision)
	return err
}
