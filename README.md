direnv - Unclutter your .profile
================================

direnv allows you to have path-dependent environment variables. It works in combination with your shell (bash or zsh) to make the magic work.

Example
-------

    $ cd ~/code/my_project
    $ ls
    bin/ lib/ Rakefile README.md
    $ echo $PATH
    /usr/local/bin:/usr/local/sbin:/usr/bin:/bin:/usr/sbin:/sbin
    $ echo PATH_add bin > .envrc
    direnv: loading /Users/zimbatm/code/my_project
    $ echo $PATH
    /Users/zimbatm/code/my_project/bin:/usr/local/bin:/usr/local/sbin:/usr/bin:/bin:/usr/sbin:/sbin
    $ cd ..
    direnv: unloading /Users/zimbatm/code/my_project
    $ echo $PATH
    /usr/local/bin:/usr/local/sbin:/usr/bin:/bin:/usr/sbin:/sbin

For more examples, check out the wiki pages: https://github.com/zimbatm/direnv/wiki

Install
-------

1a) Get the code and put direnv in your path

    git clone http://github.com/zimbatm/shell-env
    ln -s `pwd`/direnv/bin /usr/local/bin/direnv

1b) Use homebrew (for OSX users)

    brew install direnv

2) Add this line at the end of your .bashrc (after rvm, git-prompt, ...):

    eval `direnv hook $0`


Note that zsh's "named directory" feature will replace %c in your PROMPT with "~DIRENV_DIR". Until I find a solution, use %C instead if it annoys you.

Contributors (and thanks)
-------------------------

* Joshua Peek (aka. josh) for his patch and support
* Magnus Holm (aka. judofyr) for his patches and ideas
* Laurie Young (aka. wildfalcon) for fixing my engrish
* Martin Aumüller (aka. aumuell) for his patches and awesomeness
* Sam Stephenson (aka. sstephenson) for his expand_path code that I stole from https://github.com/sstephenson/bats

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

