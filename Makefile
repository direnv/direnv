DESTDIR ?= /usr/local

RONN := $(shell which ronn >/dev/null 2>&1 && echo "ronn -w --manual=direnv --organization=0x2a" || echo "@echo 'Could not generate manpage because ronn is missing. gem install ronn' || ")
RONNS = $(wildcard man/*.ronn)
ROFFS = $(RONNS:.ronn=)

.PHONY: all man html test release install gh-pages
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

html:
	$(RONN) -W5 -s toc man/*.ronn

# FIXME: restore the integration tests ./test/direnv-test.sh
test:
	go test

release: build
	./script/release `./direnv version`
	git tag v`./direnv version`

gh-pages: html
	git stash
	git checkout gh-pages
	mv man/*.html .
	git add *.html
	git commit -m "auto"
	git checkout master
	git stash pop || true

install: all
	install -d bin $(DESTDIR)/bin
	install -d man $(DESTDIR)/share/man/man1
	cp direnv $(DESTDIR)/bin
	cp -R man/*.1 $(DESTDIR)/share/man/man1

