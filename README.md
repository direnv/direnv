SHELL-ENV: Power-up your shell
==============================

shell-env let's you have path-specific environment variables in your
shell.

If you are like me and checkout tons of projects, you don't want to
clutter your .profile or .bashrc. By adding an .envrc to the projects,
you can override some variables while being in that or sub directories.
Keep your .bashrc clean !

Install
-------

Put the shell-env script in your path and add the following lines to
your bashrc:

    precmd() {
      eval `shell-env`
    }
    PROMPT_COMMAND=precmd

It should also work for zsh users but I didn't test it. For you, the
last line is not necessarz.

Usage
-----

Once the shell-env is installed, the script will look for .envrc
in the current and upper directories. If one is found, it will export
the variables to the current shell.

.envrc sample:

    export PATH="$SHELLENV/bin:$PATH"
    export RUBYOPT="-I$SHELLENV/lib"

.envrc is also compatible with rvm:

    eval `rvm --create env ruby-1.9.2@yourproject`

Features
--------

* Adapts with the current path
* Able to revert previous changes

Contribute
----------

Patches are welcome. We can also discuss new features at shellenv@librelist.com.

TODO
----

* It would make sense to port shell-env to BASH or C. There was an initial BASH
version, but I got confused with escaping and ENV diffs.

* shell-env [bash,zsh,ruby,...] for language-specific exports

FAQ
---

Q
: How does RVM update the ENV when changing path?
A
: It overrides cd with :

    cd() {
      builtin cd "$@"
      local result=$?
      __rvm_project_rvmrc
      rvm_hook="after_cd" ; source "$rvm_path/scripts/hook"
      return $result
    }

It does not work in any cases because cd is not the only command to change
directory. (see: pushd for example)

Q
: How does the magic work ?
A
: We set the PROMPT_COMMAND to a function name
  On each prompt display, bash calls the function, adapting the environment
  depending on the path.

Contributors (and thanks)
-------------------------

* Joshua Peek

LICENCE
-------

Copyright (C) 2011 Jonas Pfenniger and contributors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
