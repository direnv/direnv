package main

const STDLIB = "#!bash\n" +
	"#\n" +
	"# These are the commands available in an .envrc context\n" +
	"#\n" +
	"set -e\n" +
	"direnv=\"%s\"\n" +
	"\n" +
	"DIRENV_LOG_FORMAT=\"${DIRENV_LOG_FORMAT-direnv: %%s}\"\n" +
	"\n" +
	"# Usage: log_status [<message> ...]\n" +
	"#\n" +
	"# Logs a status message. Acts like echo,\n" +
	"# but wraps output in the standard direnv log format\n" +
	"# (controlled by $DIRENV_LOG_FORMAT), and directs it\n" +
	"# to stderr rather than stdout.\n" +
	"#\n" +
	"# Example:\n" +
	"#\n" +
	"#    log_status \"Loading ...\"\n" +
	"#\n" +
	"log_status() {\n" +
	"  if [[ -n $DIRENV_LOG_FORMAT ]]; then\n" +
	"    local msg=$*\n" +
	"    # shellcheck disable=SC2059\n" +
	"    printf \"${DIRENV_LOG_FORMAT}\\n\" \"$msg\" >&2\n" +
	"  fi\n" +
	"}\n" +
	"\n" +
	"# Usage: log_error [<message> ...]\n" +
	"#\n" +
	"# Logs an error message. Acts like echo,\n" +
	"# but wraps output in the standard direnv log format\n" +
	"# (controlled by $DIRENV_LOG_FORMAT), and directs it\n" +
	"# to stderr rather than stdout.\n" +
	"#\n" +
	"# Example:\n" +
	"#\n" +
	"#    log_error \"Unable to find specified directory!\"\n" +
	"\n" +
	"log_error() {\n" +
	"  local color_normal\n" +
	"  local color_error\n" +
	"  color_normal=$(tput sgr0)\n" +
	"  color_error=\"\\e[0;31m\"\n" +
	"  if [[ -n $DIRENV_LOG_FORMAT ]]; then\n" +
	"    local msg=$*\n" +
	"    # shellcheck disable=SC2059\n" +
	"    printf \"${color_error}${DIRENV_LOG_FORMAT}${color_normal}\\n\" \"$msg\" >&2\n" +
	"  fi\n" +
	"}\n" +
	"\n" +
	"# Usage: has <command>\n" +
	"#\n" +
	"# Returns 0 if the <command> is available. Returns 1 otherwise. It can be a\n" +
	"# binary in the PATH or a shell function.\n" +
	"#\n" +
	"# Example:\n" +
	"#\n" +
	"#    if has curl; then\n" +
	"#      echo \"Yes we do\"\n" +
	"#    fi\n" +
	"#\n" +
	"has() {\n" +
	"  type \"$1\" &>/dev/null\n" +
	"}\n" +
	"\n" +
	"# Usage: expand_path <rel_path> [<relative_to>]\n" +
	"#\n" +
	"# Outputs the absolute path of <rel_path> relative to <relative_to> or the\n" +
	"# current directory.\n" +
	"#\n" +
	"# Example:\n" +
	"#\n" +
	"#    cd /usr/local/games\n" +
	"#    expand_path ../foo\n" +
	"#    # output: /usr/local/foo\n" +
	"#\n" +
	"expand_path() {\n" +
	"  \"$direnv\" expand_path \"$@\"\n" +
	"}\n" +
	"\n" +
	"# Usage: dotenv [<dotenv>]\n" +
	"#\n" +
	"# Loads a \".env\" file into the current environment\n" +
	"#\n" +
	"dotenv() {\n" +
	"  local path=$1\n" +
	"  if [[ -z $path ]]; then\n" +
	"    path=$PWD/.env\n" +
	"  elif [[ -d $path ]]; then\n" +
	"    path=$path/.env\n" +
	"  fi\n" +
	"  if ! [[ -f $path ]]; then\n" +
	"    log_error \".env at $path not found\"\n" +
	"    return 1\n" +
	"  fi\n" +
	"  eval \"$(\"$direnv\" dotenv bash \"$@\")\"\n" +
	"  watch_file \"$path\"\n" +
	"}\n" +
	"\n" +
	"# Usage: user_rel_path <abs_path>\n" +
	"#\n" +
	"# Transforms an absolute path <abs_path> into a user-relative path if\n" +
	"# possible.\n" +
	"#\n" +
	"# Example:\n" +
	"#\n" +
	"#    echo $HOME\n" +
	"#    # output: /home/user\n" +
	"#    user_rel_path /home/user/my/project\n" +
	"#    # output: ~/my/project\n" +
	"#    user_rel_path /usr/local/lib\n" +
	"#    # output: /usr/local/lib\n" +
	"#\n" +
	"user_rel_path() {\n" +
	"  local abs_path=${1#-}\n" +
	"\n" +
	"  if [[ -z $abs_path ]]; then return; fi\n" +
	"\n" +
	"  if [[ -n $HOME ]]; then\n" +
	"    local rel_path=${abs_path#$HOME}\n" +
	"    if [[ $rel_path != \"$abs_path\" ]]; then\n" +
	"      abs_path=~$rel_path\n" +
	"    fi\n" +
	"  fi\n" +
	"\n" +
	"  echo \"$abs_path\"\n" +
	"}\n" +
	"\n" +
	"# Usage: find_up <filename>\n" +
	"#\n" +
	"# Outputs the path of <filename> when searched from the current directory up to\n" +
	"# /. Returns 1 if the file has not been found.\n" +
	"#\n" +
	"# Example:\n" +
	"#\n" +
	"#    cd /usr/local/my\n" +
	"#    mkdir -p project/foo\n" +
	"#    touch bar\n" +
	"#    cd project/foo\n" +
	"#    find_up bar\n" +
	"#    # output: /usr/local/my/bar\n" +
	"#\n" +
	"find_up() {\n" +
	"  (\n" +
	"    cd \"$(pwd -P 2>/dev/null)\"\n" +
	"    while true; do\n" +
	"      if [[ -f $1 ]]; then\n" +
	"        echo \"$PWD/$1\"\n" +
	"        return 0\n" +
	"      fi\n" +
	"      if [[ $PWD = / ]] || [[ $PWD = // ]]; then\n" +
	"        return 1\n" +
	"      fi\n" +
	"      cd ..\n" +
	"    done\n" +
	"  )\n" +
	"}\n" +
	"\n" +
	"# Usage: source_env <file_or_dir_path>\n" +
	"#\n" +
	"# Loads another \".envrc\" either by specifying its path or filename.\n" +
	"source_env() {\n" +
	"  local rcpath=${1/#\\~/$HOME}\n" +
	"  local rcfile\n" +
	"  if ! [[ -f $rcpath ]]; then\n" +
	"    rcpath=$rcpath/.envrc\n" +
	"  fi\n" +
	"\n" +
	"  rcfile=$(user_rel_path \"$rcpath\")\n" +
	"  watch_file \"$rcpath\"\n" +
	"\n" +
	"  pushd \"$(pwd -P 2>/dev/null)\" > /dev/null\n" +
	"    pushd \"$(dirname \"$rcpath\")\" > /dev/null\n" +
	"    if [[ -f ./$(basename \"$rcpath\") ]]; then\n" +
	"      log_status \"loading $rcfile\"\n" +
	"      # shellcheck source=/dev/null\n" +
	"      . \"./$(basename \"$rcpath\")\"\n" +
	"    else\n" +
	"      log_status \"referenced $rcfile does not exist\"\n" +
	"    fi\n" +
	"    popd > /dev/null\n" +
	"  popd > /dev/null\n" +
	"}\n" +
	"\n" +
	"# Usage: watch_file <filename>\n" +
	"#\n" +
	"# Adds <path> to the list of files that direnv will watch for changes - useful when the contents\n" +
	"# of a file influence how variables are set - especially in direnvrc\n" +
	"#\n" +
	"watch_file() {\n" +
	"  local file=${1/#\\~/$HOME}\n" +
	"\n" +
	"  eval \"$($direnv watch \"$file\")\"\n" +
	"}\n" +
	"\n" +
	"\n" +
	"# Usage: source_up [<filename>]\n" +
	"#\n" +
	"# Loads another \".envrc\" if found with the find_up command.\n" +
	"#\n" +
	"source_up() {\n" +
	"  local file=$1\n" +
	"  local dir\n" +
	"  if [[ -z $file ]]; then\n" +
	"    file=.envrc\n" +
	"  fi\n" +
	"  dir=$(cd .. && find_up \"$file\")\n" +
	"  if [[ -n $dir ]]; then\n" +
	"    source_env \"$(user_rel_path \"$dir\")\"\n" +
	"  fi\n" +
	"}\n" +
	"\n" +
	"# Usage: direnv_load <command-generating-dump-output>\n" +
	"#   e.g: direnv_load opam-env exec -- direnv dump\n" +
	"#\n" +
	"# Applies the environment generated by running <argv> as a\n" +
	"# command. This is useful for adopting the environment of a child\n" +
	"# process - cause that process to run \"direnv dump\" and then wrap\n" +
	"# the results with direnv_load.\n" +
	"#\n" +
	"direnv_load() {\n" +
	"  local exports\n" +
	"  exports=$(\"$direnv\" apply_dump <(\"$@\"))\n" +
	"  local es=$?\n" +
	"  if [[ $es -ne 0 ]]; then\n" +
	"    return $es\n" +
	"  fi\n" +
	"  eval \"$exports\"\n" +
	"}\n" +
	"\n" +
	"# Usage: PATH_add <path>\n" +
	"#\n" +
	"# Prepends the expanded <path> to the PATH environment variable. It prevents a\n" +
	"# common mistake where PATH is replaced by only the new <path>.\n" +
	"#\n" +
	"# Example:\n" +
	"#\n" +
	"#    pwd\n" +
	"#    # output: /home/user/my/project\n" +
	"#    PATH_add bin\n" +
	"#    echo $PATH\n" +
	"#    # output: /home/user/my/project/bin:/usr/bin:/bin\n" +
	"#\n" +
	"PATH_add() {\n" +
	"  PATH=$(expand_path \"$1\"):$PATH\n" +
	"  export PATH\n" +
	"}\n" +
	"\n" +
	"# Usage: path_add <varname> <path>\n" +
	"#\n" +
	"# Works like PATH_add except that it's for an arbitrary <varname>.\n" +
	"path_add() {\n" +
	"  local old_paths=\"${!1}\"\n" +
	"  local dir\n" +
	"  dir=$(expand_path \"$2\")\n" +
	"\n" +
	"  if [[ -z $old_paths ]]; then\n" +
	"    old_paths=\"$dir\"\n" +
	"  else\n" +
	"    old_paths=\"$dir:$old_paths\"\n" +
	"  fi\n" +
	"\n" +
	"  export \"$1=$old_paths\"\n" +
	"}\n" +
	"\n" +
	"# Usage: load_prefix <prefix_path>\n" +
	"#\n" +
	"# Expands some common path variables for the given <prefix_path> prefix. This is\n" +
	"# useful if you installed something in the <prefix_path> using\n" +
	"# $(./configure --prefix=<prefix_path> && make install) and want to use it in\n" +
	"# the project.\n" +
	"#\n" +
	"# Variables set:\n" +
	"#\n" +
	"#    CPATH\n" +
	"#    LD_LIBRARY_PATH\n" +
	"#    LIBRARY_PATH\n" +
	"#    MANPATH\n" +
	"#    PATH\n" +
	"#    PKG_CONFIG_PATH\n" +
	"#\n" +
	"# Example:\n" +
	"#\n" +
	"#    ./configure --prefix=$HOME/rubies/ruby-1.9.3\n" +
	"#    make && make install\n" +
	"#    # Then in the .envrc\n" +
	"#    load_prefix ~/rubies/ruby-1.9.3\n" +
	"#\n" +
	"load_prefix() {\n" +
	"  local dir\n" +
	"  dir=$(expand_path \"$1\")\n" +
	"  path_add CPATH \"$dir/include\"\n" +
	"  path_add LD_LIBRARY_PATH \"$dir/lib\"\n" +
	"  path_add LIBRARY_PATH \"$dir/lib\"\n" +
	"  path_add MANPATH \"$dir/man\"\n" +
	"  path_add MANPATH \"$dir/share/man\"\n" +
	"  path_add PATH \"$dir/bin\"\n" +
	"  path_add PKG_CONFIG_PATH \"$dir/lib/pkgconfig\"\n" +
	"}\n" +
	"\n" +
	"# Usage: layout <type>\n" +
	"#\n" +
	"# A semantic dispatch used to describe common project layouts.\n" +
	"#\n" +
	"layout() {\n" +
	"  local name=$1\n" +
	"  shift\n" +
	"  eval \"layout_$name\" \"$@\"\n" +
	"}\n" +
	"\n" +
	"# Usage: layout go\n" +
	"#\n" +
	"# Sets the GOPATH environment variable to the current directory.\n" +
	"#\n" +
	"layout_go() {\n" +
	"  path_add GOPATH \"$PWD\"\n" +
	"  PATH_add bin\n" +
	"}\n" +
	"\n" +
	"# Usage: layout node\n" +
	"#\n" +
	"# Adds \"$PWD/node_modules/.bin\" to the PATH environment variable.\n" +
	"layout_node() {\n" +
	"  PATH_add node_modules/.bin\n" +
	"}\n" +
	"\n" +
	"# Usage: layout perl\n" +
	"#\n" +
	"# Setup environment variables required by perl's local::lib\n" +
	"# See http://search.cpan.org/dist/local-lib/lib/local/lib.pm for more details\n" +
	"#\n" +
	"layout_perl() {\n" +
	"  local libdir=$PWD/.direnv/perl5\n" +
	"  export LOCAL_LIB_DIR=$libdir\n" +
	"  export PERL_MB_OPT=\"--install_base '$libdir'\"\n" +
	"  export PERL_MM_OPT=\"INSTALL_BASE=$libdir\"\n" +
	"  path_add PERL5LIB \"$libdir/lib/perl5\"\n" +
	"  path_add PERL_LOCAL_LIB_ROOT \"$libdir\"\n" +
	"  PATH_add \"$libdir/bin\"\n" +
	"}\n" +
	"\n" +
	"# Usage: layout python <python_exe>\n" +
	"#\n" +
	"# Creates and loads a virtualenv environment under\n" +
	"# \"$PWD/.direnv/python-$python_version\".\n" +
	"# This forces the installation of any egg into the project's sub-folder.\n" +
	"#\n" +
	"# It's possible to specify the python executable if you want to use different\n" +
	"# versions of python.\n" +
	"#\n" +
	"layout_python() {\n" +
	"  local python=${1:-python}\n" +
	"  local old_env=$PWD/.direnv/virtualenv\n" +
	"  unset PYTHONHOME\n" +
	"  if [[ -d $old_env && $python = python ]]; then\n" +
	"    export VIRTUAL_ENV=$old_env\n" +
	"  else\n" +
	"    local python_version\n" +
	"    python_version=$(\"$python\" -c \"import platform as p;print(p.python_version())\")\n" +
	"    if [[ -z $python_version ]]; then\n" +
	"      log_error \"Could not find python's version\"\n" +
	"      return 1\n" +
	"    fi\n" +
	"\n" +
	"    export VIRTUAL_ENV=$PWD/.direnv/python-$python_version\n" +
	"    if [[ ! -d $VIRTUAL_ENV ]]; then\n" +
	"      virtualenv \"--python=$python\" \"$VIRTUAL_ENV\"\n" +
	"    fi\n" +
	"  fi\n" +
	"  PATH_add \"$VIRTUAL_ENV/bin\"\n" +
	"}\n" +
	"\n" +
	"# Usage: layout python3\n" +
	"#\n" +
	"# A shortcut for $(layout python python3)\n" +
	"#\n" +
	"layout_python3() {\n" +
	"  layout_python python3\n" +
	"}\n" +
	"\n" +
	"# Usage: layout ruby\n" +
	"#\n" +
	"# Sets the GEM_HOME environment variable to \"$PWD/.direnv/ruby/RUBY_VERSION\".\n" +
	"# This forces the installation of any gems into the project's sub-folder.\n" +
	"# If you're using bundler it will create wrapper programs that can be invoked\n" +
	"# directly instead of using the $(bundle exec) prefix.\n" +
	"#\n" +
	"layout_ruby() {\n" +
	"  if ruby -e \"exit Gem::VERSION > '2.2.0'\" 2>/dev/null; then\n" +
	"    export GEM_HOME=$PWD/.direnv/ruby\n" +
	"  else\n" +
	"    local ruby_version\n" +
	"    ruby_version=$(ruby -e\"puts (defined?(RUBY_ENGINE) ? RUBY_ENGINE : 'ruby') + '-' + RUBY_VERSION\")\n" +
	"    export GEM_HOME=$PWD/.direnv/ruby-${ruby_version}\n" +
	"  fi\n" +
	"  export BUNDLE_BIN=$PWD/.direnv/bin\n" +
	"\n" +
	"  PATH_add \"$GEM_HOME/bin\"\n" +
	"  PATH_add \"$BUNDLE_BIN\"\n" +
	"}\n" +
	"\n" +
	"# Usage: use <program_name> [<version>]\n" +
	"#\n" +
	"# A semantic command dispatch intended for loading external dependencies into\n" +
	"# the environment.\n" +
	"#\n" +
	"# Example:\n" +
	"#\n" +
	"#    use_ruby() {\n" +
	"#      echo \"Ruby $1\"\n" +
	"#    }\n" +
	"#    use ruby 1.9.3\n" +
	"#    # output: Ruby 1.9.3\n" +
	"#\n" +
	"use() {\n" +
	"  local cmd=$1\n" +
	"  log_status \"using $*\"\n" +
	"  shift\n" +
	"  \"use_$cmd\" \"$@\"\n" +
	"}\n" +
	"\n" +
	"# Usage: use rbenv\n" +
	"#\n" +
	"# Loads rbenv which add the ruby wrappers available on the PATH.\n" +
	"#\n" +
	"use_rbenv() {\n" +
	"  eval \"$(rbenv init -)\"\n" +
	"}\n" +
	"\n" +
	"# Usage: rvm [...]\n" +
	"#\n" +
	"# Should work just like in the shell if you have rvm installed.#\n" +
	"#\n" +
	"rvm() {\n" +
	"  unset rvm\n" +
	"  if [[ -n ${rvm_scripts_path:-} ]]; then\n" +
	"    # shellcheck source=/dev/null\n" +
	"    source \"${rvm_scripts_path}/rvm\"\n" +
	"  elif [[ -n ${rvm_path:-} ]]; then\n" +
	"    # shellcheck source=/dev/null\n" +
	"    source \"${rvm_path}/scripts/rvm\"\n" +
	"  else\n" +
	"    # shellcheck source=/dev/null\n" +
	"    source \"$HOME/.rvm/scripts/rvm\"\n" +
	"  fi\n" +
	"  rvm \"$@\"\n" +
	"}\n" +
	"\n" +
	"# Usage: use node\n" +
	"# Loads NodeJS version from a `.node-version` or `.nvmrc` file.\n" +
	"#\n" +
	"# Usage: use node <version>\n" +
	"# Loads specified NodeJS version.\n" +
	"#\n" +
	"# If you specify a partial NodeJS version (i.e. `4.2`), a fuzzy match\n" +
	"# is performed and the highest matching version installed is selected.\n" +
	"#\n" +
	"# Environment Variables:\n" +
	"#\n" +
	"# - $NODE_VERSIONS (required)\n" +
	"#   You must specify a path to your installed NodeJS versions via the `$NODE_VERSIONS` variable.\n" +
	"#\n" +
	"# - $NODE_VERSION_PREFIX (optional) [default=\"node-v\"]\n" +
	"#   Overrides the default version prefix.\n" +
	"\n" +
	"use_node() {\n" +
	"  local version=$1\n" +
	"  local via=\"\"\n" +
	"  local node_wanted\n" +
	"  local node_prefix\n" +
	"\n" +
	"  if [[ -z $NODE_VERSIONS ]] || [[ ! -d $NODE_VERSIONS ]]; then\n" +
	"    log_error \"You must specify a \\$NODE_VERSIONS environment variable and the directory specified must exist!\"\n" +
	"    return 1\n" +
	"  fi\n" +
	"\n" +
	"  if [[ -z $version ]] && [[ -f .nvmrc ]]; then\n" +
	"    version=$(< .nvmrc)\n" +
	"    via=\".nvmrc\"\n" +
	"  fi\n" +
	"\n" +
	"  if [[ -z $version ]] && [[ -f .node-version ]]; then\n" +
	"    version=$(< .node-version)\n" +
	"    via=\".node-version\"\n" +
	"  fi\n" +
	"\n" +
	"  if [[ -z $version ]]; then\n" +
	"    log_error \"I do not know which NodeJS version to load because one has not been specified!\"\n" +
	"    return 1\n" +
	"  fi\n" +
	"\n" +
	"  node_wanted=${NODE_VERSION_PREFIX-\"node-v\"}$version\n" +
	"  node_prefix=$(find \"$NODE_VERSIONS\" -maxdepth 1 -mindepth 1 -type d -name \"$node_wanted*\" | sort -r -t . -k 1,1n -k 2,2n -k 3,3n | head -1)\n" +
	"\n" +
	"  if [[ ! -d $node_prefix ]]; then\n" +
	"    log_error \"Unable to find NodeJS version ($version) in ($NODE_VERSIONS)!\"\n" +
	"    return 1\n" +
	"  fi\n" +
	"\n" +
	"  if [[ ! -x $node_prefix/bin/node ]]; then\n" +
	"    log_error \"Unable to load NodeJS binary (node) for version ($version) in ($NODE_VERSIONS)!\"\n" +
	"    return 1\n" +
	"  fi\n" +
	"\n" +
	"  load_prefix \"$node_prefix\"\n" +
	"\n" +
	"  if [[ -z $via ]]; then\n" +
	"    log_status \"Successfully loaded NodeJS $(node --version), from prefix ($node_prefix)\"\n" +
	"  else\n" +
	"    log_status \"Successfully loaded NodeJS $(node --version) (via $via), from prefix ($node_prefix)\"\n" +
	"  fi\n" +
	"}\n" +
	"\n" +
	"# Usage: use_nix [...]\n" +
	"#\n" +
	"# Load environment variables from `nix-shell`.\n" +
	"# If you have a `default.nix` or `shell.nix` these will be\n" +
	"# used by default, but you can also specify packages directly\n" +
	"# (e.g `use nix -p ocaml`).\n" +
	"#\n" +
	"use_nix() {\n" +
	"  direnv_load nix-shell --show-trace \"$@\" --run 'direnv dump'\n" +
	"  if [[ $# = 0 ]]; then\n" +
	"    watch_file default.nix\n" +
	"    watch_file shell.nix\n" +
	"  fi\n" +
	"}\n" +
	"\n" +
	"## Load the global ~/.direnvrc if present\n" +
	"if [[ -f ${XDG_CONFIG_HOME:-$HOME/.config}/direnv/direnvrc ]]; then\n" +
	"  source_env \"${XDG_CONFIG_HOME:-$HOME/.config}/direnv/direnvrc\" >&2\n" +
	"elif [[ -f $HOME/.direnvrc ]]; then\n" +
	"  source_env \"$HOME/.direnvrc\" >&2\n" +
	"fi\n"
