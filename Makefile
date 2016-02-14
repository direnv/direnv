DESTDIR ?= /usr/local

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

man:
	@which md2man-rake >/dev/null || (echo "Could not generate man page because md2man is missing, gem install md2man"; false)
	md2man-rake md2man:man

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

