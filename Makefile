VERSION := 0.1.`git log --oneline | wc -l | sed -e "s/ //g"`

.PHONY: all
all:
	@echo TODO

.PHONY: release
release:
	git tag v$(VERSION)

.PHONY: test
test:
	./test/direnv-test.sh
