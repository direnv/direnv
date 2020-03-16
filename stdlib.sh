#!/usr/bin/env bash
#
# These are the commands available in an .envrc context
#
# ShellCheck exceptions:
#
# SC1090: Can't follow non-constant source. Use a directive to specify location.
# SC1091: Not following: (file missing)
# SC1117: Backslash is literal in "\n". Prefer explicit escaping: "\\n".
# SC2059: Don't use variables in the printf format string. Use printf "..%s.." "$foo".
shopt -s gnu_errfmt
shopt -s nullglob


# NOTE: don't touch the RHS, it gets replaced at runtime
direnv="$(command -v direnv)"

# Config, change in the direnvrc
DIRENV_LOG_FORMAT="${DIRENV_LOG_FORMAT-direnv: %s}"

# Where direnv configuration should be stored
direnv_config_dir=${XDG_CONFIG_HOME:-$HOME/.config}/direnv

# This variable can be used by programs to detect when they are running inside
# of a .envrc evaluation context. It is ignored by the direnv diffing
# algorithm and so it won't be re-exported.
export DIRENV_IN_ENVRC=1

# Usage: direnv_layout_dir
#
# Prints the folder path that direnv should use to store layout content.
# This needs to be a function as $PWD might change during source_env/up.
#
# The output defaults to $PWD/.direnv.

direnv_layout_dir() {
  echo "${direnv_layout_dir:-$PWD/.direnv}"
}

# Usage: log_status [<message> ...]
#
# Logs a status message. Acts like echo,
# but wraps output in the standard direnv log format
# (controlled by $DIRENV_LOG_FORMAT), and directs it
# to stderr rather than stdout.
#
# Example:
#
#    log_status "Loading ..."
#
log_status() {
  if [[ -n $DIRENV_LOG_FORMAT ]]; then
    local msg=$*
    # shellcheck disable=SC2059,SC1117
    printf "${DIRENV_LOG_FORMAT}\n" "$msg" >&2
  fi
}

# Usage: log_error [<message> ...]
#
# Logs an error message. Acts like echo,
# but wraps output in the standard direnv log format
# (controlled by $DIRENV_LOG_FORMAT), and directs it
# to stderr rather than stdout.
#
# Example:
#
#    log_error "Unable to find specified directory!"

log_error() {
  local color_normal
  local color_error
  color_normal=$(tput sgr0)
  color_error=$(tput setaf 1)
  if [[ -n $DIRENV_LOG_FORMAT ]]; then
    local msg=$*
    # shellcheck disable=SC2059,SC1117
    printf "${color_error}${DIRENV_LOG_FORMAT}${color_normal}\n" "$msg" >&2
  fi
}

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

# Usage: join_args [args...]
#
# Joins all the passed arguments into a single string that can be evaluated by bash
#
# This is useful when one has to serialize an array of arguments back into a string
join_args() {
  printf '%q ' "$@"
}

# Usage: expand_path <rel_path> [<relative_to>]
#
# Outputs the absolute path of <rel_path> relative to <relative_to> or the
# current directory.
#
# Example:
#
#    cd /usr/local/games
#    expand_path ../foo
#    # output: /usr/local/foo
#
expand_path() {
  "$direnv" expand_path "$@"
}

# Usage: dotenv [<dotenv>]
#
# Loads a ".env" file into the current environment
#
dotenv() {
  local path=${1:-}
  if [[ -z $path ]]; then
    path=$PWD/.env
  elif [[ -d $path ]]; then
    path=$path/.env
  fi
  if ! [[ -f $path ]]; then
    log_error ".env at $path not found"
    return 1
  fi
  eval "$("$direnv" dotenv bash "$@")"
  watch_file "$path"
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
  local abs_path=${1#-}

  if [[ -z $abs_path ]]; then return; fi

  if [[ -n $HOME ]]; then
    local rel_path=${abs_path#$HOME}
    if [[ $rel_path != "$abs_path" ]]; then
      abs_path=~$rel_path
    fi
  fi

  echo "$abs_path"
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
    while true; do
      if [[ -f $1 ]]; then
        echo "$PWD/$1"
        return 0
      fi
      if [[ $PWD == / ]] || [[ $PWD == // ]]; then
        return 1
      fi
      cd ..
    done
  )
}

# Usage: source_env <file_or_dir_path>
#
# Loads another ".envrc" either by specifying its path or filename.
#
# NOTE: the other ".envrc" is not checked by the security framework.
source_env() {
  local rcpath=${1/#\~/$HOME}
  local rcfile
  if [[ -d $rcpath ]]; then
    rcpath=$rcpath/.envrc
  fi
  if [[ ! -e $rcpath ]]; then
    log_status "referenced $rcpath does not exist"
    return 1
  fi

  rcfile=$(user_rel_path "$rcpath")
  watch_file "$rcpath"

  pushd "$(pwd 2>/dev/null)" >/dev/null || return 1
  pushd "$(dirname "$rcpath")" >/dev/null || return 1
  if [[ -f ./$(basename "$rcpath") ]]; then
    log_status "loading $rcfile"
    # shellcheck disable=SC1090
    . "./$(basename "$rcpath")"
  else
    log_status "referenced $rcfile does not exist"
  fi
  popd >/dev/null || return 1
  popd >/dev/null || return 1
}

# Usage: source_hash <file_or_dir_path> <shasum>
#
# Loads another ".envrc" either by specifying its path or filename.
# The other ".envrc" is validated using shasum check.
#
source_hash() {
  local rcpath=${1/#\~/$HOME}
  local hash=${2}
  local rcfile
  if [[ -d $rcpath ]]; then
    rcpath=$rcpath/.envrc
  fi
  if [[ ! -e $rcpath ]]; then
    log_status "referenced $rcpath does not exist"
    return 1
  fi
  rchash=$(shasum -t "$rcpath" | cut -d \  -f 1)
  if [[ "$hash" != "$rchash" ]]; then
    log_error "referenced $rcpath has change, please rehash"
    return 1
  fi

  source_env "$rcpath"
}

# Usage: watch_file <filename> [<filename> ...]
#
# Adds each <filename> to the list of files that direnv will watch for changes -
# useful when the contents of a file influence how variables are set -
# especially in direnvrc
#
watch_file() {
  eval "$("$direnv" watch bash "$@")"
}

# Usage: source_up [<filename>]
#
# Loads another ".envrc" if found with the find_up command.
#
# NOTE: the other ".envrc" is not checked by the security framework.
source_up() {
  local dir file=${1:-.envrc}
  dir=$(cd .. && find_up "$file")
  if [[ -n $dir ]]; then
    source_env "$dir"
  fi
}

# Usage: direnv_load <command-generating-dump-output>
#   e.g: direnv_load opam-env exec -- "$direnv" dump
#
# Applies the environment generated by running <argv> as a
# command. This is useful for adopting the environment of a child
# process - cause that process to run "direnv dump" and then wrap
# the results with direnv_load.
#
# shellcheck disable=SC1090
direnv_load() {
  # Backup watches in case of `nix-shell --pure`
  local prev_watches=$DIRENV_WATCHES
  local temp_dir output_file script_file exit_code

  # Prepare a temporary place for dumps and such.
  temp_dir=$(mktemp -dt direnv.XXXXXX) || {
    log_error "Could not create temporary directory."
    return 1
  }
  output_file="$temp_dir/output"
  script_file="$temp_dir/script"

  # Chain the following commands explicitly so that we can capture the exit code
  # of the whole chain. Crucially this ensures that we don't return early (via
  # `set -e`, for example) and hence always remove the temporary directory.
  touch "$output_file" &&
    DIRENV_DUMP_FILE_PATH="$output_file" "$@" &&
    { test -s "$output_file" || {
        log_error "Environment not dumped; did you invoke 'direnv dump'?"
        false
      }
    } &&
    "$direnv" apply_dump "$output_file" > "$script_file" &&
    source "$script_file" ||
      exit_code=$?

  # Scrub temporary directory
  rm -rf "$temp_dir"

  # Restore watches if the dump wiped them
  if [[ -z "${DIRENV_WATCHES:-}" ]]; then
    export DIRENV_WATCHES=$prev_watches
  fi

  # Exit accordingly
  return ${exit_code:-0}
}

# Usage: direnv_apply_dump <file>
#
# Loads the output of `direnv dump` that was stored in a file.
direnv_apply_dump() {
  local path=$1
  eval "$("$direnv" apply_dump "$path")"
}

# Usage: PATH_add <path> [<path> ...]
#
# Prepends the expanded <path> to the PATH environment variable, in order.
# It prevents a common mistake where PATH is replaced by only the new <path>,
# or where a trailing colon is left in PATH, resulting in the current directory
# being considered in the PATH.  Supports adding multiple directories at once.
#
# Example:
#
#    pwd
#    # output: /my/project
#    PATH_add bin
#    echo $PATH
#    # output: /my/project/bin:/usr/bin:/bin
#    PATH_add bam boum
#    echo $PATH
#    # output: /my/project/bam:/my/project/boum:/my/project/bin:/usr/bin:/bin
#
PATH_add() {
  path_add PATH "$@"
}

# Usage: path_add <varname> <path> [<path> ...]
#
# Works like PATH_add except that it's for an arbitrary <varname>.
path_add() {
  local path i var_name="$1"
  # split existing paths into an array
  declare -a path_array
  IFS=: read -ra path_array <<<"${!1}"
  shift

  # prepend the passed paths in the right order
  for ((i = $#; i > 0; i--)); do
    path_array=("$(expand_path "${!i}")" "${path_array[@]}")
  done

  # join back all the paths
  path=$(
    IFS=:
    echo "${path_array[*]}"
  )

  # and finally export back the result to the original variable
  export "$var_name=$path"
}

# Usage: MANPATH_add <path>
#
# Prepends a path to the MANPATH environment variable while making sure that
# `man` can still lookup the system manual pages.
#
# If MANPATH is not empty, man will only look in MANPATH.
# So if we set MANPATH=$path, man will only look in $path.
# Instead, prepend to `man -w` (which outputs man's default paths).
#
MANPATH_add() {
  local old_paths="${MANPATH:-$(man -w)}"
  local dir
  dir=$(expand_path "$1")
  export "MANPATH=$dir:$old_paths"
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
  local dir
  dir=$(expand_path "$1")
  MANPATH_add "$dir/man"
  MANPATH_add "$dir/share/man"
  path_add CPATH "$dir/include"
  path_add LD_LIBRARY_PATH "$dir/lib"
  path_add LIBRARY_PATH "$dir/lib"
  path_add PATH "$dir/bin"
  path_add PKG_CONFIG_PATH "$dir/lib/pkgconfig"
}

# Usage: layout <type>
#
# A semantic dispatch used to describe common project layouts.
#
layout() {
  local name=$1
  shift
  eval "layout_$name" "$@"
}

# Usage: layout go
#
# Sets the GOPATH environment variable to the current directory.
#
layout_go() {
  path_add GOPATH "$PWD"
  PATH_add bin
}

# Usage: layout node
#
# Adds "$PWD/node_modules/.bin" to the PATH environment variable.
layout_node() {
  PATH_add node_modules/.bin
}

# Usage: layout perl
#
# Setup environment variables required by perl's local::lib
# See http://search.cpan.org/dist/local-lib/lib/local/lib.pm for more details
#
layout_perl() {
  local libdir
  libdir=$(direnv_layout_dir)/perl5
  export LOCAL_LIB_DIR=$libdir
  export PERL_MB_OPT="--install_base '$libdir'"
  export PERL_MM_OPT="INSTALL_BASE=$libdir"
  path_add PERL5LIB "$libdir/lib/perl5"
  path_add PERL_LOCAL_LIB_ROOT "$libdir"
  PATH_add "$libdir/bin"
}

# Usage: layout php
#
# Adds "$PWD/vendor/bin" to the PATH environment variable
layout_php() {
  PATH_add vendor/bin
}

# Usage: layout python <python_exe>
#
# Creates and loads a virtual environment under
# "$direnv_layout_dir/python-$python_version".
# This forces the installation of any egg into the project's sub-folder.
# For python older then 3.3 this requires virtualenv to be installed.
#
# It's possible to specify the python executable if you want to use different
# versions of python.
#
layout_python() {
  local old_env
  local python=${1:-python}
  [[ $# -gt 0 ]] && shift
  old_env=$(direnv_layout_dir)/virtualenv
  unset PYTHONHOME
  if [[ -d $old_env && $python == python ]]; then
    VIRTUAL_ENV=$old_env
  else
    local python_version ve
    # shellcheck disable=SC2046
    read -r python_version ve <<<$($python -c "import pkgutil as u, platform as p;ve='venv' if u.find_loader('venv') else ('virtualenv' if u.find_loader('virtualenv') else '');print(p.python_version()+' '+ve)")
    if [[ -z $python_version ]]; then
      log_error "Could not find python's version"
      return 1
    fi

    VIRTUAL_ENV=$(direnv_layout_dir)/python-$python_version
    case $ve in
      "venv")
        if [[ ! -d $VIRTUAL_ENV ]]; then
          $python -m venv "$@" "$VIRTUAL_ENV"
        fi
        ;;
      "virtualenv")
        if [[ ! -d $VIRTUAL_ENV ]]; then
          $python -m virtualenv "$@" "$VIRTUAL_ENV"
        fi
        ;;
      *)
        log_error "Error: neither venv nor virtualenv are available."
        return 1
        ;;
    esac
  fi
  export VIRTUAL_ENV
  PATH_add "$VIRTUAL_ENV/bin"
}

# Usage: layout python2
#
# A shortcut for $(layout python python2)
#
layout_python2() {
  layout_python python2 "$@"
}

# Usage: layout python3
#
# A shortcut for $(layout python python3)
#
layout_python3() {
  layout_python python3 "$@"
}

# Usage: layout anaconda <environment_name> [<conda_exe>]
#
# Activates anaconda for the named environment. If the environment
# hasn't been created, it will be using the environment.yml file in
# the current directory. <conda_exe> is optional and will default to
# the one found in the system environment.
#
layout_anaconda() {
  local env_name=$1
  local env_loc
  local conda
  if [[ $# -gt 1 ]]; then
    conda=${2}
  else
    conda=$(command -v conda)
  fi
  PATH_add "$(dirname "$conda")"
  env_loc=$("$conda" env list | grep -- '^'"$env_name"'\s')
  if [[ ! "$env_loc" == $env_name*$env_name ]]; then
    if [[ -e environment.yml ]]; then
      log_status "creating conda environment"
      "$conda" env create
    else
      log_error "Could not find environment.yml"
      return 1
    fi
  fi

  # shellcheck disable=SC1091
  source activate "$env_name"
}

# Usage: layout pipenv
#
# Similar to layout_python, but uses Pipenv to build a
# virtualenv from the Pipfile located in the same directory.
#
layout_pipenv() {
  PIPENV_PIPFILE="${PIPENV_PIPFILE:-Pipfile}"
  if [[ ! -f "$PIPENV_PIPFILE" ]]; then
    log_error "No Pipfile found.  Use \`pipenv\` to create a \`$PIPENV_PIPFILE\` first."
    exit 2
  fi

  VIRTUAL_ENV=$(pipenv --venv 2>/dev/null ; true)

  if [[ -z $VIRTUAL_ENV || ! -d $VIRTUAL_ENV ]]; then
    pipenv install --dev
    VIRTUAL_ENV=$(pipenv --venv)
  fi

  PATH_add "$VIRTUAL_ENV/bin"
  export PIPENV_ACTIVE=1
  export VIRTUAL_ENV
}

# Usage: layout pyenv <python version number> [<python version number> ...]
#
# Example:
#
#    layout pyenv 3.6.7
#
# Uses pyenv and layout_python to create and load a virtual environment under
# "$direnv_layout_dir/python-$python_version".
#
layout_pyenv() {
  unset PYENV_VERSION
  # layout_python prepends each python version to the PATH, so we add each
  # version in reverse order so that the first listed version ends up
  # first in the path
  local i
  for ((i = $#; i > 0; i--)); do
    local python_version=${!i}
    local pyenv_python
    pyenv_python=$(pyenv root)/versions/${python_version}/bin/python
    if [[ -x "$pyenv_python" ]]; then
      if layout_python "$pyenv_python"; then
        # e.g. Given "use pyenv 3.6.9 2.7.16", PYENV_VERSION becomes "3.6.9:2.7.16"
        PYENV_VERSION=${python_version}${PYENV_VERSION:+:$PYENV_VERSION}
      fi
    else
      log_error "pyenv: version '$python_version' not installed"
      return 1
    fi
  done

  [[ -n "$PYENV_VERSION" ]] && export PYENV_VERSION
}

# Usage: layout ruby
#
# Sets the GEM_HOME environment variable to "$(direnv_layout_dir)/ruby/RUBY_VERSION".
# This forces the installation of any gems into the project's sub-folder.
# If you're using bundler it will create wrapper programs that can be invoked
# directly instead of using the $(bundle exec) prefix.
#
layout_ruby() {
  BUNDLE_BIN=$(direnv_layout_dir)/bin

  if ruby -e "exit Gem::VERSION > '2.2.0'" 2>/dev/null; then
    GEM_HOME=$(direnv_layout_dir)/ruby
  else
    local ruby_version
    ruby_version=$(ruby -e"puts (defined?(RUBY_ENGINE) ? RUBY_ENGINE : 'ruby') + '-' + RUBY_VERSION")
    GEM_HOME=$(direnv_layout_dir)/ruby-${ruby_version}
  fi

  export BUNDLE_BIN
  export GEM_HOME

  PATH_add "$GEM_HOME/bin"
  PATH_add "$BUNDLE_BIN"
}

# Usage: layout julia
#
# Sets the JULIA_PROJECT environment variable to the current directory.
layout_julia() {
  export JULIA_PROJECT=$PWD
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
  local cmd=$1
  log_status "using $*"
  shift
  "use_$cmd" "$@"
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
  if [[ -n ${rvm_scripts_path:-} ]]; then
    # shellcheck disable=SC1090
    source "${rvm_scripts_path}/rvm"
  elif [[ -n ${rvm_path:-} ]]; then
    # shellcheck disable=SC1090
    source "${rvm_path}/scripts/rvm"
  else
    # shellcheck disable=SC1090
    source "$HOME/.rvm/scripts/rvm"
  fi
  rvm "$@"
}

# Usage: use node
# Loads NodeJS version from a `.node-version` or `.nvmrc` file.
#
# Usage: use node [<version>]
# Loads specified NodeJS version.
#
# If you specify a partial NodeJS version (i.e. `4.2`), a fuzzy match
# is performed and the highest matching version installed is selected.
#
# Environment Variables:
#
# - $NODE_VERSIONS (required)
#   You must specify a path to your installed NodeJS versions via the `$NODE_VERSIONS` variable.
#
# - $NODE_VERSION_PREFIX (optional) [default="node-v"]
#   Overrides the default version prefix.

use_node() {
  local version=${1:-}
  local via=""
  local node_version_prefix=${NODE_VERSION_PREFIX-node-v}
  local node_wanted
  local node_prefix

  if [[ -z ${NODE_VERSIONS:-} || ! -d $NODE_VERSIONS ]]; then
    log_error "You must specify a \$NODE_VERSIONS environment variable and the directory specified must exist!"
    return 1
  fi

  if [[ -z $version && -f .nvmrc ]]; then
    version=$(<.nvmrc)
    via=".nvmrc"
  fi

  if [[ -z $version && -f .node-version ]]; then
    version=$(<.node-version)
    via=".node-version"
  fi

  if [[ -z $version ]]; then
    log_error "I do not know which NodeJS version to load because one has not been specified!"
    return 1
  fi

  node_wanted=${node_version_prefix}${version}
  node_prefix=$(
    # Look for matching node versions in $NODE_VERSIONS path
    # Strip possible "/" suffix from $NODE_VERSIONS, then use that to
    # Strip $NODE_VERSIONS/$NODE_VERSION_PREFIX prefix from line.
    # Sort by version: split by "." then reverse numeric sort for each piece of the version string
    # The first one is the highest
    find "$NODE_VERSIONS" -maxdepth 1 -mindepth 1 -type d -name "$node_wanted*" \
      | while IFS= read -r line; do echo "${line#${NODE_VERSIONS%/}/${node_version_prefix}}"; done \
      | sort -t . -k 1,1rn -k 2,2rn -k 3,3rn \
      | head -1
  )

  node_prefix="${NODE_VERSIONS}/${node_version_prefix}${node_prefix}"

  if [[ ! -d $node_prefix ]]; then
    log_error "Unable to find NodeJS version ($version) in ($NODE_VERSIONS)!"
    return 1
  fi

  if [[ ! -x $node_prefix/bin/node ]]; then
    log_error "Unable to load NodeJS binary (node) for version ($version) in ($NODE_VERSIONS)!"
    return 1
  fi

  load_prefix "$node_prefix"

  if [[ -z $via ]]; then
    log_status "Successfully loaded NodeJS $(node --version), from prefix ($node_prefix)"
  else
    log_status "Successfully loaded NodeJS $(node --version) (via $via), from prefix ($node_prefix)"
  fi
}

# Usage: use_nix [...]
#
# Load environment variables from `nix-shell`.
# If you have a `default.nix` or `shell.nix` these will be
# used by default, but you can also specify packages directly
# (e.g `use nix -p ocaml`).
#
use_nix() {
  direnv_load nix-shell --show-trace "$@" --run "$(join_args "$direnv" dump)"
  if [[ $# == 0 ]]; then
    watch_file default.nix shell.nix
  fi
}

# Usage: use_guix [...]
#
# Load environment variables from `guix environment`.
# Any arguments given will be passed to guix environment. For example,
# `use guix hello` would setup an environment with the dependencies of
# the hello package. To create an environment including hello, the
# `--ad-hoc` flag is used `use guix --ad-hoc hello`. Other options
# include `--load` which allows loading an environment from a
# file. For a full list of options, consult the documentation for the
# `guix environment` command.
use_guix() {
  eval "$(guix environment "$@" --search-paths)"
}

# Usage: direnv_version <version_at_least>
#
# Checks that the direnv version is at least old as <version_at_least>.
direnv_version() {
  "$direnv" version "$@"
}

# Usage: __main__ <cmd> [...<args>]
#
# Used by rc.go
__main__() {
  # reserve stdout for dumping
  exec 3>&1
  exec 1>&2

  __dump_at_exit() {
    local ret=$?
    "$direnv" dump json 3
    trap - EXIT
    exit "$ret"
  }
  trap __dump_at_exit EXIT

  # load direnv libraries
  for lib in "$direnv_config_dir/lib/"*.sh; do
    # shellcheck disable=SC1090
    source "$lib"
  done

  # load the global ~/.direnvrc if present
  if [[ -f $direnv_config_dir/direnvrc ]]; then
    # shellcheck disable=SC1090
    source "$direnv_config_dir/direnvrc" >&2
  elif [[ -f $HOME/.direnvrc ]]; then
    # shellcheck disable=SC1090
    source "$HOME/.direnvrc" >&2
  fi

  # and finally load the .envrc
  "$@"
}
