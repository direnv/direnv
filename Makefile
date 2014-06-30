DESTDIR ?= /usr/local

MAN_MD = $(wildcard man/*.md)
ROFFS = $(MAN_MD:.md=)

.PHONY: all man html test install dist
#all: build man test
all: build man

build: direnv

direnv: *.go
	go fmt
	go build -o direnv

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
	gox -build-toolchain
	gox -output "dist/{{.Dir}}.{{.OS}}-{{.Arch}}"

