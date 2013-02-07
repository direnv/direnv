direnv_hook() { eval `direnv private export` };
[[ -z $precmd_functions ]] && precmd_functions=();
precmd_functions=($precmd_functions direnv_hook)