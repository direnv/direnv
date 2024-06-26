############################################################################
# Variables
############################################################################

# Set this to change the target installation path
PREFIX   = /usr/local
BINDIR   = ${PREFIX}/bin
SHAREDIR = ${PREFIX}/share
MANDIR   = ${SHAREDIR}/man
DISTDIR ?= dist

# filename of the executable
exe = direnv$(shell go env GOEXE)

# Override the go executable
GO = go

# BASH_PATH can also be passed to hard-code the path to bash at build time

SHELL = bash

############################################################################
# Common
############################################################################

.PHONY: all
all: build man

export GO111MODULE=on

############################################################################
# Build
############################################################################

.PHONY: build
build: direnv

.PHONY: clean
clean:
	rm -rf \
		.gopath \
		direnv

GO_LDFLAGS =

ifeq ($(shell uname), Darwin)
	# Fixes DYLD_INSERT_LIBRARIES issues
	# See https://github.com/direnv/direnv/issues/194
	GO_LDFLAGS += -linkmode=external
endif

ifdef BASH_PATH
	GO_LDFLAGS += -X main.bashPath=$(BASH_PATH)
endif

ifneq ($(strip $(GO_LDFLAGS)),)
	GO_BUILD_FLAGS = -ldflags '$(GO_LDFLAGS)'
endif

SOURCES = $(wildcard *.go internal/*/*.go pkg/*/*.go)

direnv: $(SOURCES)
	$(GO) build $(GO_BUILD_FLAGS) -o $(exe)

############################################################################
# Format all the things
############################################################################
.PHONY: fmt fmt-go fmt-sh
fmt: fmt-go fmt-sh

fmt-go:
	$(GO) fmt

fmt-sh:
	@command -v shfmt >/dev/null || (echo "Could not format stdlib.sh because shfmt is missing. Run: go get -u mvdan.cc/sh/cmd/shfmt"; false)
	shfmt -i 2 -w stdlib.sh

.PHONY: fmt-go-lint
fmt-go-lint: fmt-go
	golangci-lint run --fix

############################################################################
# Documentation
############################################################################

man_md = $(wildcard man/*.md)
roffs = $(man_md:.md=)

.PHONY: man
man: $(roffs)

%.1: %.1.md
	@command -v go-md2man >/dev/null || (echo "Could not generate man page because go-md2man is missing. Run: go get -u github.com/cpuguy83/go-md2man/v2"; false)
	go-md2man -in $< -out $@

############################################################################
# Testing
############################################################################

tests = \
				test-shellcheck \
				test-stdlib \
				test-go \
				test-go-lint \
				test-go-fmt \
				test-bash \
				test-elvish \
				test-fish \
				test-tcsh \
				test-zsh \
				test-pwsh \
				test-mx

# Skip few checks for IBM Z mainframe's z/OS aka OS/390
ifeq ($(shell uname), OS/390)
	tests = \
		test-stdlib \
		test-go \
		test-go-fmt \
		test-bash
endif

# Skip few checks for MSYS2
ifneq ($(MSYSTEM),)
	tests = \
		test-shellcheck \
		test-stdlib \
		test-go \
		test-go-lint \
		test-go-fmt \
		test-bash
endif

# Skip few checks for Cygwin
ifeq ($(shell uname -o), Cygwin)
	tests = \
		test-shellcheck \
		test-stdlib \
		test-go \
		test-go-lint \
		test-go-fmt \
		test-bash
endif

.PHONY: $(tests)
test: build $(tests)
	@echo
	@echo SUCCESS!

test-shellcheck:
	shellcheck stdlib.sh
	shellcheck ./test/stdlib.bash

test-stdlib: build
	./test/stdlib.bash

test-go:
	$(GO) test -v ./...

test-go-lint:
	golangci-lint run

test-bash:
	bash ./test/direnv-test.bash

# Needs elvish 0.12+
test-elvish:
	elvish ./test/direnv-test.elv

test-fish:
	fish ./test/direnv-test.fish

test-tcsh:
	tcsh -e ./test/direnv-test.tcsh

test-zsh:
	zsh ./test/direnv-test.zsh

test-pwsh:
	pwsh ./test/direnv-test.ps1

test-mx:
	murex -trypipe ./test/direnv-test.mx

############################################################################
# Cygpath (MSYS2/Cygwin)
############################################################################

.ONESHELL:
.PHONY: init-cygpath-msys2-ci
init-cygpath-msys2-ci:
	@set -e
	# pacman -S fish
	# pacman -S zsh
	# pacman -S ruby
	# pacman -S python
	# pacman -S go
	# pacman -S make

.ONESHELL:
.PHONY: init-cygpath
init-cygpath:
	@set -e
	command -v scoop >/dev/null || (echo "Please install Scoop first."; false)
	scoop install main/make
	scoop install main/golangci-lint
	scoop install main/shfmt
	scoop install main/shellcheck
	scoop install main/go
	scoop install main/msys2
	scoop install main/python
	scoop install main/ruby
	scoop install main/hyperfine

.ONESHELL:
.PHONY: test-cygpath-env
test-cygpath-env:
	@set -e
	command -v cygpath >/dev/null || (echo "cygpath is missing"; false)
	command -v bash >/dev/null || (echo "bash is missing"; false)
	command -v make >/dev/null || (echo "make is missing"; false)
	command -v shfmt >/dev/null || (echo "shfmt is missing, do you want to run make init-cygpath?"; false)
	command -v golangci-lint >/dev/null || (echo "golangci-lint is missing, do you want to run make init-cygpath?"; false)
	command -v shellcheck >/dev/null || (echo "shellcheck is missing, do you want to run make init-cygpath?"; false)
	echo "Environment check done."

.PHONY: test-cygpath-go
test-cygpath-go: test-cygpath-env fmt fmt-go-lint
	$(GO) test -v ./internal/cmd -run TestCygpath

.PHONY: test-cygpath-go-bench
test-cygpath-go-bench: test-cygpath-env fmt fmt-go-lint
	$(GO) test -run XXX -bench=. ./internal/cmd

.PHONY: test-cygpath-bash
test-cygpath-bash: test-cygpath-env fmt-sh build test-shellcheck test-stdlib test-bash

.PHONY: test-cygpath
test-cygpath: test-cygpath-go test-cygpath-bash

# Run like this for accurate measurements on MSYS2:
# MSYS2_ENV_CONV_EXCL="*" make test-cygpath-bench
.ONESHELL:
.PHONY: test-cygpath-bench
test-cygpath-bench: SHELL:=/usr/bin/bash
test-cygpath-bench: build
	@set -euxo pipefail
	direnv_path="$(realpath ./direnv.exe)"
	export direnv_path
	echo "$$direnv_path"
	echo "MSYS2_ENV_CONV_EXCL=${MSYS2_ENV_CONV_EXCL:-}"
	bench_1() { 
		eval "$(MSYS2_ENV_CONV_EXCL="*" "$$direnv_path" export bash)"
	}
	export -f bench_1
	bench_4() { 
		MSYS2_ENV_CONV_EXCL="*" "$$direnv_path" exec . hostname
	}
	export -f bench_4
	bench_2() { 
		eval "$("$$direnv_path" export bash)"
	}
	export -f bench_2
	bench_3() { 
		"$$direnv_path" exec . hostname
	}
	export -f bench_3
	cd ./test/scenarios/cygpath-perf 
	MSYS2_ENV_CONV_EXCL="*" "$$direnv_path" allow
	hyperfine --show-output --export-markdown /tmp/bench_1.md --shell=bash --runs 10 --warmup 1 bench_1
	hyperfine --show-output --export-markdown /tmp/bench_4.md --shell=bash --runs 10 --warmup 1 bench_4
	unset MSYS2_ENV_CONV_EXCL
	"$$direnv_path" allow
	hyperfine --show-output --export-markdown /tmp/bench_2.md --shell=bash --runs 10 --warmup 1 bench_2
	hyperfine --show-output --export-markdown /tmp/bench_3.md --shell=bash --runs 10 --warmup 1 bench_3
	cat /tmp/bench_1.md
	cat /tmp/bench_4.md
	cat /tmp/bench_2.md
	cat /tmp/bench_3.md

############################################################################
# Installation
############################################################################

.PHONY: install
install: all
	install -d $(DESTDIR)$(BINDIR)
	install $(exe) $(DESTDIR)$(BINDIR)
	install -d $(DESTDIR)$(MANDIR)/man1
	cp -R man/*.1 $(DESTDIR)$(MANDIR)/man1
	install -d $(DESTDIR)$(SHAREDIR)/fish/vendor_conf.d
	echo "$(BINDIR)/direnv hook fish | source" > $(DESTDIR)$(SHAREDIR)/fish/vendor_conf.d/direnv.fish

.PHONY: dist
dist:
	@command -v gox >/dev/null || (echo "Could not generate dist because gox is missing. Run: go get -u github.com/mitchellh/gox"; false)
	CGO_ENABLED=0 GOFLAGS="-trimpath" \
		gox -rebuild -ldflags="-s -w" -output "$(DISTDIR)/direnv.{{.OS}}-{{.Arch}}" \
		-osarch darwin/amd64 \
		-osarch darwin/arm64 \
		-osarch freebsd/386 \
		-osarch freebsd/amd64 \
		-osarch freebsd/arm \
		-osarch linux/386 \
		-osarch linux/amd64 \
		-osarch linux/arm \
		-osarch linux/arm64 \
		-osarch linux/mips \
		-osarch linux/mips64 \
		-osarch linux/mips64le \
		-osarch linux/mipsle \
		-osarch linux/ppc64 \
		-osarch linux/ppc64le \
		-osarch linux/s390x \
		-osarch netbsd/386 \
		-osarch netbsd/amd64 \
		-osarch netbsd/arm \
		-osarch openbsd/386 \
		-osarch openbsd/amd64 \
		-osarch windows/386 \
		-osarch windows/amd64 \
		&& true
