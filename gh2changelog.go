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

// GH2Changelog is to output changelogs
type GH2Changelog struct {
	gitPath                 string
	repoPath                string
	owner, repo, remoteName string
	c                       *commander
	gh                      *github.Client
}

// Options is for functional option
type Option func(*GH2Changelog)

// New returns new GH2Changelog
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

// Draft gets draft changelog
func (gch *GH2Changelog) Draft(ctx context.Context, nextTag string) (string, string, error) {
	vers := (&gitsemvers.Semvers{GitPath: gch.gitPath}).VersionStrings()
	var previousTag *string
	if len(vers) > 0 {
		previousTag = &vers[0]
	}

	releaseBranch, err := gch.defaultBranch()
	if err != nil {
		return "", "", err
	}
	releases, _, err := gch.gh.Repositories.GenerateReleaseNotes(
		ctx, gch.owner, gch.repo, &github.GenerateNotesOptions{
			TagName:         nextTag,
			PreviousTagName: previousTag,
			TargetCommitish: &releaseBranch,
		})
	if err != nil {
		return "", "", err
	}
	return convertKeepAChangelogFormat(releases.Body, time.Now()), releases.Body, nil
}

// Unreleased gets unreleased changelog
func (gch *GH2Changelog) Unreleased(ctx context.Context) (string, string, error) {
	const tentativeTag = "v999999.999.999"
	body, orig, err := gch.Draft(ctx, tentativeTag)
	if err != nil {
		return "", "", err
	}
	bodies := strings.Split(body, "\n")
	for i, b := range bodies {
		if strings.HasPrefix(b, `## [`+tentativeTag+`](http`) {
			bodies[i] = "## [Unreleased]"
		}
	}
	return strings.Join(bodies, "\n") + "\n", orig, nil
}

// Latest gets latest changelog
func (gch *GH2Changelog) Latest(ctx context.Context) (string, string, error) {
	vers := (&gitsemvers.Semvers{GitPath: gch.gitPath}).VersionStrings()
	if len(vers) == 0 {
		return "", "", errors.New("no change log found. Never released yet")
	}
	ver := vers[0]
	date, _, err := gch.c.GitE("log", "-1", "--format=%ai", "--date=iso", ver)
	if err != nil {
		return "", "", err
	}
	d, _ := time.Parse("2006-01-02 15:04:05 -0700", date)
	releases, _, err := gch.gh.Repositories.GenerateReleaseNotes(
		ctx, gch.owner, gch.repo, &github.GenerateNotesOptions{
			TagName: ver,
		})
	if err != nil {
		return "", "", err
	}
	return strings.TrimSpace(convertKeepAChangelogFormat(releases.Body, d)) + "\n", releases.Body, nil
}

// Changelogs gets changelogs
func (gch *GH2Changelog) Changelogs(ctx context.Context, limit int) ([]string, []string, error) {
	vers := (&gitsemvers.Semvers{GitPath: gch.gitPath}).VersionStrings()
	// logs := []string{"# Changelog\n"}
	var (
		logs     []string
		origLogs []string
	)
	for i, ver := range vers {
		if limit != -1 && i > limit {
			break
		}
		date, _, err := gch.c.GitE("log", "-1", "--format=%ai", "--date=iso", ver)
		if err != nil {
			return nil, nil, err
		}
		d, _ := time.Parse("2006-01-02 15:04:05 -0700", date)
		releases, _, err := gch.gh.Repositories.GenerateReleaseNotes(
			ctx, gch.owner, gch.repo, &github.GenerateNotesOptions{
				TagName: ver,
			})
		if err != nil {
			return nil, nil, err
		}
		origLogs = append(origLogs, releases.Body)
		logs = append(logs, strings.TrimSpace(convertKeepAChangelogFormat(releases.Body, d))+"\n")
	}
	return logs, origLogs, nil
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

var headBranchReg = regexp.MustCompile(`(?m)^\s*HEAD branch: (.*)$`)

func (gch *GH2Changelog) defaultBranch() (string, error) {
	// `git symbolic-ref refs/remotes/origin/HEAD` sometimes doesn't work
	// So use `git remote show origin` for detecting default branch
	show, _, err := gch.c.GitE("remote", "show", gch.remoteName)
	if err != nil {
		return "", fmt.Errorf("failed to detect defaut branch: %w", err)
	}
	m := headBranchReg.FindStringSubmatch(show)
	if len(m) < 2 {
		return "", fmt.Errorf("failed to detect default branch from remote: %s", gch.remoteName)
	}
	return m[1], nil
}
