
SOURCE=$(wildcard src/*c)

all: shenv | vendor/lua-5.1.4

shenv: $(SOURCE)
	$(CC) -o $@ $(SOURCE)

# TODO: integrate the lib/ source into a file
src/lib.c:
	@echo "Yoo"

vendor:
	mkdir vendor

vendor/lua-5.1.4.tar.gz: | vendor
	wget -O $@ "http://www.lua.org/ftp/lua-5.1.4.tar.gz"

vendor/patch-lua-5.1.4-2: | vendor
	wget -O $@ "http://www.lua.org/ftp/patch-lua-5.1.4-2"

vendor/lua-5.1.4: vendor/lua-5.1.4.tar.gz vendor/patch-lua-5.1.4-2
	cd $(dir $@); tar xzvf $(abspath $<)
	cd $@/src; patch -p0 < $(abspath vendor/patch-lua-5.1.4-2)
