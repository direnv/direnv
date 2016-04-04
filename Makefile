DESTDIR ?= /usr/local

MAN_MD = $(wildcard man/*.md)
ROFFS = $(MAN_MD:.md=)

ifeq ($(shell uname), Darwin)
	# Fixes DYLD_INSERT_LIBRARIES issues
	# See https://github.com/direnv/direnv/issues/194
	GO_FLAGS += -ldflags -linkmode=external
endif

.PHONY: all man html test install dist
#all: build man test
all: build man

build: direnv

stdlib.go: stdlib.sh
	cat $< | ./script/str2go main STDLIB $< > $@

direnv: stdlib.go *.go
	go fmt
	go build $(GO_FLAGS) -o direnv

clean:
	rm -f direnv

%.1: %.1.md
	@which md2man-roff >/dev/null || (echo "Could not generate man page because md2man is missing, gem install md2man"; false)
	md2man-roff $< > $@

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

