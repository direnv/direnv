package main

const HOOK_BASH = "direnv() {\n" +
	"  if [ `type -f direnv-$1` == \"function\" ]; then\n" +
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
	"PROMPT_COMMAND=\"eval \\`direnv private export\\`;$PROMPT_COMMAND\"\n"
