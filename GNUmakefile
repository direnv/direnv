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
	@command -v shfmt >/dev/null || (echo "Could not format stdlib.sh because shfmt is missing. Run: go install mvdan.cc/sh/cmd/shfmt@latest"; false)
	shfmt -i 2 -w stdlib.sh

############################################################################
# Documentation
############################################################################

man_md = $(wildcard man/*.md)
roffs = $(man_md:.md=)

.PHONY: man
man: $(roffs)

%.1: %.1.md
	@command -v go-md2man >/dev/null || (echo "Could not generate man page because go-md2man is missing. Run: go install github.com/cpuguy83/go-md2man/v2@latest"; false)
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
		-osarch windows/arm64 \
		&& true

.PHONY: create-release
create-release: dist
	@if [ -z "$$GITHUB_REF_NAME" ]; then \
		echo "GITHUB_REF_NAME is not set. This target is meant to be run in GitHub Actions."; \
		exit 1; \
	fi
	@echo "Extracting release notes from CHANGELOG.md..."
	@release_notes=$$(awk '/^======/{if(headers>0) exit} /^======/{headers++; next} headers>0' CHANGELOG.md); \
	gh release create "$$GITHUB_REF_NAME" \
		--title "Release $$GITHUB_REF_NAME" \
		--notes "$$release_notes" \
		--verify-tag
	gh release upload "$$GITHUB_REF_NAME" $(DISTDIR)/direnv.*
