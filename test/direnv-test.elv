#!/usr/bin/env elvish

E:TEST_DIR = (path-dir (src)[path])
E:PATH = (path-dir $E:TEST_DIR):$E:PATH

cd $E:TEST_DIR

## reset the direnv loading if any
set-env DIRENV_CONFIG $pwd
unset-env DIRENV_BASH
unset-env DIRENV_DIR
unset-env DIRENV_MTIME
unset-env DIRENV_WATCHES
unset-env DIRENV_DIFF

set-env XDG_CONFIG_HOME $E:TEST_DIR
mkdir -p $E:XDG_CONFIG_HOME/direnv
touch $E:XDG_CONFIG_HOME/direnv/direnvrc

fn direnv-eval {
	try {
		m = (direnv export elvish | from-json)
		keys $m | each [k]{
			if (==s $k 'null') {
				unset-env $k
			} else {
				set-env $k $m[$k]
			}
		}
	} except e {
		nop
	}
}

fn test-debug {
	if (==s $E:DIRENV_DEBUG "1") {
		echo
	}
}

fn test-eq [a b]{
	if (!=s $a $b) {
		fail "FAILED: '"$a"' == '"$b"'"
	}
}

fn test-neq [a b]{
	if (==s $a $b) {
		fail "FAILED: '"$a"' != '"$b"'"
	}
}

fn test-scenario [name fct]{
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
} except e {
	nop
}

direnv-eval

test-scenario base {
	echo "Setting up"
	direnv-eval
	test-eq $E:HELLO "world"

	E:WATCHES=$E:DIRENV_WATCHES

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

## TODO: special-vars
## TODO: dump
## TODO: empty-var
## TODO: empty-var-unset

test-scenario "missing-file-source-env" {
	direnv-eval
}
