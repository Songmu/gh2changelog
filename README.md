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

gh2changelog short description

## Synopsis

```go
// simple usage here
```

## Description

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

## Author

[Songmu](https://github.com/Songmu)
