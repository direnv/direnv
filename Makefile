VERSION := 0.1.$(shell git log --oneline | wc -l | sed -e "s/ //g")

DESTDIR ?= /usr/local

RONN := $(shell which ronn >/dev/null 2>&1 && echo "ronn -w --manual=direnv --organization=0x2a" || echo "@echo 'Could not generate manpage because ronn is missing. gem install ronn' || ")
RONNS = $(wildcard man/*.ronn)
ROFFS = $(RONNS:.ronn=)

.PHONY: all man html test release install gh-pages
all: man test

%.1: %.1.ronn
	$(RONN) -r $<

man: $(ROFFS)

html:
	$(RONN) -W5 -s toc man/*.ronn

# Maybe use https://github.com/bmizerany/roundup
test:
	./test/direnv-test.sh

release:
	git tag v$(VERSION)

gh-pages: html
	git stash
	git checkout gh-pages
	mv man/*.html .
	git add *.html
	git commit -m "$(VERSION)"
	git checkout master
	git stash pop || true

install:
	install -d bin $(DESTDIR)/bin
	install -d libexec $(DESTDIR)/libexec
	install -d man $(DESTDIR)/share/man/man1
	cp -R bin/* $(DESTDIR)/bin
	cp -R libexec/* $(DESTDIR)/libexec
	cp -R man/*.1 $(DESTDIR)/share/man/man1
