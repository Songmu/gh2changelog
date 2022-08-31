gh2changelog
=======

[![Test Status](https://github.com/Songmu/gh2changelog/workflows/test/badge.svg?branch=main)][actions]
[![Coverage Status](https://codecov.io/gh/Songmu/gh2changelog/branch/main/graph/badge.svg)][codecov]
[![MIT License](https://img.shields.io/github/license/Songmu/gh2changelog)][license]
[![PkgGoDev](https://pkg.go.dev/badge/github.com/Songmu/gh2changelog)][PkgGoDev]

[actions]: https://github.com/Songmu/gh2changelog/actions?workflow=test
[codecov]: https://codecov.io/gh/Songmu/gh2changelog
[license]: https://github.com/Songmu/gh2changelog/blob/main/LICENSE
[PkgGoDev]: https://pkg.go.dev/github.com/Songmu/gh2changelog

gh2changelog generates keep a changelog like output from GitHub Releases

## Synopsis

```console
% gh2changelog
# Changelog

## [v0.0.14](https://github.com/Songmu/tagpr/compare/v0.0.13...v0.0.14) - 2022-08-28
- fix version file detection in releasing by @Songmu in https://github.com/Songmu/tagpr/pull/70
...
```

## Description

The "gh2changelog" outputs a changelog of the "Keep a changelog" like format and writes it to CHANGELOG.md.
To generate a changelog, use "generate-notes" in the GitHub REST API's Releases.

- https://keepachangelog.com/
- https://docs.github.com/ja/rest/releases/releases#generate-release-notes-content-for-a-release

## Options

```
  -all          outputs all changelogs
  -alone        only outputs the specified changelog without merging with CHANGELOG.md.
  -git string   git path (default "git")
  -latest       get latest changelog section
  -limit int    outputs the specified number of most recent changelogs
  -next string  tag to be released next
  -repo string  local repository path (default ".")
  -tag string   specify existing tag
  -unreleased   output unreleased
  -verbose      verbose
  -w            write result to CHANGELOG.md
```

## GITHUB Token

When github's api token is required, it is used in the following order of priority.

- command line option `--token`
- enviroment variable `GITHUB_TOKEN`
- `git config github.token`

## Installation

```console
# Install the latest version. (Install it into ./bin/ by default).
% curl -sfL https://raw.githubusercontent.com/Songmu/gh2changelog/main/install.sh | sh -s

# Specify installation directory ($(go env GOPATH)/bin/) and version.
% curl -sfL https://raw.githubusercontent.com/Songmu/gh2changelog/main/install.sh | sh -s -- -b $(go env GOPATH)/bin [vX.Y.Z]

# In alpine linux (as it does not come with curl by default)
% wget -O - -q https://raw.githubusercontent.com/Songmu/gh2changelog/main/install.sh | sh -s [vX.Y.Z]

# go install
% go install github.com/Songmu/gh2changelog/cmd/gh2changelog@latest
```

Built binaries are available on gihub releases.
<https://github.com/Songmu/gh2changelog/releases>

## Author

[Songmu](https://github.com/Songmu)
