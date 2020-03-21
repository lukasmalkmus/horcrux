# TOOLCHAIN
GO				:= CGO_ENABLED=0 GOFLAGS=-mod=vendor GOBIN=$(CURDIR)/bin go
GO_BIN_IN_PATH	:= CGO_ENABLED=0 GOFLAGS=-mod=vendor go
GO_NO_VENDOR	:= CGO_ENABLED=0 GOBIN=$(CURDIR)/bin go
GOFMT			:= $(GO)fmt

# ENVIRONMENT
VERBOSE 	=
GOPATH		:= $(GOPATH)
GOOS		?= $(shell echo $(shell uname -s) | tr A-Z a-z)
GOARCH		?= amd64
MOD_NAME	:= github.com/lukasmalkmus/horcrux

# APPLICATION INFORMATION
BUILD_DATE	:= $(shell date -u '+%Y-%m-%d_%H:%M:%S')
REVISION	:= $(shell git rev-parse --short HEAD)
RELEASE		:= $(shell cat RELEASE)
USER		:= $(shell whoami)

# TOOLS
BENCHSTAT		:= bin/benchstat
GOLANGCI_LINT	:= bin/golangci-lint
GOTESTSUM		:= bin/gotestsum

# MISC
COVERPROFILE	:= coverage.out
DIRTY			:= $(shell git diff-index --quiet HEAD || echo "untracked")

# FLAGS
GOFLAGS			:= -buildmode=exe -tags=netgo -installsuffix=cgo -trimpath \
					-ldflags='-s -w -extldflags "-static" \
					-X $(MOD_NAME)/pkg/version.release=$(RELEASE) \
					-X $(MOD_NAME)/pkg/version.revision=$(REVISION) \
					-X $(MOD_NAME)/pkg/version.buildDate=$(BUILD_DATE) \
					-X $(MOD_NAME)/pkg/version.buildUser=$(USER)'

GOTESTSUM_FLAGS	:= --jsonfile tests.json --junitfile junit.xml
GO_TEST_FLAGS 	:= -race -coverprofile=$(COVERPROFILE)
GO_BENCH_FLAGS  := -run='NONE' -bench=. -benchmem -count=10

# DEPENDENCIES
GOMODDEPS = go.mod go.sum

# Enable verbose test output if explicitly set.
ifdef VERBOSE
	GOTESTSUM_FLAGS	+= --format=standard-verbose
endif

# FUNCS
# func go-list-pkg-sources(package)
go-list-pkg-sources = $(GO) list $(GOFLAGS) -f '{{ range $$index, $$filename := .GoFiles }}{{ $$.Dir }}/{{ $$filename }} {{end}}' $(1)
# func go-pkg-sourcefiles(package)
go-pkg-sourcefiles = $(shell $(call go-list-pkg-sources,$(strip $1)))

.PHONY: all
all: dep fmt lint test build ## Run dep, fmt, lint, test and build.

.PHONY: bench
bench: $(BENCHSTAT) ## Run all benchmarks and compare the benchmark results of a dirty workspace to the ones of a clean workspace if available.
	@echo ">> running benchmarks"
	@mkdir -p benchmarks
ifneq ($(DIRTY),untracked)
	@$(GO) test $(GO_BENCH_FLAGS) ./... > benchmarks/$(REVISION).txt
	@$(BENCHSTAT) benchmarks/$(REVISION).txt
else
	@$(GO) test $(GO_BENCH_FLAGS) ./... > benchmarks/$(REVISION)-dirty.txt
ifneq (,$(wildcard benchmarks/$(REVISION).txt))
	@$(BENCHSTAT) benchmarks/$(REVISION).txt benchmarks/$(REVISION)-dirty.txt
else
	@$(BENCHSTAT) benchmarks/$(REVISION)-dirty.txt
endif
endif

.PHONY: build
build: .build/horcrux-$(GOOS)-$(GOARCH) ## Build all binaries.

.PHONY: clean
clean: ## Remove build and test artifacts.
	@echo ">> cleaning up artifacts"
	@rm -rf .build $(COVERPROFILE) tests.json junit.xml

.PHONY: cover
cover: $(COVERPROFILE) ## Calculate the code coverage score.
	@echo ">> calculating code coverage"
	@$(GO) tool cover -func=$(COVERPROFILE)

.PHONY: dep-clean
dep-clean: ## Remove obsolete dependencies.
	@echo ">> cleaning dependencies"
	@$(GO) mod tidy

.PHONY: dep-upgrade
dep-upgrade: ## Upgrade all direct dependencies to their latest version.
	@echo ">> upgrading dependencies"
	@$(GO) get $(shell $(GO_NO_VENDOR) list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all)
	@$(GO) mod vendor
	@make dep

.PHONY: dep
dep: dep-clean dep.stamp ## Install and verify dependencies and remove obsolete ones.

dep.stamp: $(GOMODDEPS)
	@echo ">> installing dependencies"
	@$(GO) mod download
	@$(GO) mod verify
	@$(GO) mod vendor
	@touch $@

.PHONY: fmt
fmt: ## Format and simplify the source code using `gofmt`.
	@echo ">> formatting code"
	@! $(GOFMT) -s -w $(shell find . -path -prune -o -name '*.go' -print) | grep '^'

.PHONY: install
install: $(GOPATH)/bin/horcrux ## Install the binary into the $GOPATH/bin directory.

.PHONY: lint
lint: $(GOLANGCI_LINT) ## Lint the source code.
	@echo ">> linting code"
	@$(GOLANGCI_LINT) run

.PHONY: test
test: $(GOTESTSUM) ## Run all tests. Run with VERBOSE=1 to get verbose test output ('-v' flag).
	@echo ">> running tests"
	@$(GOTESTSUM) $(GOTESTSUM_FLAGS) -- $(GO_TEST_FLAGS) ./...

.PHONY: tools
tools: $(BENCHSTAT) $(GOLANGCI_LINT) $(GOTESTSUM) ## Install all tools into the projects local $GOBIN directory.

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# BUILD TARGETS

.build/horcrux-darwin-amd64: dep.stamp $(call go-pkg-sourcefiles, ./...)
	@echo ">> building horcrux production binary for darwin/amd64"
	@GOOS=darwin GOARCH=amd64 $(GO) build $(GOFLAGS) -o .build/horcrux-darwin-amd64 ./cmd/horcrux

.build/horcrux-linux-amd64: dep.stamp $(call go-pkg-sourcefiles, ./...)
	@echo ">> building horcrux production binary for linux/amd64"
	@GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -o .build/horcrux-linux-amd64 ./cmd/horcrux

.build/horcrux-windows-amd64: dep.stamp $(call go-pkg-sourcefiles, ./...)
	@echo ">> building horcrux production binary for windows/amd64"
	@GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) -o .build/horcrux-windows-amd64 ./cmd/horcrux

# INSTALL TARGETS

$(GOPATH)/bin/horcrux: dep.stamp $(call go-pkg-sourcefiles, ./...)
	@echo ">> installing horcrux binary"
	@$(GO_BIN_IN_PATH) install $(GOFLAGS) ./cmd/horcrux

# TEST TARGETS

$(COVERPROFILE):
	@make test

# TOOLS

$(BENCHSTAT): dep.stamp $(call go-pkg-sourcefiles, ./vendor/golang.org/x/perf/cmd/benchstat)
	@echo ">> installing benchstat"
	@$(GO) install golang.org/x/perf/cmd/benchstat

$(GOLANGCI_LINT): dep.stamp $(call go-pkg-sourcefiles, ./vendor/github.com/golangci/golangci-lint/cmd/golangci-lint)
	@echo ">> installing golangci-lint"
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint

$(GOTESTSUM): dep.stamp $(call go-pkg-sourcefiles, ./vendor/gotest.tools/gotestsum)
	@echo ">> installing gotestsum"
	@$(GO) install gotest.tools/gotestsum
