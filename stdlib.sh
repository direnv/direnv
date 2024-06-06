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
shopt -s extglob

# NOTE: don't touch the RHS, it gets replaced at runtime
direnv="$(command -v direnv)"

# Config, change in the direnvrc
DIRENV_LOG_FORMAT="${DIRENV_LOG_FORMAT-direnv: %s}"

# Where direnv configuration should be stored
direnv_config_dir="${DIRENV_CONFIG:-${XDG_CONFIG_HOME:-$HOME/.config}/direnv}"

# This variable can be used by programs to detect when they are running inside
# of a .envrc evaluation context. It is ignored by the direnv diffing
# algorithm and so it won't be re-exported.
export DIRENV_IN_ENVRC=1

__env_strictness() {
  local mode tmpfile old_shell_options
  local -i res

  tmpfile="$(mktemp)"
  res=0
  mode="$1"
  shift

  set +o | grep 'pipefail\|nounset\|errexit' >"$tmpfile"
  old_shell_options=$(<"$tmpfile")
  rm -f "$tmpfile"

  case "$mode" in
  strict)
    set -o errexit -o nounset -o pipefail
    ;;
  unstrict)
    set +o errexit +o nounset +o pipefail
    ;;
  *)
    log_error "Unknown strictness mode '${mode}'."
    exit 1
    ;;
  esac

  if (($#)); then
    "${@}"
    res=$?
    eval "$old_shell_options"
  fi

  # Force failure if the inner script has failed and the mode is strict
  if [[ $mode = strict && $res -gt 0 ]]; then
    exit 1
  fi

  return $res
}

# Usage: strict_env [<command> ...]
#
# Turns on shell execution strictness. This will force the .envrc
# evaluation context to exit immediately if:
#
# - any command in a pipeline returns a non-zero exit status that is
#   not otherwise handled as part of `if`, `while`, or `until` tests,
#   return value negation (`!`), or part of a boolean (`&&` or `||`)
#   chain.
# - any variable that has not explicitly been set or declared (with
#   either `declare` or `local`) is referenced.
#
# If followed by a command-line, the strictness applies for the duration
# of the command.
#
# Example:
#
#    strict_env
#    has curl
#
#    strict_env has curl
strict_env() {
  __env_strictness strict "$@"
}

# Usage: unstrict_env [<command> ...]
#
# Turns off shell execution strictness. If followed by a command-line, the
# strictness applies for the duration of the command.
#
# Example:
#
#    unstrict_env
#    has curl
#
#    unstrict_env has curl
unstrict_env() {
  if (($#)); then
    __env_strictness unstrict "$@"
  else
    set +o errexit +o nounset +o pipefail
  fi
}

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
    local msg=$* color_normal=''
    if [[ -t 2 ]]; then
      color_normal="\e[m"
    fi
    # shellcheck disable=SC2059,SC1117
    printf "${color_normal}${DIRENV_LOG_FORMAT}\n" "$msg" >&2
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
  if [[ -n $DIRENV_LOG_FORMAT ]]; then
    local msg=$* color_normal='' color_error=''
    if [[ -t 2 ]]; then
      color_normal="\e[m"
      color_error="\e[38;5;1m"
    fi
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
  local REPLY
  realpath.absolute "${2+"$2"}" "${1+"$1"}"
  echo "$REPLY"
}

# --- vendored from https://github.com/bashup/realpaths
realpath.dirname() {
  REPLY=.
  ! [[ $1 =~ /+[^/]+/*$|^//$ ]] || REPLY="${1%"${BASH_REMATCH[0]}"}"
  REPLY=${REPLY:-/}
}
realpath.basename() {
  REPLY=/
  ! [[ $1 =~ /*([^/]+)/*$ ]] || REPLY="${BASH_REMATCH[1]}"
}

realpath.absolute() {
  REPLY=$PWD
  local eg=extglob
  ! shopt -q $eg || eg=
  ${eg:+shopt -s $eg}
  while (($#)); do case $1 in
    // | //[^/]*)
      REPLY=//
      set -- "${1:2}" "${@:2}"
      ;;
    /*)
      REPLY=/
      set -- "${1##+(/)}" "${@:2}"
      ;;
    */*) set -- "${1%%/*}" "${1##"${1%%/*}"+(/)}" "${@:2}" ;;
    '' | .) shift ;;
    ..)
      realpath.dirname "$REPLY"
      shift
      ;;
    *)
      REPLY="${REPLY%/}/$1"
      shift
      ;;
    esac done
  ${eg:+shopt -u $eg}
}
# ---

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
  watch_file "$path"
  if ! [[ -f $path ]]; then
    log_error ".env at $path not found"
    return 1
  fi
  eval "$("$direnv" dotenv bash "$@")"
}

# Usage: dotenv_if_exists [<filename>]
#
# Loads a ".env" file into the current environment, but only if it exists.
#
dotenv_if_exists() {
  local path=${1:-}
  if [[ -z $path ]]; then
    path=$PWD/.env
  elif [[ -d $path ]]; then
    path=$path/.env
  fi
  watch_file "$path"
  if ! [[ -f $path ]]; then
    return
  fi
  eval "$("$direnv" dotenv bash "$@")"
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
    local rel_path=${abs_path#"$HOME"}
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
  if has cygpath; then
    rcpath=$(cygpath -u "$rcpath")
  fi

  local REPLY
  if [[ -d $rcpath ]]; then
    rcpath=$rcpath/.envrc
  fi
  if [[ ! -e $rcpath ]]; then
    log_status "referenced $rcpath does not exist"
    return 1
  fi

  realpath.dirname "$rcpath"
  local rcpath_dir=$REPLY
  realpath.basename "$rcpath"
  local rcpath_base=$REPLY

  local rcfile
  rcfile=$(user_rel_path "$rcpath")
  watch_file "$rcpath"

  pushd "$(pwd 2>/dev/null)" >/dev/null || return 1
  pushd "$rcpath_dir" >/dev/null || return 1
  if [[ -f ./$rcpath_base ]]; then
    log_status "loading $(user_rel_path "$(expand_path "$rcpath_base")")"
    # shellcheck disable=SC1090
    . "./$rcpath_base"
  else
    log_status "referenced $rcfile does not exist"
  fi
  popd >/dev/null || return 1
  popd >/dev/null || return 1
}

# Usage: source_env_if_exists <filename>
#
# Loads another ".envrc", but only if it exists.
#
# NOTE: contrary to source_env, this only works when passing a path to a file,
#       not a directory.
#
# Example:
#
#    source_env_if_exists .envrc.private
#
source_env_if_exists() {
  watch_file "$1"
  if [[ -f "$1" ]]; then source_env "$1"; fi
}

# Usage: env_vars_required <varname> [<varname> ...]
#
# Logs error for every variable not present in the environment or having an empty value.
# Typically this is used in combination with source_env and source_env_if_exists.
#
# Example:
#
#     # expect .envrc.private to provide tokens
#     source_env .envrc.private
#     # check presence of tokens
#     env_vars_required GITHUB_TOKEN OTHER_TOKEN
#
env_vars_required() {
  local environment
  local -i ret
  environment=$(env)
  ret=0

  for var in "$@"; do
    if [[ "$environment" != *"$var="* || -z ${!var:-} ]]; then
      log_error "env var $var is required but missing/empty"
      ret=1
    fi
  done
  return "$ret"
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

# Usage: watch_dir <dir>
#
# Adds <dir> to the list of dirs that direnv will recursively watch for changes
watch_dir() {
  eval "$("$direnv" watch-dir bash "$1")"
}

# Usage: _source_up [<filename>] [true|false]
#
# Private helper function for source_up and source_up_if_exists. The second
# parameter determines if it's an error for the file we're searching for to
# not exist.
_source_up() {
  local envrc file=${1:-.envrc}
  local ok_if_not_exist=${2}
  envrc=$(cd .. && (find_up "$file" || true))
  if [[ -n $envrc ]]; then
    source_env "$envrc"
  elif $ok_if_not_exist; then
    return 0
  else
    log_error "No ancestor $file found"
    return 1
  fi
}

# Usage: source_up [<filename>]
#
# Loads another ".envrc" if found with the find_up command. Returns 1 if no
# file is found.
#
# NOTE: the other ".envrc" is not checked by the security framework.
source_up() {
  _source_up "${1:-}" false
}

# Usage: source_up_if_exists [<filename>]
#
# Loads another ".envrc" if found with the find_up command. If one is not
# found, nothing happens.
#
# NOTE: the other ".envrc" is not checked by the security framework.
source_up_if_exists() {
  _source_up "${1:-}" true
}

# Usage: fetchurl <url> [<integrity-hash>]
#
# Fetches a URL and outputs a file with its content. If the <integrity-hash>
# is given it will also validate the content of the file before returning it.
fetchurl() {
  "$direnv" fetchurl "$@"
}

# Usage: source_url <url> <integrity-hash>
#
# Fetches a URL and evaluates its content.
source_url() {
  local url=$1 integrity_hash=${2:-} path
  if [[ -z $url ]]; then
    log_error "source_url: <url> argument missing"
    return 1
  fi
  if [[ -z $integrity_hash ]]; then
    log_error "source_url: <integrity-hash> argument missing. Use \`direnv fetchurl $url\` to find out the hash."
    return 1
  fi

  log_status "loading $url ($integrity_hash)"
  path=$(fetchurl "$url" "$integrity_hash")
  # shellcheck disable=SC1090
  source "$path"
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
  local temp_dir output_file script_file exit_code old_direnv_dump_file_path

  # Prepare a temporary place for dumps and such.
  temp_dir=$(mktemp -dt direnv.XXXXXX) || {
    log_error "Could not create temporary directory."
    return 1
  }
  output_file="$temp_dir/output"
  script_file="$temp_dir/script"
  old_direnv_dump_file_path=${DIRENV_DUMP_FILE_PATH:-}

  # Chain the following commands explicitly so that we can capture the exit code
  # of the whole chain. Crucially this ensures that we don't return early (via
  # `set -e`, for example) and hence always remove the temporary directory.
  touch "$output_file" &&
    DIRENV_DUMP_FILE_PATH="$output_file" "$@" &&
    {
      test -s "$output_file" || {
        log_error "Environment not dumped; did you invoke 'direnv dump'?"
        false
      }
    } &&
    "$direnv" apply_dump "$output_file" >"$script_file" &&
    source "$script_file" ||
    exit_code=$?

  # Scrub temporary directory
  rm -rf "$temp_dir"

  # Restore watches if the dump wiped them
  if [[ -z "${DIRENV_WATCHES:-}" ]]; then
    export DIRENV_WATCHES=$prev_watches
  fi

  # Restore DIRENV_DUMP_FILE_PATH if needed
  if [[ -n "$old_direnv_dump_file_path" ]]; then
    export DIRENV_DUMP_FILE_PATH=$old_direnv_dump_file_path
  else
    unset DIRENV_DUMP_FILE_PATH
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
  IFS=: read -ra path_array <<<"${!1-}"
  shift

  # prepend the passed paths in the right order
  for ((i = $#; i > 0; i--)); do
    path_array=("$(expand_path "${!i}")" ${path_array[@]+"${path_array[@]}"})
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

# Usage: PATH_rm <pattern> [<pattern> ...]
# Removes directories that match any of the given shell patterns from
# the PATH environment variable. Order of the remaining directories is
# preserved in the resulting PATH.
#
# Bash pattern syntax:
#   https://www.gnu.org/software/bash/manual/html_node/Pattern-Matching.html
#
# Example:
#
#   echo $PATH
#   # output: /dontremove/me:/remove/me:/usr/local/bin/:...
#   PATH_rm '/remove/*'
#   echo $PATH
#   # output: /dontremove/me:/usr/local/bin/:...
#
PATH_rm() {
  path_rm PATH "$@"
}

# Usage: path_rm <varname> <pattern> [<pattern> ...]
#
# Works like PATH_rm except that it's for an arbitrary <varname>.
path_rm() {
  local path i discard var_name="$1"
  # split existing paths into an array
  declare -a path_array
  IFS=: read -ra path_array <<<"${!1}"
  shift

  patterns=("$@")
  results=()

  # iterate over path entries, discard entries that match any of the patterns
  # shellcheck disable=SC2068
  for path in ${path_array[@]+"${path_array[@]}"}; do
    discard=false
    # shellcheck disable=SC2068
    for pattern in ${patterns[@]+"${patterns[@]}"}; do
      if [[ "$path" == +($pattern) ]]; then
        discard=true
        break
      fi
    done
    if ! $discard; then
      results+=("$path")
    fi
  done

  # join the result paths
  result=$(
    IFS=:
    echo "${results[*]}"
  )

  # and finally export back the result to the original variable
  export "$var_name=$result"
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
  local REPLY
  realpath.absolute "$1"
  MANPATH_add "$REPLY/man"
  MANPATH_add "$REPLY/share/man"
  path_add CPATH "$REPLY/include"
  path_add LD_LIBRARY_PATH "$REPLY/lib"
  path_add LIBRARY_PATH "$REPLY/lib"
  path_add PATH "$REPLY/bin"
  path_add PKG_CONFIG_PATH "$REPLY/lib/pkgconfig"
}

# Usage: semver_search <directory> <folder_prefix> <partial_version>
#
# Search a directory for the highest version number in SemVer format (X.Y.Z).
#
# Examples:
#
# $ tree .
# .
# |-- dir
#     |-- program-1.4.0
#     |-- program-1.4.1
#     |-- program-1.5.0
# $ semver_search "dir" "program-" "1.4.0"
# 1.4.0
# $ semver_search "dir" "program-" "1.4"
# 1.4.1
# $ semver_search "dir" "program-" "1"
# 1.5.0
#
semver_search() {
  local version_dir=${1:-}
  local prefix=${2:-}
  local partial_version=${3:-}
  # Look for matching versions in $version_dir path
  # Strip possible "/" suffix from $version_dir, then use that to
  # strip $version_dir/$prefix prefix from line.
  # Sort by version: split by "." then reverse numeric sort for each piece of the version string
  # The first one is the highest
  find "$version_dir" -maxdepth 1 -mindepth 1 -type d -name "${prefix}${partial_version}*" |
    while IFS= read -r line; do echo "${line#"${version_dir%/}"/"${prefix}"}"; done |
    sort -t . -k 1,1rn -k 2,2rn -k 3,3rn |
    head -1
}

# Usage: layout <type>
#
# A semantic dispatch used to describe common project layouts.
#
layout() {
  local funcname="layout_$1"
  shift
  "$funcname" "$@"
  local layout_dir
  layout_dir=$(direnv_layout_dir)
  if [[ -d "$layout_dir" && ! -f "$layout_dir/CACHEDIR.TAG" ]]; then
    echo 'Signature: 8a477f597d28d172789f06886806bc55
# This file is a cache directory tag created by direnv.
# For information about cache directory tags, see:
#	http://www.brynosaurus.com/cachedir/' >"$layout_dir/CACHEDIR.TAG"
  fi
}

# Usage: layout go
#
# Adds "$(direnv_layout_dir)/go" to the GOPATH environment variable.
# Furthermore "$(direnv_layout_dir)/go/bin" is set as the value for the GOBIN environment variable and added to the PATH environment variable.
layout_go() {
  path_add GOPATH "$(direnv_layout_dir)/go"

  bindir="$(direnv_layout_dir)/go/bin"
  PATH_add "$bindir"
  export GOBIN="$bindir"
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
# Creates and loads a virtual environment.
# You can specify the path of the virtual environment through VIRTUAL_ENV
# environment variable, otherwise it will be set to
# "$direnv_layout_dir/python-$python_version".
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
    read -r python_version ve <<<$($python -c "import importlib.util as u, platform as p;ve='venv' if u.find_spec('venv') else ('virtualenv' if u.find_spec('virtualenv') else '');print('.'.join(p.python_version_tuple()[:2])+' '+ve)")
    if [[ -z $python_version ]]; then
      log_error "Could not find python's version"
      return 1
    fi

    if [[ -n "${VIRTUAL_ENV:-}" ]]; then
      local REPLY
      realpath.absolute "$VIRTUAL_ENV"
      VIRTUAL_ENV=$REPLY
    else
      VIRTUAL_ENV=$(direnv_layout_dir)/python-$python_version
    fi
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

# Usage: layout anaconda <env_spec> [<conda_exe>]
#
# Activates anaconda for the provided environment.
# The <env_spec> can be one of the following:
#   1. Name of an environment
#   2. Prefix path to an environment
#   3. Path to a yml-formatted file specifying the environment
#
# Environment creation will use environment.yml, if
# available, when a name or prefix is provided. Otherwise,
# an empty environment will be created.
#
# <conda_exe> is optional and will default to the one
# found in the system environment.
#
layout_anaconda() {
  local env_spec=$1
  local env_name
  local env_loc
  local env_config
  local conda
  local REPLY
  if [[ $# -gt 1 ]]; then
    conda=${2}
  else
    conda=$(command -v conda)
  fi
  realpath.dirname "$conda"
  PATH_add "$REPLY"

  if [[ "${env_spec##*.}" == "yml" ]]; then
    env_config=$env_spec
  elif [[ "${env_spec%%/*}" == "." ]]; then
    # "./foo" relative prefix
    realpath.absolute "$env_spec"
    env_loc="$REPLY"
  elif [[ ! "$env_spec" == "${env_spec#/}" ]]; then
    # "/foo" absolute prefix
    env_loc="$env_spec"
  elif [[ -n "$env_spec" ]]; then
    # "name" specified
    env_name="$env_spec"
  else
    # Need at least one
    env_config=environment.yml
  fi

  # If only config, it needs a name field
  if [[ -n "$env_config" ]]; then
    if [[ -e "$env_config" ]]; then
      env_name="$(grep -- '^name:' "$env_config")"
      env_name="${env_name/#name:*([[:space:]])/}"
      if [[ -z "$env_name" ]]; then
        log_error "Unable to find 'name' in '$env_config'"
        return 1
      fi
    else
      log_error "Unable to find config '$env_config'"
      return 1
    fi
  fi

  # Try to find location based on name
  if [[ -z "$env_loc" ]]; then
    # Update location if already created
    env_loc=$("$conda" env list | grep -- '^'"$env_name"'\s')
    env_loc="${env_loc##* }"
  fi

  # Check for environment existence
  if [[ ! -d "$env_loc" ]]; then

    # Create if necessary
    if [[ -z "$env_config" ]] && [[ -n "$env_name" ]]; then
      if [[ -e environment.yml ]]; then
        "$conda" env create --file environment.yml --name "$env_name"
      else
        "$conda" create -y --name "$env_name"
      fi
    elif [[ -n "$env_config" ]]; then
      "$conda" env create --file "$env_config"
    elif [[ -n "$env_loc" ]]; then
      if [[ -e environment.yml ]]; then
        "$conda" env create --file environment.yml --prefix "$env_loc"
      else
        "$conda" create -y --prefix "$env_loc"
      fi
    fi

    if [[ -z "$env_loc" ]]; then
      # Update location if already created
      env_loc=$("$conda" env list | grep -- '^'"$env_name"'\s')
      env_loc="${env_loc##* }"
    fi
  fi

  eval "$("$conda" shell.bash activate "$env_loc")"
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

  VIRTUAL_ENV=$(
    pipenv --venv 2>/dev/null
    true
  )

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
# Uses pyenv and layout_python to create and load a virtual environment.
# You can specify the path of the virtual environment through VIRTUAL_ENV
# environment variable, otherwise it will be set to
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

# Usage: use julia [<version>]
# Loads specified Julia version.
#
# Environment Variables:
#
# - $JULIA_VERSIONS (required)
#   You must specify a path to your installed Julia versions with the `$JULIA_VERSIONS` variable.
#
# - $JULIA_VERSION_PREFIX (optional) [default="julia-"]
#   Overrides the default version prefix.
#
use_julia() {
  local version=${1:-}
  local julia_version_prefix=${JULIA_VERSION_PREFIX-julia-}
  local search_version
  local julia_prefix

  if [[ -z ${JULIA_VERSIONS:-} || -z $version ]]; then
    log_error "Must specify the \$JULIA_VERSIONS environment variable and a Julia version!"
    return 1
  fi

  julia_prefix="${JULIA_VERSIONS}/${julia_version_prefix}${version}"

  if [[ ! -d ${julia_prefix} ]]; then
    search_version=$(semver_search "${JULIA_VERSIONS}" "${julia_version_prefix}" "${version}")
    julia_prefix="${JULIA_VERSIONS}/${julia_version_prefix}${search_version}"
  fi

  if [[ ! -d $julia_prefix ]]; then
    log_error "Unable to find Julia version ($version) in ($JULIA_VERSIONS)!"
    return 1
  fi

  if [[ ! -x $julia_prefix/bin/julia ]]; then
    log_error "Unable to load Julia binary (julia) for version ($version) in ($JULIA_VERSIONS)!"
    return 1
  fi

  PATH_add "$julia_prefix/bin"
  MANPATH_add "$julia_prefix/share/man"

  log_status "Successfully loaded $(julia --version), from prefix ($julia_prefix)"
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
    # shellcheck disable=SC1090,SC1091
    source "${rvm_scripts_path}/rvm"
  elif [[ -n ${rvm_path:-} ]]; then
    # shellcheck disable=SC1090,SC1091
    source "${rvm_path}/scripts/rvm"
  else
    # shellcheck disable=SC1090,SC1091
    source "$HOME/.rvm/scripts/rvm"
  fi
  rvm "$@"
}

# Usage: use node [<version>]
#
# Loads the specified NodeJS version into the environment.
#
# If a partial NodeJS version is passed (i.e. `4.2`), a fuzzy match
# is performed and the highest matching version installed is selected.
#
# If no version is passed, it will look at the '.nvmrc' or '.node-version'
# files in the current directory if they exist.
#
# Environment Variables:
#
# - $NODE_VERSIONS (required)
#   Points to a folder that contains all the installed Node versions. That
#   folder must exist.
#
# - $NODE_VERSION_PREFIX (optional) [default="node-v"]
#   Overrides the default version prefix.
#
use_node() {
  local version=${1:-}
  local via=""
  local node_version_prefix=${NODE_VERSION_PREFIX-node-v}
  local search_version
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

  version=${version#v}

  if [[ -z $version ]]; then
    log_error "I do not know which NodeJS version to load because one has not been specified!"
    return 1
  fi

  # Search for the highest version matching $version in the folder
  search_version=$(semver_search "$NODE_VERSIONS" "${node_version_prefix}" "${version}")
  node_prefix="${NODE_VERSIONS}/${node_version_prefix}${search_version}"

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

# Usage: use nodenv <node version number>
#
# Example:
#
#    use nodenv 15.2.1
#
# Uses nodenv, use_node and layout_node to add the chosen node version and
# "$PWD/node_modules/.bin" to the PATH
#
use_nodenv() {
  local node_version="${1}"
  local node_versions_dir
  local nodenv_version
  node_versions_dir="$(nodenv root)/versions"
  nodenv_version="${node_versions_dir}/${node_version}"
  if [[ -e "$nodenv_version" ]]; then
    # Put the selected node version in the PATH
    NODE_VERSIONS="${node_versions_dir}" NODE_VERSION_PREFIX="" use_node "${node_version}"
    # Add $PWD/node_modules/.bin to the PATH
    layout_node
  else
    log_error "nodenv: version '$node_version' not installed.  Use \`nodenv install ${node_version}\` to install it first."
    return 1
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

# Usage: use_flake [<installable>]
#
# Load the build environment of a derivation similar to `nix develop`.
#
# By default it will load the current folder flake.nix devShell attribute. Or
# pass an "installable" like "nixpkgs#hello" to load all the build
# dependencies of the hello package from the latest nixpkgs.
#
# Note that the flakes feature is hidden behind an experimental flag, which
# you will have to enable on your own. Flakes is not considered stable yet.
use_flake() {
  watch_file flake.nix
  watch_file flake.lock
  mkdir -p "$(direnv_layout_dir)"
  eval "$(nix print-dev-env --profile "$(direnv_layout_dir)/flake-profile" "$@")"
  nix profile wipe-history --profile "$(direnv_layout_dir)/flake-profile"
}

# Usage: use_guix [...]
#
# Load environment variables from `guix shell`.
# Any arguments given will be passed to guix shell. For example,
# `use guix hello` would setup an environment including the hello
# package. To create an environment with the hello dependencies, the
# `--development` flag is used `use guix --development hello`. Other
# options include `--file` which allows loading an environment from a
# file. For a full list of options, consult the documentation for the
# `guix shell` command.
use_guix() {
  eval "$(guix shell "$@" --search-paths)"
}

# Usage: use_vim [<vimrc_file>]
#
# Prepends the specified vim script (or .vimrc.local by default) to the
# `DIRENV_EXTRA_VIMRC` environment variable.
#
# This variable is understood by the direnv/direnv.vim extension. When found,
# it will source it after opening files in the directory.
use_vim() {
  local extra_vimrc=${1:-.vimrc.local}
  path_add DIRENV_EXTRA_VIMRC "$extra_vimrc"
}

# Usage: direnv_version <version_at_least>
#
# Checks that the direnv version is at least old as <version_at_least>.
direnv_version() {
  "$direnv" version "$@"
}

# Usage: on_git_branch [<branch_name>] OR on_git_branch -r [<regexp>]
#
# Returns 0 if within a git repository with given `branch_name`. If no branch
# name is provided, then returns 0 when within _any_ branch. Requires the git
# command to be installed. Returns 1 otherwise.
#
# When the `-r` flag is specified, then the second argument is interpreted as a
# regexp pattern for matching a branch name.
#
# Regardless, when a branch is specified, then `.git/HEAD` is watched so that
# entering/exiting a branch triggers a reload.
#
# Example (.envrc):
#
#    if on_git_branch; then
#      echo "Thanks for contributing to a GitHub project!"
#    fi
#
#    if on_git_branch child_changes; then
#      export MERGE_BASE_BRANCH=parent_changes
#    fi
#
#    if on_git_branch -r '.*py2'; then
#      layout python2
#    else
#      layout python
#    fi
on_git_branch() {
  local git_dir
  if ! has git; then
    log_error "on_git_branch needs git, which could not be found on your system"
    return 1
  elif ! git_dir=$(git rev-parse --absolute-git-dir 2>/dev/null); then
    log_error "on_git_branch could not locate the .git directory corresponding to the current working directory"
    return 1
  elif [ -z "$1" ]; then
    return 0
  elif [[ "$1" = "-r" && -z "$2" ]]; then
    log_error "missing regexp pattern after \`-r\` flag"
    return 1
  fi
  watch_file "$git_dir/HEAD"
  local git_branch
  git_branch=$(git branch --show-current)
  if [ "$1" = '-r' ]; then
    [[ "$git_branch" =~ $2 ]]
  else
    [ "$1" = "$git_branch" ]
  fi
}

# Usage: __main__ <cmd> [...<args>]
#
# Used by rc.go
__main__() {
  # reserve stdout for dumping
  exec 3>&1
  exec 1>&2

  # shellcheck disable=SC2317
  __dump_at_exit() {
    local ret=$?
    "$direnv" dump json "" >&3
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
    # shellcheck disable=SC1090,SC1091
    source "$direnv_config_dir/direnvrc" >&2
  elif [[ -f $HOME/.direnvrc ]]; then
    # shellcheck disable=SC1090,SC1091
    source "$HOME/.direnvrc" >&2
  fi

  # and finally load the .envrc
  "$@"
}
