SHELL-ENV: Power-up your shell (mushroom edition)
==============================

shell-env is a really small tool that opens lots of new possibilities.

Usage
-----

Once the shell-env is installed, the script will look for .envrc
in the current and upper directories. If one is found, it will export
the variables to the current shell.

It is also usable by scripts by invoking shell-env.

.envrc is .rvmrc compatible:
eval `rvm --create env ruby-1.9.2@yourproject`

Features
--------

* Adapts with the current path
* Able to revert previous changes

Install
-------

put shell-env in your path, and put the following lines in your .bashrc:

precmd() {
  eval `shell-env`
}
PROMPT_COMMAND=precmd

for zsh, you should normally just need to set the `precmd` function, but
I didn't test it. Please report if it works ! 
(ref: http://www.zsh.org/mla/users/1997/msg00267.html)


TODO
----

* Add filters, there are some env vars line PWD that shouldn't be changed

* It would make sense to port shell-env to BASH or C. There was an initial BASH
version, but I got confused with escaping and ENV diffs.

* Support out-of-envrc variable changes. Currently, if the env is switched,
  user variables are also removed.

* shell-env [bash,zsh,ruby,...] for language-specific exports

Inspirations
------------

Homebrew, RVM

FAQ
---

Q: How does RVM update the ENV when changing path?
A: It overrides cd with :

cd() {
  builtin cd "$@"
  local result=$?
  __rvm_project_rvmrc
  rvm_hook="after_cd" ; source "$rvm_path/scripts/hook"
  return $result
}

It does not work in any cases because cd is not the only command to change
directory. (see: pushd for example)

Q: How does the magic work ?
A: We set the PROMPT_COMMAND to a function name
  On each prompt display, bash calls the function, adapting the environment
  depending on the path.
