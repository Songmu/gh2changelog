package gh2changelog

import (
	"io"

	"github.com/google/go-github/v49/github"
)

// GitPath sets a git executable path
func GitPath(p string) Option {
	return func(gch *GH2Changelog) {
		gch.gitPath = p
	}
}

// RepoPath sets a repository path
func RepoPath(p string) Option {
	return func(gch *GH2Changelog) {
		gch.repoPath = p
	}
}

// SetOutputs sets a stdout and a stderr
func SetOutputs(outStream, errStream io.Writer) Option {
	return func(gch *GH2Changelog) {
		gch.outStream = outStream
		gch.errStream = errStream
	}
}

// GitHubClient sets a github.Client
func GitHubClient(cli *github.Client) Option {
	return func(gch *GH2Changelog) {
		gch.gen = cli.Repositories
	}
}
