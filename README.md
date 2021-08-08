# lukasmalkmus/horcrux

> A security question based secret sharing utility.

[![Go Workflow][go_workflow_badge]][go_workflow]
[![Coverage Status][coverage_badge]][coverage]
[![Go Report][report_badge]][report]
[![Go Reference][gopkg_badge]][gopkg]
[![Latest Release][release_badge]][release]
[![License][license_badge]][license]

---

## Table of Contents

1. [Introduction](#introduction)
1. [Usage](#usage)
1. [Contributing](#contributing)
1. [License](#license)

## Introduction

_horcrux_ is a security question based secret sharing utility. The idea and
package code is inspired and mostly taken from the abandoned [horcrux][1]
package by [Coda Hale][2].

_horcrux_ splits a secret into multiple fragments and associates every fragment
with a security question. The answer to that question is used to encrypt the
fragment using ChaCha20Poly1305. Only a given number of fragments is needed to
fully restore the original secret.

  [1]: https://github.com/codahale/horcrux
  [2]: https://github.com/codahale

## Installation

### Download and install the pre-compiled binary manually

Binary releases are available on [GitHub Releases][3].

  [3]: https://github.com/lukasmalkmus/horcrux/releases/latest

### Install using [Homebrew][4]

```shell
brew tap lukasmalkmus/tap
brew install horcrux
```

  [4]: https://brew.sh

To update:

```shell
brew upgrade horcrux
```

### Install using `go get`

```shell
go get -u github.com/lukasmalkmus/horcrux/cmd/horcrux
```

### Install from source

```shell
git clone https://github.com/lukasmalkmus/horcrux.git
cd horcrux
make install # Build and install binary into $GOPATH
```

### Run the Docker image

Docker images are available on the [GitHub Container Registry][5].

```shell
docker pull ghcr.io/lukasmalkmus/horcrux
docker run ghcr.io/lukasmalkmus/horcrux
```

  [5]: https://github.com/lukasmalkmus/horcrux/pkgs/container/horcrux

### Validate installation

In all cases the installation can be validated by running `horcrux -v` in the
terminal:

```shell
horcrux version 1.0.0
```

## Usage

```shell
horcrux [flags] [commands]
```

Help on flags and commands:

```shell
horcrux --help
```

### Performance

As of today, the implementation isn't suitable for large files. Shamir's Secret
Sharing algorithm is very computation intesive and takes most of the time.
Below are some benchmarks (MacBook Pro, 2,8 GHz Quad-Core i7, 16 GB):

```
name          time/op
Split64KB-8    364ms ± 6%
Split1MB-8     476ms ± 1%
Split128MB-8   18.4s ± 1%
Split1GB-8      160s ± 8%

name          alloc/op
Split64KB-8    135MB ± 0%
Split1MB-8     145MB ± 0%
Split128MB-8  1.48GB ± 0%
Split1GB-8    10.9GB ± 0%

name          allocs/op
Split64KB-8    65.6k ± 0%
Split1MB-8     1.05M ± 0%
Split128MB-8    134M ± 0%
Split1GB-8     1.07G ± 0%
```

Splitting a 1GB file takes up to 3 minutes. Also the memory consumption is a lot
higher than the size of the file which is being processed. An `io.Reader`
based implementation is needed to fix this but this requires multiple tweaks,
especially to Shamir's Secret Sharing implementation.

## Contributing

Feel free to submit PRs or to fill issues. Every kind of help is appreciated. 

Before committing, `make` should run without any issues.

## License

&copy; Lukas Malkmus, 2021

Distributed under MIT License (`The MIT License`).

See [LICENSE](LICENSE) for more information.

[![License Status][license_status_badge]][license_status]

<!-- Badges -->

[gopkg]: https://pkg.go.dev/github.com/lukasmalkmus/horcrux
[gopkg_badge]: https://img.shields.io/badge/doc-reference-007d9c?logo=go&logoColor=white&style=flat-square
[go_workflow]: https://github.com/lukasmalkmus/horcrux/actions/workflows/push.yml
[go_workflow_badge]: https://img.shields.io/github/workflow/status/lukasmalkmus/horcrux/Push?style=flat-square&ghcache=unused
[coverage]: https://codecov.io/gh/lukasmalkmus/horcrux
[coverage_badge]: https://img.shields.io/codecov/c/github/lukasmalkmus/horcrux.svg?style=flat-square&ghcache=unused
[report]: https://goreportcard.com/report/github.com/lukasmalkmus/horcrux
[report_badge]: https://goreportcard.com/badge/github.com/lukasmalkmus/horcrux?style=flat-square&ghcache=unused
[release]: https://github.com/lukasmalkmus/horcrux/releases/latest
[release_badge]: https://img.shields.io/github/release/lukasmalkmus/horcrux.svg?style=flat-square&ghcache=unused
[license]: https://opensource.org/licenses/MIT
[license_badge]: https://img.shields.io/github/license/lukasmalkmus/horcrux.svg?color=blue&style=flat-square&ghcache=unused
[license_status]: https://app.fossa.com/projects/git%2Bgithub.com%2Flukasmalkmus%2Fhorcrux?ref=badge_shield
[license_badge]: https://img.shields.io/github/license/lukasmalkmus/horcrux.svg?color=blue&style=flat-square&ghcache=unused
[license_status]: https://app.fossa.com/projects/git%2Bgithub.com%2Flukasmalkmus%2Fhorcrux
[license_status_badge]: https://app.fossa.com/api/projects/git%2Bgithub.com%2Flukasmalkmus%2Fhorcrux.svg?type=large&ghcache=unused
