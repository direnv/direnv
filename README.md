direnv -- Unclutter your .profile
=================================

`direnv` is an environment variable manager for your shell. It knows how to
hook into bash, zsh and fish shell to load or unload environment variables
depending on your current directory. This allows to have project-specific
environment variables and not clutter the "~/.profile" file.

Before each prompt it checks for the existence of an ".envrc" file in the
current and parent directories. If the file exists, it is loaded into a bash
sub-shell and all exported variables are then captured by direnv and then made
available to your shell.

Because direnv is compiled into a single static executable it is fast enough
to be unnoticeable on each prompt. It is also language agnostic and can be
used to build solutions similar to rbenv, pyenv, phpenv, ...


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

### From source

Dependencies: make, golang

```bash
git clone http://github.com/zimbatm/direnv
cd direnv
make install
# or symlink bin/direnv into your $PATH
```

### Packaged

There's package definitions on Homebrew, Arch's AUR and NixOS's nixpkgs.

Links to binary builds are also available on each release.

## Setup

For direnv to work properly it needs to be hooked into the shell. Each shell
has it's own extension mechanism:

### BASH

Add the following line at the end of your "~/.bashrc" file:

`eval "$(direnv hook bash)"`

Make sure it appears even after rvm, git-prompt and other shell extensions
that manipulate your prompt.

### ZSH

Add the previous line at the end of you "~/.zshrc" file:

`eval "$(direnv hook zsh)"`

### FISH

Add the previous line at the end of your "~/.config/fish/config.fish" file:

`eval (direnv hook fish)`

## Usage

In some target folder, create an ".envrc" file and add some export(1)
directives in it.

On the next prompt you will notice that direnv complains about the ".envrc"
being blocked. This is the security mechanism to avoid loading new files
automatically. Otherwise and git repo that you pull, or tar archive that you
unpack, would be able to wipe your hard drive once you `cd` into it.

So here we are pretty sure that it won't do anything bad. Type `direnv allow .`
and watch direnv loading your new environment. Note that `direnv edit .` is a
handy shortcut that open the file in your $EDITOR and automatically allows it
if the file's modification time has changed.

Now that the environment is loaded you can notice that once your `cd` out
of the directory it automatically gets unloaded. If you `cd` back into it it's
loaded again. That's the base of the mechanism that allows you to build cool
things.

Exporting variables by hand is a bit repetitive so direnv provides a set of
utility functions that are made available in the context of the ".envrc" file.
Check the direnv-stdlib(1) man page for more details. You can also define
your own extensions inside a "~/.direnvrc" file.

### Loading layered .envrc

Lets say you have the following structure:

- `/a/.envrc`
- `/a/b/.envrc`

If you add the following line in `/a/b/.envrc`, you can load both of the `.envrc` when you are in `/a/b`:

```
source_env ..
```

## Contribute

Bug reports, contributions and forks are welcome.

All bugs or other forms of discussion happen on
<http://github.com/zimbatm/direnv/issues>

There is also a wiki available where you can share your usage patterns or
other tips and tricks <https://github.com/zimbatm/direnv/wiki>

Or drop by on the [#direnv channel on FreeNode](irc://#direnv@FreeNode) to
have a chat.

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

Copyright (C) 2014 zimbatm and contributors under the MIT licence.
