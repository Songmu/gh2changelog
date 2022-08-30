package gh2changelog

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/Songmu/gitsemvers"
	"github.com/google/go-github/v45/github"
)

type GH2Changelog struct {
	gitPath                 string
	repoPath                string
	owner, repo, remoteName string
	c                       *commander
	gh                      *github.Client
}

type Option func(*GH2Changelog)

func New(ctx context.Context, opts ...Option) (*GH2Changelog, error) {
	gch := &GH2Changelog{
		gitPath:  "git",
		repoPath: ".",
		c:        &commander{dir: ".", outStream: io.Discard, errStream: io.Discard},
	}
	for _, opt := range opts {
		opt(gch)
	}

	var err error
	gch.remoteName, err = gch.detectRemote()
	if err != nil {
		return nil, err
	}
	remoteURL, _, err := gch.c.GitE("config", "remote."+gch.remoteName+".url")
	if err != nil {
		return nil, err
	}
	u, err := parseGitURL(remoteURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse remote")
	}
	m := strings.Split(strings.TrimPrefix(u.Path, "/"), "/")
	if len(m) < 2 {
		return nil, fmt.Errorf("failed to detect owner and repo from remote URL")
	}
	gch.owner = m[0]
	repo := m[1]
	if u.Scheme == "ssh" || u.Scheme == "git" {
		repo = strings.TrimSuffix(repo, ".git")
	}
	gch.repo = repo

	if gch.gh == nil {
		cli, err := ghClient(ctx, "", u.Hostname())
		if err != nil {
			return nil, err
		}
		gch.gh = cli
	}
	return gch, nil
}

func (gch *GH2Changelog) getLogs(ctx context.Context, limit int) ([]string, error) {
	vers := (&gitsemvers.Semvers{GitPath: "git"}).VersionStrings()
	// logs := []string{"# Changelog\n"}
	var logs []string
	for i, ver := range vers {
		if limit != -1 && i > limit {
			break
		}
		date, _, _ := gch.c.GitE("log", "-1", "--format=%ai", "--date=iso", ver)
		d, _ := time.Parse("2006-01-02 15:04:05 -0700", date)
		releases, _, _ := gch.gh.Repositories.GenerateReleaseNotes(
			ctx, gch.owner, gch.repo, &github.GenerateNotesOptions{
				TagName: ver,
			})
		logs = append(logs, strings.TrimSpace(convertKeepAChangelogFormat(releases.Body, d))+"\n")
	}
	return logs, nil
}

func (gch *GH2Changelog) detectRemote() (string, error) {
	remotesStr, _, err := gch.c.GitE("remote")
	if err != nil {
		return "", fmt.Errorf("failed to detect remote: %s", err)
	}
	remotes := strings.Fields(remotesStr)
	if len(remotes) < 1 {
		return "", errors.New("failed to detect remote")
	}
	for _, r := range remotes {
		if r == "origin" {
			return r, nil
		}
	}
	// the last output is the first added remote
	return remotes[len(remotes)-1], nil
}

var (
	hasSchemeReg  = regexp.MustCompile("^[^:]+://")
	scpLikeURLReg = regexp.MustCompile("^([^@]+@)?([^:]+):(/?.+)$")
)

func parseGitURL(u string) (*url.URL, error) {
	if !hasSchemeReg.MatchString(u) {
		if m := scpLikeURLReg.FindStringSubmatch(u); len(m) == 4 {
			u = fmt.Sprintf("ssh://%s%s/%s", m[1], m[2], strings.TrimPrefix(m[3], "/"))
		}
	}
	return url.Parse(u)
}
