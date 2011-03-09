direnv - Unclutter your .profile
================================

direnv is a small POSIX utility that works in combination with your shell (bash or zsh).
It allows you to have path-dependent environment variables, to load and unload them
when navigating trough your filesystem.

Usage
-----

Usually, I set an .envrc file in projects that have a bin/ or lib/ directory.

When navigating in the project's root direcory or it's children, direnv will
add the bin/ directory to my path and the lib/ directory to the target-language's libpath.
That way they are handily available, no need to set something by hand.

I also use direnv to set the GEM_HOME environment variable to make project-dependent gemsets.
The gemset moves with the project unlike in rvm. Feature request: the cache directory could
be set to $HOME/.gem/cache, it would really help.

Install
-------

1) Get the code

  git clone http://github.com/zimbatm/shell-env

2) Put `direnv` in your PATH

For example by symlinking it in your ~/bin or /usr/local/bin directory

3) Add those lines to your .bashrc:

    precmd() {
      eval `direnv export`
    }
    PROMPT_COMMAND=precmd

zsh users can use the same code and forget the last line, precmd is a magic function. (TO BE TESTED, please report)

Contributors (and thanks)
-------------------------

* Joshua Peek (aka. josh) for his patch
* Magnus Holm (aka. judofyr) for his patches and ideas

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

