direnv() {
  if [ `type -f direnv-$1` == "function" ]; then
    shift $@;
    direnv-$1 "$@";
  else
    `which direnv` "$@";
  fi
};
direnv-switch() {
  if [ -n "$DIRENV_BACKUP" ]; then
    echo "You need to be in a folder to load a context"
  else
    export DIRENV_CONTEXT=$1;
  fi
};

PROMPT_COMMAND="eval \`direnv private export\`;$PROMPT_COMMAND"
