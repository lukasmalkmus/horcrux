# lukasmalkmus/horcrux

> A security question based secret sharing utility. - by **[Lukas Malkmus]**

[![Build Status][build_badge]][build]
[![Coverage Status][coverage_badge]][coverage]
[![Go Report][report_badge]][report]
[![GoDoc][docs_badge]][docs]
[![Latest Release][release_badge]][release]
[![License][license_badge]][license]
[![License Status][license_status_badge]][license_status]

---

## Table of Contents

1. [Introduction](#introduction)
1. [Usage](#usage)
1. [Contributing](#contributing)
1. [License](#license)

## Introduction

*horcrux* is a security question based secret sharing utility. The idea and
package code is inspired and mostly taken from the abandoned [horcrux] package
by [Coda Hale].

*horcrux* splits a secret into multiple fragments and associates every fragment
with a security question. The answer to that question is used to encrypt the
fragment using ChaCha20Poly1305. Only a given number of fragments is needed to
fully restore the original secret.

## Usage

### Installation

The easiest way to run *horcrux* is by grabbing the latest standalone binary
from the [release page][release].

This project uses native [go mod] support for vendoring and requires a working
`go` toolchain installation when installing via `go get` or from source.

#### Install using `go get`

```bash
GO111MODULE=on go install github.com/lukasmalkmus/horcrux/cmd/horcrux
```

#### Install from source

```bash
git clone https://github.com/lukasmalkmus/horcrux.git
cd horcrux
make # Run a full build including code formatting, linting and testing
make build # Build production binary
make install # Build and install binary into $GOPATH
```

#### Validate installation

The installation can be validated by running `horcrux version` in the terminal.

### Using the application

```bash
horcrux [flags] [commands]
```

Help on flags and commands:

```bash
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

Splitting a 1 GB file takes up to 3 minutes. Also the memory consumption is a
bit higher than the size of the file which is being processed. An `io.Reader`
based implementation is needed to fix this but this requires multiple tweaks,
especially Shamir's Secret Sharing implementation.

Splitting 

## Contributing

Feel free to submit PRs or to fill Issues. Every kind of help is appreciated.

## License

© Lukas Malkmus, 2019

Distributed under MIT License (`The MIT License`).

See [LICENSE](LICENSE) for more information.

[![License Status Large][license_status_large_badge]][license_status_large]

<!-- Links -->
[Lukas Malkmus]: https://github.com/lukasmalkmus
[Coda Hale]: https://github.com/codahale
[horcrux]: https://github.com/codahale/horcrux
[go mod]: https://golang.org/cmd/go/#hdr-Module_maintenance


<!-- Badges -->
[build]: https://travis-ci.com/lukasmalkmus/horcrux
[build_badge]: https://img.shields.io/travis/com/lukasmalkmus/horcrux.svg?style=flat-square
[coverage]: https://codecov.io/gh/lukasmalkmus/horcrux
[coverage_badge]: https://img.shields.io/codecov/c/github/lukasmalkmus/horcrux.svg?style=flat-square
[report]: https://goreportcard.com/report/github.com/lukasmalkmus/horcrux
[report_badge]: https://goreportcard.com/badge/github.com/lukasmalkmus/horcrux?style=flat-square
[docs]: https://godoc.org/github.com/lukasmalkmus/horcrux
[docs_badge]: https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square
[release]: https://github.com/lukasmalkmus/horcrux/releases
[release_badge]: https://img.shields.io/github/release/lukasmalkmus/horcrux.svg?style=flat-square
[license]: https://opensource.org/licenses/MIT
[license_badge]: https://img.shields.io/github/license/lukasmalkmus/horcrux.svg?color=blue&style=flat-square
[license_status]: https://app.fossa.com/projects/git%2Bgithub.com%2Flukasmalkmus%2Fhorcrux?ref=badge_shield
[license_status_badge]: https://app.fossa.com/api/projects/git%2Bgithub.com%2Flukasmalkmus%2Fhorcrux.svg
[license_status_large]: https://app.fossa.com/projects/git%2Bgithub.com%2Flukasmalkmus%2Fhorcrux?ref=badge_large
[license_status_large_badge]: https://app.fossa.com/api/projects/git%2Bgithub.com%2Flukasmalkmus%2Fhorcrux.svg?type=large
