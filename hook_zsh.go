package main

const HOOK_ZSH = "direnv() {\n" +
	"  if [ `type -w direnv-$1 | cut -d ' ' -f 2` == \"function\" ]; then\n" +
	"    shift $@;\n" +
	"    direnv-$1 \"$@\";\n" +
	"  else\n" +
	"    `which direnv` \"$@\";\n" +
	"  fi\n" +
	"};\n" +
	"direnv-switch() {\n" +
	"  if [ -n \"$DIRENV_BACKUP\" ]; then\n" +
	"    echo \"You need to be in a folder to load a context\"\n" +
	"  else\n" +
	"    export DIRENV_CONTEXT=$1;\n" +
	"  fi\n" +
	"};\n" +
	"\n" +
	"direnv_hook() { eval `direnv private export` };\n" +
	"[[ -z $precmd_functions ]] && precmd_functions=();\n" +
	"precmd_functions=($precmd_functions direnv_hook)\n"
