direnv -- Unclutter your .profile
=================================

`direnv` is an environment switcher for the shell. It knows how to hook into
bash, zsh, tcsh, fish shell and elvish to load or unload environment variables
depending on the current directory. This allows project-specific
environment variables without cluttering the `~/.profile` file.

Before each prompt, direnv checks for the existence of a ".envrc" file in the
current and parent directories. If the file exists (and is authorized), it is
loaded into a **bash** sub-shell and all exported variables are then captured by
direnv and then made available to the current shell.

Because direnv is compiled into a single static executable, it is fast enough
to be unnoticeable on each prompt. It is also language-agnostic and can be
used to build solutions similar to rbenv, pyenv and phpenv.


## Example

```
$ cd ~/my_project
$ echo ${FOO-nope}
nope
$ echo export FOO=foo > .envrc
.envrc is not allowed
$ direnv allow .
direnv: reloading
direnv: loading .envrc
direnv export: +FOO
$ echo ${FOO-nope}
foo
$ cd ..
direnv: unloading
direnv export: ~PATH
$ echo ${FOO-nope}
nope
```

## Install

### From system packages

direnv is packaged for a variety of systems:

* [Fedora](https://apps.fedoraproject.org/packages/direnv)
* [Arch AUR](https://aur.archlinux.org/packages/direnv/)
* [Debian](https://packages.debian.org/search?keywords=direnv&searchon=names&suite=all&section=all)
* [Gentoo go-overlay](https://github.com/Dr-Terrible/go-overlay)
* [NetBSD pkgsrc-wip](http://www.pkgsrc.org/wip/)
* [NixOS](https://nixos.org/nixos/packages.html)
* [OSX Homebrew](http://brew.sh/)
* [OpenSuSe](https://build.opensuse.org/package/show/openSUSE%3AFactory/direnv)
* [MacPorts](https://www.macports.org/)
* [Ubuntu](https://packages.ubuntu.com/search?keywords=direnv&searchon=names&suite=all&section=all)
* [GNU Guix](https://www.gnu.org/software/guix/)

Non-official builds for Debian and Fedora are also available at
https://dl.equinox.io/zimbatm/direnv/stable

### From binary builds

Binary builds for a variety of architectures are also available for
[each release](https://github.com/direnv/direnv/releases).

Fetch the binary, `chmod +x direnv` and put it somewhere in your PATH.

### Compile from source

Setup go environment https://golang.org/doc/install

> go >= 1.9 is required

For example using $HOME/go for your workspace

    $ export GOPATH=$HOME/go

Create the directory:

    $ mkdir -p $HOME/go/src/github.com/direnv

Clone project into that directory:

    $ git clone git@github.com:direnv/direnv.git $HOME/go/src/github.com/direnv/direnv

Build by just typing make:

    $ cd $HOME/go/src/github.com/direnv/direnv
    $ make

To install to /usr/local:

    $ make install

Or to a different location like `~/.local`:

    $ make install DESTDIR=~/.local


## Setup

For direnv to work properly it needs to be hooked into the shell. Each shell
has its own extension mechanism:

### BASH

Add the following line at the end of the `~/.bashrc` file:

```sh
eval "$(direnv hook bash)"
```

Make sure it appears even after rvm, git-prompt and other shell extensions
that manipulate the prompt.

### ZSH

Add the following line at the end of the `~/.zshrc` file:

```sh
eval "$(direnv hook zsh)"
```

### FISH

Add the following line at the end of the `~/.config/fish/config.fish` file:

```fish
eval (direnv hook fish)
```

### TCSH

Add the following line at the end of the `~/.cshrc` file:

```sh
eval `direnv hook tcsh`
```

### Elvish

Run:

```
$> direnv hook elvish > ~/.elvish/lib/direnv.elv
```

and add the following line to your `~/.elvish/rc.elv` file:

```
use direnv
```

## Usage

In some target folder, create a ".envrc" file and add some export(1)
directives in it.

Note that the contents of the `.envrc` file must be valid bash syntax,
regardless of the shell you are using.
This is because direnv always executes the `.envrc` with bash (a sort of
lowest common denominator of UNIX shells) so that direnv can work across shells.
If you try to use some syntax that doesn't work in bash (like zsh's
nested expansions), you will [run into
trouble](https://github.com/direnv/direnv/issues/199).

On the next prompt you will notice that direnv complains about the ".envrc"
being blocked. This is the security mechanism to avoid loading new files
automatically. Otherwise any git repo that you pull, or tar archive that you
unpack, would be able to wipe your hard drive once you `cd` into it.

So here we are pretty sure that it won't do anything bad. Type `direnv allow .`
and watch direnv loading your new environment. Note that `direnv edit .` is a
handy shortcut that opens the file in your $EDITOR and automatically allows it
if the file's modification time has changed.

Now that the environment is loaded, you will notice that once you `cd` out
of the directory it automatically gets unloaded. If you `cd` back into it, it's
loaded again. That's the basis of the mechanism that allows you to build cool
things.

### The stdlib

Exporting variables by hand is a bit repetitive so direnv provides a set of
utility functions that are made available in the context of the ".envrc" file.

As an example, the `PATH_add` function is used to expand and prepend a path to
the $PATH environment variable. Instead of `export $PATH=$PWD/bin:$PATH` you
can write `PATH_add bin`. It's shorter and avoids a common mistake where
`$PATH=bin`.

To find the documentation for all available functions check the
direnv-stdlib(1) man page.

It's also possible to create your own extensions by creating a bash file at
`~/.config/direnv/direnvrc` or `~/.direnvrc`. This file is loaded before your
".envrc" and thus allows you to make your own extensions to direnv.

#### Loading layered .envrc

Let's say you have the following structure:

- "/a/.envrc"
- "/a/b/.envrc"

If you add the following line in "/a/b/.envrc", you can load both of the
".envrc" files when you are in `/a/b`:

```sh
source_env ..
```
In the general case `source_up` will load any .envrc higher up in the folder structure. This allows you to truely enable arbitrary heirarchical stuctures of `.envrc` usage.

```sh
source_up
```

## Common things people don't know

Based on GitHub issues interractions, here are the top things that have been confusing for users:

1. direnv has a standard library of functions, a collection of utilities that I found useful to have and accumulated over the years. If you know how to read bash, you can find it here: https://github.com/direnv/direnv/blob/master/stdlib.sh

2. It's possible to override the stdlib with your own set of function by adding a bash file to either `~/.config/direnv/direnvrc` or `~/.direnvrc`. These will become available to all your `.envrc` files.

3. direnv is actually creating a new bash process to load the stdlib, direnvrc and `.envrc`, and only exports the environment diff back to the original shell. This allows direnv to record the environment changes accurately and also work with all sorts of shells. It also means that aliases and functions are not exportable right now.

## Similar projects

* [Environment Modules](http://modules.sourceforge.net/) - one of the oldest (in a good way) environment-loading systems
* [autoenv](https://github.com/kennethreitz/autoenv) - lightweight; doesn't support unloads
* [zsh-autoenv](https://github.com/Tarrasch/zsh-autoenv) - a feature-rich mixture of autoenv and [smartcd](https://github.com/cxreg/smartcd): enter/leave events, nesting, stashing (Zsh-only).

## Contribute

Bug reports, contributions and forks are welcome.

All bugs or other forms of discussion happen on
<http://github.com/direnv/direnv/issues>.

There is a wiki available where you can share your usage patterns or
other tips and tricks <https://github.com/direnv/direnv/wiki>

For longer form discussions you can also write to <mailto:direnv-discuss@googlegroups.com>

Or drop by on [IRC (#direnv on freenode)](irc://irc.freenode.net/#direnv) to
have a chat. If you ask a question make sure to stay around as not everyone is
active all day.

[![Build Status](https://api.travis-ci.org/direnv/direnv.png?branch=master)](http://travis-ci.org/direnv/direnv)

## COPYRIGHT

Copyright (C) 2014 shared by all
[contributors](https://github.com/direnv/direnv/graphs/contributors) under
the MIT licence.
