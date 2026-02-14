package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/Songmu/tagpr/gh2changelog"
)

const cmdName = "gh2changelog"

func run(ctx context.Context, argv []string, outStream, errStream io.Writer) error {
	log.SetOutput(errStream)
	flagSet := flag.NewFlagSet(
		fmt.Sprintf("%s (v%s rev:%s)", cmdName, version, revision), flag.ContinueOnError)
	flagSet.SetOutput(errStream)
	
	// Define flags
	showVersion := flagSet.Bool("version", false, "display version")
	gitCmd := flagSet.String("git", "git", "git path")
	repoDir := flagSet.String("repo", ".", "local repository path")
	verboseMode := flagSet.Bool("verbose", false, "verbose")
	tagName := flagSet.String("tag", "", "specify existing tag")
	nextTag := flagSet.String("next", "", "tag to be released next")
	showUnreleased := flagSet.Bool("unreleased", false, "output unreleased")
	getLatest := flagSet.Bool("latest", false, "get latest changelog section")
	limitCount := flagSet.Int("limit", 0, "outputs the specified number of most recent changelogs")
	showAll := flagSet.Bool("all", false, "outputs all changelogs")
	aloneMode := flagSet.Bool("alone", false, "only outputs the specified changelog without merging with CHANGELOG.md.")
	writeFile := flagSet.Bool("w", false, "write result to CHANGELOG.md")

	if err := flagSet.Parse(argv); err != nil {
		return err
	}
	
	if *showVersion {
		return printVersion(outStream)
	}

	// Build options
	options := []gh2changelog.Option{gh2changelog.GitPath(*gitCmd), gh2changelog.RepoPath(*repoDir)}
	if *verboseMode {
		options = append(options, gh2changelog.SetOutputs(outStream, errStream))
	}
	
	gch, err := gh2changelog.New(ctx, options...)
	if err != nil {
		return err
	}
	
	if *showAll {
		if *limitCount != 0 {
			log.Println("Both the limit and all options are specified, but the all option takes precedence.")
		}
		*limitCount = -1
	}

	// Handle limit option
	if *limitCount != 0 {
		logs, _, err := gch.Changelogs(ctx, *limitCount)
		if err != nil {
			return err
		}
		if *showUnreleased {
			unreleasedLog, _, err := gch.Unreleased(ctx)
			if err != nil {
				return err
			}
			logs = append([]string{unreleasedLog}, logs...)
		}
		return output(gch, outStream, strings.Join(logs, "\n"), gh2changelog.Trunc, *writeFile)
	}

	// Build output sections
	var sections []string

	if *nextTag != "" {
		if *showUnreleased {
			log.Println("Both unreleased and next options are specified, but next takes precedence.")
		}
		draftLog, _, err := gch.Draft(ctx, *nextTag, "", time.Now())
		if err != nil {
			return err
		}
		sections = append(sections, draftLog)
	} else if *showUnreleased {
		unreleasedLog, _, err := gch.Unreleased(ctx)
		if err != nil {
			return err
		}
		sections = append(sections, unreleasedLog)
	}

	if *getLatest || (len(sections) < 1 && *tagName == "") {
		latestLog, _, err := gch.Latest(ctx)
		if err != nil {
			return err
		}
		sections = append(sections, latestLog)
	}

	if *tagName != "" {
		tagLog, _, err := gch.Changelog(ctx, *tagName)
		if err != nil {
			return err
		}
		sections = append(sections, tagLog)
	}

	combinedOutput := strings.Join(sections, "\n")

	if *aloneMode {
		if *writeFile {
			log.Println("Both alone and w options are specified, but alone takes precedence.")
		}
		_, err = fmt.Fprint(outStream, combinedOutput)
		return err
	}
	
	return output(gch, outStream, combinedOutput, 0, *writeFile)
}

func output(gch *gh2changelog.GH2Changelog, outStream io.Writer, content string, mode int, shouldWrite bool) error {
	if !shouldWrite {
		mode |= gh2changelog.DryRun
	}
	result, err := gch.Update(content, mode)
	if err != nil {
		return err
	}
	if !shouldWrite {
		_, err = fmt.Fprint(outStream, result)
		return err
	}
	return nil
}

func printVersion(out io.Writer) error {
	_, err := fmt.Fprintf(out, "%s v%s (rev:%s)\n", cmdName, version, revision)
	return err
}
