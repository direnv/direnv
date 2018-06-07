package main

const STDLIB = "#!bash\n" +
	"#\n" +
	"# These are the commands available in an .envrc context\n" +
	"#\n" +
	"set -e\n" +
	"# NOTE: don't touch the RHS, it gets replaced at runtime\n" +
	"direnv=\"$(which direnv)\"\n" +
	"\n" +
	"# Config, change in the direnvrc\n" +
	"DIRENV_LOG_FORMAT=\"${DIRENV_LOG_FORMAT-direnv: %s}\"\n" +
	"\n" +
	"# Usage: direnv_layout_dir\n" +
	"#\n" +
	"# Prints the folder path that direnv should use to store layout content.\n" +
	"# This needs to be a function as $PWD might change during source_env/up.\n" +
	"#\n" +
	"# The output defaults to $PWD/.direnv.\n" +
	"\n" +
	"direnv_layout_dir() {\n" +
	"  echo \"${direnv_layout_dir:-$PWD/.direnv}\"\n" +
	"}\n" +
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
	"  color_error=$(tput setaf 1)\n" +
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
	"# Usage: join_args [args...]\n" +
	"#\n" +
	"# Joins all the passed arguments into a single string that can be evaluated by bash\n" +
	"#\n" +
	"# This is useful when one has to serialize an array of arguments back into a string\n" +
	"join_args() {\n" +
	"  printf '%q ' \"$@\"\n" +
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
	"#\n" +
	"# NOTE: the other \".envrc\" is not checked by the security framework.\n" +
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
	"  pushd \"$(pwd 2>/dev/null)\" > /dev/null\n" +
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
	"  eval \"$(\"$direnv\" watch \"$file\")\"\n" +
	"}\n" +
	"\n" +
	"\n" +
	"# Usage: source_up [<filename>]\n" +
	"#\n" +
	"# Loads another \".envrc\" if found with the find_up command.\n" +
	"#\n" +
	"# NOTE: the other \".envrc\" is not checked by the security framework.\n" +
	"source_up() {\n" +
	"  local file=$1\n" +
	"  local dir\n" +
	"  if [[ -z $file ]]; then\n" +
	"    file=.envrc\n" +
	"  fi\n" +
	"  dir=$(cd .. && find_up \"$file\")\n" +
	"  if [[ -n $dir ]]; then\n" +
	"    source_env \"$dir\"\n" +
	"  fi\n" +
	"}\n" +
	"\n" +
	"# Usage: direnv_load <command-generating-dump-output>\n" +
	"#   e.g: direnv_load opam-env exec -- \"$direnv\" dump\n" +
	"#\n" +
	"# Applies the environment generated by running <argv> as a\n" +
	"# command. This is useful for adopting the environment of a child\n" +
	"# process - cause that process to run \"direnv dump\" and then wrap\n" +
	"# the results with direnv_load.\n" +
	"#\n" +
	"direnv_load() {\n" +
	"  local exports\n" +
	"  # backup and restore watches in case of nix-shell --pure\n" +
	"  local __watches=$DIRENV_WATCHES\n" +
	"\n" +
	"  exports=$(\"$direnv\" apply_dump <(\"$@\"))\n" +
	"  local es=$?\n" +
	"  if [[ $es -ne 0 ]]; then\n" +
	"    return $es\n" +
	"  fi\n" +
	"  eval \"$exports\"\n" +
	"\n" +
	"  export DIRENV_WATCHES=$__watches\n" +
	"}\n" +
	"\n" +
	"# Usage: PATH_add <path> [<path> ...]\n" +
	"#\n" +
	"# Prepends the expanded <path> to the PATH environment variable, in order.\n" +
	"# It prevents a common mistake where PATH is replaced by only the new <path>,\n" +
	"# or where a trailing colon is left in PATH, resulting in the current directory\n" +
	"# being considered in the PATH.  Supports adding multiple directories at once.\n" +
	"#\n" +
	"# Example:\n" +
	"#\n" +
	"#    pwd\n" +
	"#    # output: /my/project\n" +
	"#    PATH_add bin\n" +
	"#    echo $PATH\n" +
	"#    # output: /my/project/bin:/usr/bin:/bin\n" +
	"#    PATH_add bam boum\n" +
	"#    echo $PATH\n" +
	"#    # output: /my/project/bam:/my/project/boum:/my/project/bin:/usr/bin:/bin\n" +
	"#\n" +
	"PATH_add() {\n" +
	"  path_add PATH \"$@\"\n" +
	"}\n" +
	"\n" +
	"# Usage: path_add <varname> <path> [<path> ...]\n" +
	"#\n" +
	"# Works like PATH_add except that it's for an arbitrary <varname>.\n" +
	"path_add() {\n" +
	"  local var_name=\"$1\"\n" +
	"  # split existing paths into an array\n" +
	"  declare -a path_array\n" +
	"  IFS=: read -ra path_array <<< \"${!1}\"\n" +
	"  shift\n" +
	"\n" +
	"  # prepend the passed paths in the right order\n" +
	"  for (( i=$# ; i>0 ; i-- )); do\n" +
	"    path_array=( \"$(expand_path \"${!i}\")\" \"${path_array[@]}\" )\n" +
	"  done\n" +
	"\n" +
	"  # join back all the paths\n" +
	"  local path=$(IFS=:; echo \"${path_array[*]}\")\n" +
	"\n" +
	"  # and finally export back the result to the original variable\n" +
	"  export \"$var_name=$path\"\n" +
	"}\n" +
	"\n" +
	"# Usage: MANPATH_add <path>\n" +
	"#\n" +
	"# Prepends a path to the MANPATH environment variable while making sure that\n" +
	"# `man` can still lookup the system manual pages.\n" +
	"#\n" +
	"# If MANPATH is not empty, man will only look in MANPATH.\n" +
	"# So if we set MANPATH=$path, man will only look in $path.\n" +
	"# Instead, prepend to `man -w` (which outputs man's default paths).\n" +
	"#\n" +
	"MANPATH_add() {\n" +
	"  local old_paths=\"${MANPATH:-$(man -w)}\"\n" +
	"  local dir\n" +
	"  dir=$(expand_path \"$1\")\n" +
	"  export \"MANPATH=$dir:$old_paths\"\n" +
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
	"  MANPATH_add \"$dir/man\"\n" +
	"  MANPATH_add \"$dir/share/man\"\n" +
	"  path_add CPATH \"$dir/include\"\n" +
	"  path_add LD_LIBRARY_PATH \"$dir/lib\"\n" +
	"  path_add LIBRARY_PATH \"$dir/lib\"\n" +
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
	"  local libdir=$(direnv_layout_dir)/perl5\n" +
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
	"# \"$direnv_layout_dir/python-$python_version\".\n" +
	"# This forces the installation of any egg into the project's sub-folder.\n" +
	"#\n" +
	"# It's possible to specify the python executable if you want to use different\n" +
	"# versions of python.\n" +
	"#\n" +
	"layout_python() {\n" +
	"  local python=${1:-python}\n" +
	"  [[ $# -gt 0 ]] && shift\n" +
	"  local old_env=$(direnv_layout_dir)/virtualenv\n" +
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
	"    export VIRTUAL_ENV=$(direnv_layout_dir)/python-$python_version\n" +
	"    if [[ ! -d $VIRTUAL_ENV ]]; then\n" +
	"      virtualenv \"--python=$python\" \"$@\" \"$VIRTUAL_ENV\"\n" +
	"    fi\n" +
	"  fi\n" +
	"  PATH_add \"$VIRTUAL_ENV/bin\"\n" +
	"}\n" +
	"\n" +
	"# Usage: layout python2\n" +
	"#\n" +
	"# A shortcut for $(layout python python2)\n" +
	"#\n" +
	"layout_python2() {\n" +
	"  layout_python python2 \"$@\"\n" +
	"}\n" +
	"\n" +
	"# Usage: layout python3\n" +
	"#\n" +
	"# A shortcut for $(layout python python3)\n" +
	"#\n" +
	"layout_python3() {\n" +
	"  layout_python python3 \"$@\"\n" +
	"}\n" +
	"\n" +
	"# Usage: layout anaconda <enviroment_name> [<conda_exe>]\n" +
	"#\n" +
	"# Activates anaconda for the named environment. If the environment\n" +
	"# hasn't been created, it will be using the environment.yml file in\n" +
	"# the current directory. <conda_exe> is optional and will default to\n" +
	"# the one found in the system environment.\n" +
	"#\n" +
	"layout_anaconda() {\n" +
	"  local env_name=$1\n" +
	"  local conda\n" +
	"  if [[ $# -gt 1 ]]; then\n" +
	"    conda=${2}\n" +
	"  else\n" +
	"    conda=$(command -v conda)\n" +
	"  fi\n" +
	"  PATH_add $(dirname \"$conda\")\n" +
	"  local env_loc=$(\"$conda\" env list | grep -- \"$env_name\")\n" +
	"  if [[ ! \"$env_loc\" == $env_name*$env_name ]]; then\n" +
	"    if [[ -e environment.yml ]]; then\n" +
	"      log_status \"creating conda enviroment\"\n" +
	"      \"$conda\" env create\n" +
	"    else\n" +
	"      log_error \"Could not find environment.yml\"\n" +
	"      return 1\n" +
	"    fi\n" +
	"  fi\n" +
	"\n" +
	"  source activate \"$env_name\"\n" +
	"}\n" +
	"\n" +
	"# Usage: layout pipenv\n" +
	"#\n" +
	"# Similar to layout_python, but uses Pipenv to build a\n" +
	"# virtualenv from the Pipfile located in the same directory.\n" +
	"#\n" +
	"layout_pipenv() {\n" +
	"  if [[ ! -f Pipfile ]]; then\n" +
	"    log_error 'No Pipfile found.  Use `pipenv` to create a Pipfile first.'\n" +
	"    exit 2\n" +
	"  fi\n" +
	"\n" +
	"  local VENV=$(pipenv --bare --venv 2>/dev/null)\n" +
	"  if [[ -z $VENV || ! -d $VENV ]]; then\n" +
	"    pipenv install --dev\n" +
	"  fi\n" +
	"\n" +
	"  export VIRTUAL_ENV=$(pipenv --venv)\n" +
	"  PATH_add \"$VIRTUAL_ENV/bin\"\n" +
	"  export PIPENV_ACTIVE=1\n" +
	"}\n" +
	"\n" +
	"# Usage: layout ruby\n" +
	"#\n" +
	"# Sets the GEM_HOME environment variable to \"$(direnv_layout_dir)/ruby/RUBY_VERSION\".\n" +
	"# This forces the installation of any gems into the project's sub-folder.\n" +
	"# If you're using bundler it will create wrapper programs that can be invoked\n" +
	"# directly instead of using the $(bundle exec) prefix.\n" +
	"#\n" +
	"layout_ruby() {\n" +
	"  if ruby -e \"exit Gem::VERSION > '2.2.0'\" 2>/dev/null; then\n" +
	"    export GEM_HOME=$(direnv_layout_dir)/ruby\n" +
	"  else\n" +
	"    local ruby_version\n" +
	"    ruby_version=$(ruby -e\"puts (defined?(RUBY_ENGINE) ? RUBY_ENGINE : 'ruby') + '-' + RUBY_VERSION\")\n" +
	"    export GEM_HOME=$(direnv_layout_dir)/ruby-${ruby_version}\n" +
	"  fi\n" +
	"  export BUNDLE_BIN=$(direnv_layout_dir)/bin\n" +
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
	"  local node_version_prefix=${NODE_VERSION_PREFIX-node-v}\n" +
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
	"  node_wanted=${node_version_prefix}${version}\n" +
	"  node_prefix=$(\n" +
	"    # Look for matching node versions in $NODE_VERSIONS path\n" +
	"    find \"$NODE_VERSIONS\" -maxdepth 1 -mindepth 1 -type d -name \"$node_wanted*\" |\n" +
	"\n" +
	"    # Strip possible \"/\" suffix from $NODE_VERSIONS, then use that to\n" +
	"    # Strip $NODE_VERSIONS/$NODE_VERSION_PREFIX prefix from line.\n" +
	"    while IFS= read -r line; do echo \"${line#${NODE_VERSIONS%/}/${node_version_prefix}}\"; done |\n" +
	"\n" +
	"    # Sort by version: split by \".\" then reverse numeric sort for each piece of the version string\n" +
	"    sort -t . -k 1,1rn -k 2,2rn -k 3,3rn |\n" +
	"\n" +
	"    # The first one is the highest\n" +
	"    head -1\n" +
	"  )\n" +
	"\n" +
	"  node_prefix=\"${NODE_VERSIONS}/${node_version_prefix}${node_prefix}\"\n" +
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
	"  direnv_load nix-shell --show-trace \"$@\" --run \"$(join_args \"$direnv\" dump)\"\n" +
	"  if [[ $# = 0 ]]; then\n" +
	"    watch_file default.nix\n" +
	"    watch_file shell.nix\n" +
	"  fi\n" +
	"}\n" +
	"\n" +
	"# Usage: use_guix [...]\n" +
	"#\n" +
	"# Load environment variables from `guix environment`.\n" +
	"# Any arguments given will be passed to guix environment. For example,\n" +
	"# `use guix hello` would setup an environment with the dependencies of\n" +
	"# the hello package. To create an environment including hello, the\n" +
	"# `--ad-hoc` flag is used `use guix --ad-hoc hello`. Other options\n" +
	"# include `--load` which allows loading an environment from a\n" +
	"# file. For a full list of options, consult the documentation for the\n" +
	"# `guix environment` command.\n" +
	"use_guix() {\n" +
	"  eval \"$(guix environment \"$@\" --search-paths)\"\n" +
	"}\n" +
	"\n" +
	"## Load the global ~/.direnvrc if present\n" +
	"if [[ -f ${XDG_CONFIG_HOME:-$HOME/.config}/direnv/direnvrc ]]; then\n" +
	"  source \"${XDG_CONFIG_HOME:-$HOME/.config}/direnv/direnvrc\" >&2\n" +
	"elif [[ -f $HOME/.direnvrc ]]; then\n" +
	"  source \"$HOME/.direnvrc\" >&2\n" +
	"fi\n" +
	""
