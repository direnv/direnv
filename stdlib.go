package main

const STDLIB = "# These are the commands available in an .envrc context\n" +
	"\n" +
	"set -e\n" +
	"\n" +
	"# Usage: has something\n" +
	"# determines if \"something\" is availabe as a command\n" +
	"has() {\n" +
	"  type \"$1\" &>/dev/null\n" +
	"}\n" +
	"\n" +
	"# Usage: expand_path ./rel/path [RELATIVE_TO]\n" +
	"# RELATIVE_TO is $PWD by default\n" +
	"expand_path() {\n" +
	"  $DIRENV_LIBEXEC/direnv private expand_path \"$@\"\n" +
	"}\n" +
	"\n" +
	"# Usage: user_rel_path /Users/you/some_path => ~/some_path\n" +
	"user_rel_path() {\n" +
	"  local path=${1#-}\n" +
	"\n" +
	"  if [ -z \"$path\" ]; then return; fi\n" +
	"\n" +
	"  if [ -n \"$HOME\" ]; then\n" +
	"    local rel_path=\"${path#$HOME}\"\n" +
	"    if [ \"$rel_path\" != \"$path\" ]; then\n" +
	"      path=\"~${rel_path}\"\n" +
	"    fi\n" +
	"  fi\n" +
	"\n" +
	"  echo $path\n" +
	"}\n" +
	"\n" +
	"# Usage: find_up some_file\n" +
	"find_up() {\n" +
	"  (\n" +
	"    cd \"`pwd -P 2>/dev/null`\"\n" +
	"    while true; do\n" +
	"      if [ -f \"$1\" ]; then\n" +
	"        echo $PWD/$1\n" +
	"        return 0\n" +
	"      fi\n" +
	"      if [ \"$PWD\" = \"/\" ] || [ \"$PWD\" = \"//\" ]; then\n" +
	"        return 1\n" +
	"      fi\n" +
	"      cd ..\n" +
	"    done\n" +
	"  )\n" +
	"}\n" +
	"\n" +
	"direnv_find_rc() {\n" +
	"  local path=`find_up .envrc`\n" +
	"  if [ -n \"$path\" ]; then\n" +
	"    cd \"$(dirname \"$path\")\"\n" +
	"    return 0\n" +
	"  else\n" +
	"    return 1\n" +
	"  fi\n" +
	"}\n" +
	"\n" +
	"# Safer PATH handling\n" +
	"#\n" +
	"# Usage: PATH_add PATH\n" +
	"# Example: PATH_add bin\n" +
	"PATH_add() {\n" +
	"  export PATH=\"`expand_path \"$1\"`:$PATH\"\n" +
	"}\n" +
	"\n" +
	"# Safer path handling\n" +
	"#\n" +
	"# Usage: path_add VARNAME PATH\n" +
	"# Example: path_add LD_LIBRARY_PATH ./lib\n" +
	"path_add() {\n" +
	"  local old_paths=${!1}\n" +
	"  local path=`expand_path \"$2\"`\n" +
	"\n" +
	"  if [ -z \"$old_paths\" ]; then\n" +
	"    old_paths=\"$path\"\n" +
	"  else\n" +
	"    old_paths=\"$path:$old_paths\"\n" +
	"  fi\n" +
	"\n" +
	"  export $1=\"$old_paths\"\n" +
	"}\n" +
	"\n" +
	"# Usage: layout ruby\n" +
	"layout_ruby() {\n" +
	"  # TODO: ruby_version should be the ABI version\n" +
	"  local ruby_version=`ruby -e\"puts (defined?(RUBY_ENGINE) ? RUBY_ENGINE : 'ruby') + '-' + RUBY_VERSION\"`\n" +
	"\n" +
	"  export GEM_HOME=$PWD/.direnv/${ruby_version}\n" +
	"  export BUNDLE_BIN=$PWD/.direnv/bin\n" +
	"\n" +
	"  PATH_add \".direnv/${ruby_version}/bin\"\n" +
	"  PATH_add \".direnv/bin\"\n" +
	"}\n" +
	"\n" +
	"layout_python() {\n" +
	"  if ! [ -d .direnv/virtualenv ]; then\n" +
	"    virtualenv --no-site-packages --distribute .direnv/virtualenv\n" +
	"    virtualenv --relocatable .direnv/virtualenv\n" +
	"  fi\n" +
	"  source .direnv/virtualenv/bin/activate\n" +
	"}\n" +
	"\n" +
	"layout_node() {\n" +
	"  PATH_add node_modules/.bin\n" +
	"}\n" +
	"\n" +
	"layout() {\n" +
	"  eval \"layout_$1\"\n" +
	"}\n" +
	"\n" +
	"# This folder contains a <program-name>/<version> structure\n" +
	"cellar_path=/usr/local/Cellar\n" +
	"set_cellar_path() {\n" +
	"  cellar_path=$1\n" +
	"}\n" +
	"\n" +
	"# Usage: use PROGRAM_NAME VERSION\n" +
	"# Example: use ruby 1.9.3\n" +
	"use() {\n" +
	"  if has use_$1 ; then\n" +
	"    echo \"Using $1 v$2\"\n" +
	"    eval \"use_$1 $2\"\n" +
	"    return $?\n" +
	"  fi\n" +
	"  \n" +
	"  local path=\"$cellar_path/$1/$2/bin\"\n" +
	"  if [ -d \"$path\" ]; then\n" +
	"    echo \"Using $1 v$2\"\n" +
	"    PATH_add \"$path\"\n" +
	"    return\n" +
	"  fi\n" +
	"    \n" +
	"  echo \"* Unable to load $path\"\n" +
	"  return 1\n" +
	"}\n" +
	"\n" +
	"# Inherit another .envrc\n" +
	"# Usage: source_env <FILE_OR_DIR_PATH>\n" +
	"source_env() {\n" +
	"  local rcfile=\"$1\"\n" +
	"  if ! [ -f \"$1\" ]; then\n" +
	"    rcfile=\"$rcfile/.envrc\"\n" +
	"  fi\n" +
	"  echo \"direnv: loading $(user_rel_path \"$rcfile\")\"\n" +
	"  pushd \"`dirname \"$rcfile\"`\" > /dev/null\n" +
	"  set +u\n" +
	"  . \"./`basename \"$rcfile\"`\"\n" +
	"  popd > /dev/null\n" +
	"}\n" +
	"\n" +
	"# Inherits the first .envrc (or given FILE_NAME) it finds in the path\n" +
	"# Usage: source_up [FILE_NAME]\n" +
	"source_up() {\n" +
	"  local file=\"$1\"\n" +
	"  if [ -z \"$file\" ]; then\n" +
	"    file=\".envrc\"\n" +
	"  fi\n" +
	"  local path=`cd .. && find_up \"$file\"`\n" +
	"  if [ -n \"$path\" ]; then\n" +
	"    source_env \"$path\"\n" +
	"  fi\n" +
	"}\n" +
	"\n" +
	"if [ -n \"${rvm_path-}\" ]; then\n" +
	"  # source rvm on first call\n" +
	"  rvm() {\n" +
	"    unset rvm\n" +
	"    set +e\n" +
	"    . \"$rvm_path/scripts/rvm\"\n" +
	"    rvm $@\n" +
	"    set -e\n" +
	"  }\n" +
	"fi\n"
