VERSION := 0.1.$(shell git log --oneline | wc -l | sed -e "s/ //g")

DESTDIR := /usr/local

RONN := ronn -w --manual=direnv --organization=0x2a
MAN_PAGES := $(shell ls libexec | sed -e "s/\(.*\)/man\/\1.1/g")

.PHONY: all
all: man test

%.1: %.1.ronn
	$(RONN) -r $@.ronn

.PHONY: man
man: $(MAN_PAGES)

.PHONY: test
test:
	./test/direnv-test.sh

.PHONY: release
release:
	git tag v$(VERSION)

.PHONY: install
install:
	install -d bin $(DESTDIR)/bin
	install -d libexec $(DESTDIR)/libexec
	install -d man $(DESTDIR)/man/man1
	cp -R bin/* $(DESTDIR)/bin
	cp -R libexec/* $(DESTDIR)/libexec
	cp -R man/*.1 $(DESTDIR)/man/man1
