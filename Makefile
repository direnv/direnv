DESTDIR ?= /usr/local

MAN_MD = $(wildcard man/*.md)
ROFFS = $(MAN_MD:.md=)

GO_LDFLAGS =

ifeq ($(shell uname), Darwin)
	# Fixes DYLD_INSERT_LIBRARIES issues
	# See https://github.com/direnv/direnv/issues/194
	GO_LDFLAGS += -linkmode=external
endif

ifdef BASH_PATH
	GO_LDFLAGS += -X main.bashPath=$(BASH_PATH)
endif

ifdef GO_LDFLAGS
	GO_FLAGS += -ldflags '$(GO_LDFLAGS)'
endif

.PHONY: all man html test install dist
#all: build man test
all: build man

build: direnv

stdlib.go: stdlib.sh
	cat $< | ./script/str2go main STDLIB $< > $@

version.go: version.txt
	echo package main > $@
	echo 'const VERSION = "$(shell cat $<)";' >> $@

direnv: stdlib.go *.go
	go fmt
	go build $(GO_FLAGS) -o direnv

clean:
	rm -f direnv

%.1: %.1.md
	@which go-md2man >/dev/null || (echo "Could not generate man page because go-md2man is missing, `go get -u https://github.com/cpuguy83/go-md2man`"; false)
	go-md2man -in $< -out $@

man: $(ROFFS)

test: build
	go test
	./test/direnv-test.sh

install: all
	install -d $(DESTDIR)/bin
	install -d $(DESTDIR)/share/man/man1
	install direnv $(DESTDIR)/bin
	cp -R man/*.1 $(DESTDIR)/share/man/man1

dist:
	go get github.com/mitchellh/gox
	gox -output "dist/direnv.{{.OS}}-{{.Arch}}"

