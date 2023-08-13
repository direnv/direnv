#!/usr/bin/env elvish

use path

set E:TEST_DIR = (path:dir (src)[name])
set-env XDG_CONFIG_HOME $E:TEST_DIR/config
set-env XDG_DATA_HOME $E:TEST_DIR/data
set E:PATH = (path:dir $E:TEST_DIR):$E:PATH

cd $E:TEST_DIR

## reset the direnv loading if any
set-env DIRENV_CONFIG $pwd
unset-env DIRENV_BASH
unset-env DIRENV_DIR
unset-env DIRENV_FILE
unset-env DIRENV_WATCHES
unset-env DIRENV_DIFF

mkdir -p $E:XDG_CONFIG_HOME/direnv
touch $E:XDG_CONFIG_HOME/direnv/direnvrc

fn direnv-eval {
	try {
		var m = (direnv export elvish | from-json)
		var k
		keys $m | each {|k|
			if $m[$k] {
				set-env $k $m[$k]
			} else {
				unset-env $k
			}
		}
	} catch e {
		nop
	}
}

fn test-debug {
	if (==s $E:DIRENV_DEBUG "1") {
		echo
	}
}

fn test-eq {|a b|
	if (!=s $a $b) {
		fail "FAILED: '"$a"' == '"$b"'"
	}
}

fn test-neq {|a b|
	if (==s $a $b) {
		fail "FAILED: '"$a"' != '"$b"'"
	}
}

fn test-scenario {|name fct|
	cd $E:TEST_DIR/scenarios/$name
	direnv allow
	test-debug
	echo "\n## Testing "$name" ##"
	test-debug

	$fct

	cd $E:TEST_DIR
	direnv-eval
}


### RUN ###

try {
	direnv allow
} catch e {
	nop
}

direnv-eval

test-scenario base {
	echo "Setting up"
	direnv-eval
	test-eq $E:HELLO "world"

	set E:WATCHES = $E:DIRENV_WATCHES

	echo "Reloading (should be no-op)"
	direnv-eval
	test-eq $E:WATCHES $E:DIRENV_WATCHES

	sleep 1

	echo "Updating envrc and reloading (should reload)"
	touch .envrc
	direnv-eval
	test-neq $E:WATCHES $E:DIRENV_WATCHES

	echo "Leaving dir (should clear env set by dir's envrc)"
	cd ..
	direnv-eval
	test-eq $E:HELLO ""
}

test-scenario inherit {
	cp ../base/.envrc ../inherited/.envrc
	direnv-eval
	echo "HELLO should be world:"$E:HELLO
	test-eq $E:HELLO "world"

	sleep 1
	echo "export HELLO=goodbye" > ../inherited/.envrc
	direnv-eval
	test-eq $E:HELLO "goodbye"
}

test-scenario "ruby-layout" {
	direnv-eval
	test-neq $E:GEM_HOME ""
}

test-scenario "space dir" {
	direnv-eval
	test-eq $E:SPACE_DIR "true"
}

test-scenario "child-env" {
	direnv-eval
	test-eq $E:PARENT_PRE "1"
	test-eq $E:CHILD "1"
	test-eq $E:PARENT_POST "1"
	test-eq $E:REMOVE_ME ""
}

test-scenario "utf-8" {
	direnv-eval
	test-eq $E:UTFSTUFF "♀♂"
}

## TODO: special-vars
## TODO: dump
## TODO: empty-var
## TODO: empty-var-unset

test-scenario "missing-file-source-env" {
	direnv-eval
}
