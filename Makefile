VERSION := 0.1.$(shell git log --oneline | wc -l | sed -e "s/ //g")

DESTDIR := /usr/local

RONN := $(shell which ronn >/dev/null 2>&1 && echo "ronn -w --manual=direnv --organization=0x2a" || echo "@echo 'Could not generate manpage because ronn is missing. gem install ronn' || ")
MAN_PAGES := $(shell ls libexec | sed -e "s/\(.*\)/man\/\1.1/g")

.PHONY: all man test release install
all: man test

%.1: %.1.ronn
	$(RONN) -r $<

man: $(MAN_PAGES)

test:
	./test/direnv-test.sh

release:
	git tag v$(VERSION)

install:
	install -d bin $(DESTDIR)/bin
	install -d libexec $(DESTDIR)/libexec
	install -d man $(DESTDIR)/man/man1
	cp -R bin/* $(DESTDIR)/bin
	cp -R libexec/* $(DESTDIR)/libexec
	cp -R man/*.1 $(DESTDIR)/man/man1
