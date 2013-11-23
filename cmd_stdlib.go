package main

import (
	"fmt"
)

// `direnv stdlib`
var CmdStdlib = &Cmd{
	Name: "stdlib",
	Desc: "Displays the stdlib available in the .envrc execution context",
	Fn: func(env Env, args []string) (err error) {
		var config *Config
		if config, err = LoadConfig(env); err != nil {
			return
		}

		fmt.Printf(STDLIB, config.SelfPath)
		return
	},
}

const STDLIB = `# These are the commands available in an .envrc context
set -e
DIRENV_PATH="%s"

# Usage: has <command>
#
# Returns 0 if the <command> is available. Returns 1 otherwise. It can be a
# binary in the PATH or a shell function.
#
# Example:
#
#    if has curl; then
#      echo "Yes we do"
#    fi
#
has() {
	type "$1" &>/dev/null
}

# Usage: expand_path <rel_path> [<relative_to>]
#
# Outputs the absolute path of <rel_path> relaitve to <relative_to> or the 
# current directory.
#
# Example:
#
#    cd /usr/local/games
#    expand_path ../foo
#    # output: /usr/local/foo
#
expand_path() {
	"$DIRENV_PATH" expand_path "$@"
}

# Usage: dotenv [<dotenv>]
#
# Loads a ".env" file into the current environment
#
dotenv() {
	eval "$("$DIRENV_PATH" dotenv "$@")"
}

# Usage: user_rel_path <abs_path>
#
# Transforms an absolute path <abs_path> into a user-relative path if
# possible.
#
# Example:
#
#    echo $HOME
#    # output: /home/user
#    user_rel_path /home/user/my/project
#    # output: ~/my/project
#    user_rel_path /usr/local/lib
#    # output: /usr/local/lib
#
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

# Usage: find_up <filename>
#
# Outputs the path of <filename> when searched from the current directory up to 
# /. Returns 1 if the file has not been found.
#
# Example:
#
#    cd /usr/local/my
#    mkdir -p project/foo
#    touch bar
#    cd project/foo
#    find_up bar
#    # output: /usr/local/my/bar
#
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

# Usage: source_env <file_or_dir_path>
#
# Loads another ".envrc" either by specifying it's path or filename.
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

# Usage: source_up [<filename>]
#
# Loads another ".envrc" if found with the find_up command.
#
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

# Usage: PATH_add <path>
#
# Prepends the expanded <path> to the PATH environment variable. It prevents a
# common mistake where PATH is replaced by only the new <path>.
#
# Example:
#
#    pwd
#    # output: /home/user/my/project
#    PATH_add bin
#    echo $PATH
#    # output: /home/user/my/project/bin:/usr/bin:/bin
#
PATH_add() {
	export PATH="$(expand_path "$1"):$PATH"
}

# Usage: path_add <varname> <path>
#
# Works like PATH_add except that it's for an arbitrary <varname>.
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

# Usage: load_prefix <prefix_path>
#
# Expands some common path variables for the given <prefix_path> prefix. This is
# useful if you installed something in the <prefix_path> using
# $(./configure --prefix=<prefix_path> && make install) and want to use it in 
# the project.
#
# Variables set:
#
#    CPATH
#    LD_LIBRARY_PATH
#    LIBRARY_PATH
#    MANPATH
#    PATH
#    PKG_CONFIG_PATH
#
# Example:
#
#    ./configure --prefix=$HOME/rubies/ruby-1.9.3
#    make && make install
#    # Then in the .envrc
#    load_prefix ~/rubies/ruby-1.9.3
#
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

# Usage: layout <type>
#
# A semantic dispatch used to describe common project layouts.
#
layout() {
	eval "layout_$1"
}

# Usage: layout ruby
#
# Sets the GEM_HOME environment variable to "$PWD/.direnv/ruby/RUBY_VERSION".
# This forces the installation of any gems into the project's sub-folder.
# If you're using bundler it will create wrapper programs that can be invoked
# directly instead of using the $(bundle exec) prefix.
#
layout_ruby() {
	local ruby_version="$(ruby -e"puts (defined?(RUBY_ENGINE) ? RUBY_ENGINE : 'ruby') + '-' + RUBY_VERSION")"

	export GEM_HOME="$PWD/.direnv/${ruby_version}"
	export BUNDLE_BIN="$PWD/.direnv/bin"

	PATH_add ".direnv/${ruby_version}/bin"
	PATH_add ".direnv/bin"
}

# Usage: layout python
#
# Creates and loads a virtualenv environment under "$PWD/.direnv/virtualenv".
# This forces the installation of any egg into the project's sub-folder.
#
layout_python() {
	if ! [ -d .direnv/virtualenv ]; then
		virtualenv --no-site-packages --distribute .direnv/virtualenv
		virtualenv --relocatable .direnv/virtualenv
	fi
	source .direnv/virtualenv/bin/activate
}

# Usage: layout node
#
# Adds "$PWD/node_modules/.bin" to the PATH environment variable.
layout_node() {
	PATH_add node_modules/.bin
}

# Usage: layout go
#
# Sets the GOPATH environment variable to the current directory.
#
layout_go() {
	path_add GOPATH "$PWD"
}

# Usage: use <program_name> [<version>]
#
# A semantic command dispatch intended for loading external dependencies into
# the environment.
#
# Example:
#
#    use_ruby() {
#      echo "Ruby $1"
#    }
#    use ruby 1.9.3
#    # output: Ruby 1.9.3
#
use() {
	local cmd="$1"
	echo "Using $@"
	shift
	use_$cmd "$@"
}

# Usage: use rbenv
#
# Loads rbenv which add the ruby wrappers available on the PATH.
#
use_rbenv() {
	eval "$(rbenv init -)"
}

# Usage: rvm [...]
#
# Should work just like in the shell if you have rvm installed.#
#
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
