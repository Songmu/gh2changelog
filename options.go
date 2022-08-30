package gh2changelog

import (
	"io"

	"github.com/google/go-github/v45/github"
)

func GitPath(p string) Option {
	return func(gch *GH2Changelog) {
		gch.gitPath = p
	}
}

func RepoPath(p string) Option {
	return func(gch *GH2Changelog) {
		gch.repoPath = p
		gch.c.dir = p
	}
}

func SetOutputs(outStream, errStream io.Writer) Option {
	return func(gch *GH2Changelog) {
		gch.c.outStream = outStream
		gch.c.errStream = errStream
	}
}

func GitHubClient(cli *github.Client) Option {
	return func(gch *GH2Changelog) {
		gch.gh = cli
	}
}
