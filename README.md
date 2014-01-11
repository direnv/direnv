direnv -- Unclutter your .profile
=================================

`direnv` is a shell extension that loads different environment variables
depending on your path.

Instead of putting every environment variable in your "~/.profile", have
directory-specific ".envrc" files for your AWS_ACCESS_KEY, LIBRARY_PATH or
other environment variables.

It does some of the job of rvm, rbenv or virtualenv but in a
language-agnostic way.

## Example

```
$ cd ~/code/my_project
$ ls
bin/ lib/ Rakefile README.md
$ echo $PATH
/usr/local/bin:/usr/local/sbin:/usr/bin:/bin:/usr/sbin:/sbin
$ direnv edit .
# Opens in your `$EDITOR .envrc`. Add:
export PATH=$PWD/bin:$PATH
$
direnv: loading .envrc
direnv export: ~PATH
$ echo $PATH
/Users/zimbatm/code/my_project/bin:/usr/local/bin:/usr/local/sbin:/usr/bin:/bin:/usr/sbin:/sbin
$ cd ..
direnv: unloading
direnv export: ~PATH
$ echo $PATH
/usr/local/bin:/usr/local/sbin:/usr/bin:/bin:/usr/sbin:/sbin
```

## Install

### 1. Install the code

```bash
git clone http://github.com/zimbatm/direnv
cd direnv
make install
# or symlink bin/direnv into your $PATH
```

Installing from the repository requires the Go language.

For Homebrew users, you can also use `brew install direnv`

For MacPorts users, install Go with `sudo port install go`, then install from the repository.

### 2. Add the hook for your shell

This is what is going to enable the direnv extension. It's going to allow
direnv to execute before every prompt command and adjust the environment.

#### BASH

Add the following line at the end of your "~/.bashrc" file:

```bash
eval "$(direnv hook bash)"
```

Make sure it appears even after rvm, git-prompt and other shell extensions
that manipulate your prompt.

#### ZSH

Add the following line at the end of you "~/.zshrc" file:

```bash
eval "$(direnv hook zsh)"
```

If you want to place it in another file replace $0 with "zsh" as zsh changes
the value dynamically.

#### FISH

Add the following line at the end of your "~/.config/fish/config.fish" file:

```
eval (direnv hook fish)
```

## Usage

Use `direnv edit .` to open an ".envrc" in your $EDITOR. This script is going
to be executed once you exit the editor. Every `export` is going to be
available in your shell until you `cd ..` out of the directory.

To make your life convenient there is a couple of additional commands in the
.envrc execution context that are loaded from the `direnv stdlib`.

## Contribute

Bug reports, contributions and forks are welcome.

For bugs, report them on <http://github.com/zimbatm/direnv/issues>

Or if you have some cool usages of direnv that you want to share, feel free
to put them in the wiki <https://github.com/zimbatm/direnv/wiki>

[![Build Status](https://api.travis-ci.org/zimbatm/direnv.png?branch=master)](http://travis-ci.org/zimbatm/direnv)

## COPYRIGHT

Thank you for making direnv better

* Alan Brenner (aka. alanbbr) for the fish shell support
* Alexander Kobel for his patches
* Brian M. Clapper (aka. bmc) for his patch
* Joshua Peek (aka. josh) for his patch and support
* Laurie Young (aka. wildfalcon) for fixing my engrish
* Magnus Holm (aka. judofyr) for his patches and ideas
* Martin Aumüller (aka. aumuell) for his patches and awesomeness
* Peter Waller (aka. pwaller) for his patches and insights
* Sam Stephenson (aka. sstephenson) for his expand_path code that I stole from https://github.com/sstephenson/bats
* Tim Cuthbertson (aka. gfxmonk) for his contribution over the last months

Copyright (C) 2013 Jonas Pfenniger and contributors under the MIT licence.
