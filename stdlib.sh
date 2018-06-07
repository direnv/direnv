#!bash
#
# These are the commands available in an .envrc context
#
set -e
# NOTE: don't touch the RHS, it gets replaced at runtime
direnv="$(which direnv)"

# Config, change in the direnvrc
DIRENV_LOG_FORMAT="${DIRENV_LOG_FORMAT-direnv: %s}"

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
    # shellcheck disable=SC2059
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
    # shellcheck disable=SC2059
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
  local path=$1
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
      if [[ $PWD = / ]] || [[ $PWD = // ]]; then
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
  if ! [[ -f $rcpath ]]; then
    rcpath=$rcpath/.envrc
  fi

  rcfile=$(user_rel_path "$rcpath")
  watch_file "$rcpath"

  pushd "$(pwd 2>/dev/null)" > /dev/null
    pushd "$(dirname "$rcpath")" > /dev/null
    if [[ -f ./$(basename "$rcpath") ]]; then
      log_status "loading $rcfile"
      # shellcheck source=/dev/null
      . "./$(basename "$rcpath")"
    else
      log_status "referenced $rcfile does not exist"
    fi
    popd > /dev/null
  popd > /dev/null
}

# Usage: watch_file <filename>
#
# Adds <path> to the list of files that direnv will watch for changes - useful when the contents
# of a file influence how variables are set - especially in direnvrc
#
watch_file() {
  local file=${1/#\~/$HOME}

  eval "$("$direnv" watch "$file")"
}


# Usage: source_up [<filename>]
#
# Loads another ".envrc" if found with the find_up command.
#
# NOTE: the other ".envrc" is not checked by the security framework.
source_up() {
  local file=$1
  local dir
  if [[ -z $file ]]; then
    file=.envrc
  fi
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
direnv_load() {
  local exports
  # backup and restore watches in case of nix-shell --pure
  local __watches=$DIRENV_WATCHES

  exports=$("$direnv" apply_dump <("$@"))
  local es=$?
  if [[ $es -ne 0 ]]; then
    return $es
  fi
  eval "$exports"

  export DIRENV_WATCHES=$__watches
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
  local var_name="$1"
  # split existing paths into an array
  declare -a path_array
  IFS=: read -ra path_array <<< "${!1}"
  shift

  # prepend the passed paths in the right order
  for (( i=$# ; i>0 ; i-- )); do
    path_array=( "$(expand_path "${!i}")" "${path_array[@]}" )
  done

  # join back all the paths
  local path=$(IFS=:; echo "${path_array[*]}")

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
  local libdir=$(direnv_layout_dir)/perl5
  export LOCAL_LIB_DIR=$libdir
  export PERL_MB_OPT="--install_base '$libdir'"
  export PERL_MM_OPT="INSTALL_BASE=$libdir"
  path_add PERL5LIB "$libdir/lib/perl5"
  path_add PERL_LOCAL_LIB_ROOT "$libdir"
  PATH_add "$libdir/bin"
}

# Usage: layout python <python_exe>
#
# Creates and loads a virtualenv environment under
# "$direnv_layout_dir/python-$python_version".
# This forces the installation of any egg into the project's sub-folder.
#
# It's possible to specify the python executable if you want to use different
# versions of python.
#
layout_python() {
  local python=${1:-python}
  [[ $# -gt 0 ]] && shift
  local old_env=$(direnv_layout_dir)/virtualenv
  unset PYTHONHOME
  if [[ -d $old_env && $python = python ]]; then
    export VIRTUAL_ENV=$old_env
  else
    local python_version
    python_version=$("$python" -c "import platform as p;print(p.python_version())")
    if [[ -z $python_version ]]; then
      log_error "Could not find python's version"
      return 1
    fi

    export VIRTUAL_ENV=$(direnv_layout_dir)/python-$python_version
    if [[ ! -d $VIRTUAL_ENV ]]; then
      virtualenv "--python=$python" "$@" "$VIRTUAL_ENV"
    fi
  fi
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

# Usage: layout anaconda <enviroment_name> [<conda_exe>]
#
# Activates anaconda for the named environment. If the environment
# hasn't been created, it will be using the environment.yml file in
# the current directory. <conda_exe> is optional and will default to
# the one found in the system environment.
#
layout_anaconda() {
  local env_name=$1
  local conda
  if [[ $# -gt 1 ]]; then
    conda=${2}
  else
    conda=$(command -v conda)
  fi
  PATH_add $(dirname "$conda")
  local env_loc=$("$conda" env list | grep -- "$env_name")
  if [[ ! "$env_loc" == $env_name*$env_name ]]; then
    if [[ -e environment.yml ]]; then
      log_status "creating conda enviroment"
      "$conda" env create
    else
      log_error "Could not find environment.yml"
      return 1
    fi
  fi

  source activate "$env_name"
}

# Usage: layout pipenv
#
# Similar to layout_python, but uses Pipenv to build a
# virtualenv from the Pipfile located in the same directory.
#
layout_pipenv() {
  if [[ ! -f Pipfile ]]; then
    log_error 'No Pipfile found.  Use `pipenv` to create a Pipfile first.'
    exit 2
  fi

  local VENV=$(pipenv --bare --venv 2>/dev/null)
  if [[ -z $VENV || ! -d $VENV ]]; then
    pipenv install --dev
  fi

  export VIRTUAL_ENV=$(pipenv --venv)
  PATH_add "$VIRTUAL_ENV/bin"
  export PIPENV_ACTIVE=1
}

# Usage: layout ruby
#
# Sets the GEM_HOME environment variable to "$(direnv_layout_dir)/ruby/RUBY_VERSION".
# This forces the installation of any gems into the project's sub-folder.
# If you're using bundler it will create wrapper programs that can be invoked
# directly instead of using the $(bundle exec) prefix.
#
layout_ruby() {
  if ruby -e "exit Gem::VERSION > '2.2.0'" 2>/dev/null; then
    export GEM_HOME=$(direnv_layout_dir)/ruby
  else
    local ruby_version
    ruby_version=$(ruby -e"puts (defined?(RUBY_ENGINE) ? RUBY_ENGINE : 'ruby') + '-' + RUBY_VERSION")
    export GEM_HOME=$(direnv_layout_dir)/ruby-${ruby_version}
  fi
  export BUNDLE_BIN=$(direnv_layout_dir)/bin

  PATH_add "$GEM_HOME/bin"
  PATH_add "$BUNDLE_BIN"
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
    # shellcheck source=/dev/null
    source "${rvm_scripts_path}/rvm"
  elif [[ -n ${rvm_path:-} ]]; then
    # shellcheck source=/dev/null
    source "${rvm_path}/scripts/rvm"
  else
    # shellcheck source=/dev/null
    source "$HOME/.rvm/scripts/rvm"
  fi
  rvm "$@"
}

# Usage: use node
# Loads NodeJS version from a `.node-version` or `.nvmrc` file.
#
# Usage: use node <version>
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
  local version=$1
  local via=""
  local node_version_prefix=${NODE_VERSION_PREFIX-node-v}
  local node_wanted
  local node_prefix

  if [[ -z $NODE_VERSIONS ]] || [[ ! -d $NODE_VERSIONS ]]; then
    log_error "You must specify a \$NODE_VERSIONS environment variable and the directory specified must exist!"
    return 1
  fi

  if [[ -z $version ]] && [[ -f .nvmrc ]]; then
    version=$(< .nvmrc)
    via=".nvmrc"
  fi

  if [[ -z $version ]] && [[ -f .node-version ]]; then
    version=$(< .node-version)
    via=".node-version"
  fi

  if [[ -z $version ]]; then
    log_error "I do not know which NodeJS version to load because one has not been specified!"
    return 1
  fi

  node_wanted=${node_version_prefix}${version}
  node_prefix=$(
    # Look for matching node versions in $NODE_VERSIONS path
    find "$NODE_VERSIONS" -maxdepth 1 -mindepth 1 -type d -name "$node_wanted*" |

    # Strip possible "/" suffix from $NODE_VERSIONS, then use that to
    # Strip $NODE_VERSIONS/$NODE_VERSION_PREFIX prefix from line.
    while IFS= read -r line; do echo "${line#${NODE_VERSIONS%/}/${node_version_prefix}}"; done |

    # Sort by version: split by "." then reverse numeric sort for each piece of the version string
    sort -t . -k 1,1rn -k 2,2rn -k 3,3rn |

    # The first one is the highest
    head -1
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
  if [[ $# = 0 ]]; then
    watch_file default.nix
    watch_file shell.nix
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

## Load the global ~/.direnvrc if present
if [[ -f ${XDG_CONFIG_HOME:-$HOME/.config}/direnv/direnvrc ]]; then
  source "${XDG_CONFIG_HOME:-$HOME/.config}/direnv/direnvrc" >&2
elif [[ -f $HOME/.direnvrc ]]; then
  source "$HOME/.direnvrc" >&2
fi
