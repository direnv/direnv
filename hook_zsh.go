package main

const HOOK_ZSH = "direnv_hook() { eval `direnv private export` };\n" +
	"[[ -z $precmd_functions ]] && precmd_functions=();\n" +
	"precmd_functions=($precmd_functions direnv_hook)"
