DESTDIR ?= /usr/local

RONN := $(shell which ronn >/dev/null 2>&1 && echo "ronn -w --manual=direnv" || echo "@echo 'Could not generate manpage because ronn is missing. gem install ronn' || ")
RONNS = $(wildcard man/*.ronn)
ROFFS = $(RONNS:.ronn=)

.PHONY: all man html test release install
#all: build man test
all: build man

build: direnv

direnv: *.go
	go fmt
	go build -o direnv

clean:
	rm -f direnv

%.1: %.1.ronn
	$(RONN) -r $<

man: $(ROFFS)

test: build
	go test
	./test/direnv-test.sh

release: build
	./script/release `./direnv version`
	git tag v`./direnv version`

install: all
	install -d bin $(DESTDIR)/bin
	install -d man $(DESTDIR)/share/man/man1
	install direnv $(DESTDIR)/bin
	cp -R man/*.1 $(DESTDIR)/share/man/man1

