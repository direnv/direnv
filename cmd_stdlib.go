package main

import (
	"fmt"
)

// `direnv stdlib`
var CmdStdlib = &Cmd{
	Name:    "stdlib",
	Desc:    "Outputs the stdlib that is available in the .envrc",
	Private: true,
	Fn: func(env Env, args []string) (err error) {
		var config *Config
		if config, err = LoadConfig(env); err != nil {
			return
		}

		fmt.Printf(STDLIB, config.SelfPath)
		return
	},
}

const STDLIB = `
# These are the commands available in an .envrc context
set -e
DIRENV_PATH="%s"

# Determines if "something" is availabe as a command
#
# Usage: has something
has() {
	type "$1" &>/dev/null
}

# Usage: expand_path ./rel/path [RELATIVE_TO]
# RELATIVE_TO is $PWD by default
expand_path() {
	"$DIRENV_PATH" expand_path "$@"
}

# Loads a .env in the current environment
#
# Usage: dotenv
dotenv() {
	eval "$("$DIRENV_PATH" dotenv "$@")"
}

# Usage: user_rel_path /Users/you/some_path => ~/some_path
user_rel_path() {
	local path="${1#-}"

	if [ -z "$path" ]; then return; fi

	if [ -n "$HOME" ]; then
		local rel_path="${path#$HOME}"
		if [ "$rel_path" != "$path" ]; then
			path="~${rel_path}"
		fi
	fi

	echo $path
}

# Usage: find_up FILENAME
find_up() {
	(
		cd "$(pwd -P 2>/dev/null)"
		while true; do
			if [ -f "$1" ]; then
				echo $PWD/$1
				return 0
			fi
			if [ "$PWD" = "/" ] || [ "$PWD" = "//" ]; then
				return 1
			fi
			cd ..
		done
	)
}

# Inherit another .envrc
#
# Usage: source_env <FILE_OR_DIR_PATH>
source_env() {
	local rcfile="$1"
	local rcpath="${1/#\~/$HOME}"
	if ! [ -f "$rcpath" ]; then
		rcfile="$rcfile/.envrc"
		rcpath="$rcpath/.envrc"
	fi
	echo "direnv: loading $rcfile"
	pushd "$(dirname "$rcpath")" > /dev/null
	. "./$(basename "$rcpath")"
	popd > /dev/null
}

# Inherits the first .envrc (or given FILENAME) it finds in the path
#
# Usage: source_up [FILENAME]
source_up() {
	local file="$1"
	if [ -z "$file" ]; then
		file=".envrc"
	fi
	local path="$(cd .. && find_up "$file")"
	if [ -n "$path" ]; then
		source_env "$(user_rel_path "$path")"
	fi
}

# Safer PATH handling
#
# Usage: PATH_add PATH
# Example: PATH_add bin
PATH_add() {
	export PATH="$(expand_path "$1"):$PATH"
}

# Safer path handling
#
# Usage: path_add VARNAME PATH
# Example: path_add LD_LIBRARY_PATH ./lib
path_add() {
	local old_paths="${!1}"
	local path="$(expand_path "$2")"

	if [ -z "$old_paths" ]; then
		old_paths="$path"
	else
		old_paths="$path:$old_paths"
	fi

	export $1="$old_paths"
}

#
# Usage: load_prefix PATH
load_prefix() {
	local path="$(expand_path "$1")"
	path_add CPATH "$path/include"
	path_add LD_LIBRARY_PATH "$path/lib"
	path_add LIBRARY_PATH "$path/lib"
	path_add MANPATH "$path/man"
	path_add MANPATH "$path/share/man"
	path_add PATH "$path/bin"
	path_add PKG_CONFIG_PATH "$path/lib/pkgconfig"
}

# Pre-programmed project layout. Add your own in your ~/.direnvrc.
#
# Usage: layout TYPE
layout() {
	eval "layout_$1"
}

# Usage: layout ruby
layout_ruby() {
	# TODO: ruby_version should be the ABI version
	local ruby_version="$(ruby -e"puts (defined?(RUBY_ENGINE) ? RUBY_ENGINE : 'ruby') + '-' + RUBY_VERSION")"

	export GEM_HOME="$PWD/.direnv/${ruby_version}"
	export BUNDLE_BIN="$PWD/.direnv/bin"

	PATH_add ".direnv/${ruby_version}/bin"
	PATH_add ".direnv/bin"
}

# Usage: layout python
layout_python() {
	if ! [ -d .direnv/virtualenv ]; then
		virtualenv --no-site-packages --distribute .direnv/virtualenv
		virtualenv --relocatable .direnv/virtualenv
	fi
	source .direnv/virtualenv/bin/activate
}

# Usage: layout node
layout_node() {
	PATH_add node_modules/.bin
}

# Intended to load external dependencies into the environment.
#
# Usage: use PROGRAM_NAME VERSION
# Example: use ruby 1.9.3
use() {
	local cmd="$1"
	echo "Using $@"
	shift
	use_$cmd "$@"
}

# Usage: use rbenv
use_rbenv() {
	eval "$(rbenv init -)"
}

# Sources rvm on first call. Should work like the rvm command-line.
rvm() {
	unset rvm
	if [ -n "${rvm_scripts_path:-}" ]; then
		source "${rvm_scripts_path}/rvm"
	elif [ -n "${rvm_path:-}" ]; then
		source "${rvm_path}/scripts/rvm"
	else
		source "$HOME/.rvm/scripts/rvm"
	fi
	rvm "$@"
}

## Load the global ~/.direnvrc if present
if [ -f "$HOME/.direnvrc" ]; then
	source_env "~/.direnvrc" >&2
fi
`
