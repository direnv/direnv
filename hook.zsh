direnv() {
  if [ `type -w direnv-$1 | cut -d ' ' -f 2` == "function" ]; then
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

direnv_hook() { eval `direnv private export` };
[[ -z $precmd_functions ]] && precmd_functions=();
precmd_functions=($precmd_functions direnv_hook)
