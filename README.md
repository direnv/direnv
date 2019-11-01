direnv -- unclutter your .profile
=================================

[![Built with Nix](https://builtwithnix.org/badge.svg)](https://builtwithnix.org)
[![Build Status](https://dev.azure.com/direnv/direnv/_apis/build/status/direnv.direnv?branchName=master)](https://dev.azure.com/direnv/direnv/_build/latest?definitionId=1&branchName=master)
[![Packaging status](https://repology.org/badge/tiny-repos/direnv.svg)](https://repology.org/project/direnv/versions)
[![latest packaged version(s)](https://repology.org/badge/latest-versions/direnv.svg)](https://repology.org/project/direnv/versions)

`direnv` is an extension for your shell. It augments existing shells with a
new feature that can load and unload environment variables depending on the
current directory.

## Use cases

* Load 12factor apps environment variables
* Create per-project isolated development environments
* Load secrets for deployment

## How it works

Before each prompt, direnv checks for the existence of a `.envrc` file in the
current and parent directories. If the file exists (and is authorized), it is
loaded into a **bash** sub-shell and all exported variables are then captured
by direnv and then made available to the current shell.

It supports hooks for all the common shells like bash, zsh, tcsh and fish.
This allows project-specific environment variables without cluttering the
`~/.profile` file.

Because direnv is compiled into a single static executable, it is fast enough
to be unnoticeable on each prompt. It is also language-agnostic and can be
used to build solutions similar to rbenv, pyenv and phpenv.

## Getting Started

### Prerequisites

* Unix-like operating system (macOS, Linux, ...)
* A supported shell (bash, zsh, tcsh, fish, elvish)

### Basic Installation

1. direnv is packaged in most distributions already. See [the installation documentation](docs/installation.md) for details.
2. [hook direnv into your shell](docs/hook.md).

Now restart your shell.

### Quick demo

To follow along in your shell once direnv is installed.

```
# Create a new folder for demo purposes.
$ mkdir ~/my-project
$ cd ~/my-project

# Show that the FOO environment variable is not loaded.
$ echo ${FOO-nope}
nope

# Create a new .envrc. This file is bash code that is going to be loaded by
# direnv.
$ echo export FOO=foo > .envrc
.envrc is not allowed

# The security mechanism didn't allow to load the .envrc. Since we trust it,
# let's allow it's execution.
$ direnv allow .
direnv: reloading
direnv: loading .envrc
direnv export: +FOO

# Show that the FOO environment variable is loaded.
$ echo ${FOO-nope}
foo

# Exit the project
$ cd ..
direnv: unloading

# And now FOO is unset again
$ echo ${FOO-nope}
nope
```

### The stdlib

Exporting variables by hand is a bit repetitive so direnv provides a set of
utility functions that are made available in the context of the `.envrc` file.

As an example, the `PATH_add` function is used to expand and prepend a path to
the $PATH environment variable. Instead of `export PATH=$PWD/bin:$PATH` you
can write `PATH_add bin`. It's shorter and avoids a common mistake where
`$PATH=bin`.

To find the documentation for all available functions check the
[direnv-stdlib(1) man page](man/direnv-stdlib.1.md).

It's also possible to create your own extensions by creating a bash file at
`~/.config/direnv/direnvrc` or `~/.direnvrc`. This file is loaded before your
`.envrc` and thus allows you to make your own extensions to direnv.

## Docs

* [Install direnv](docs/installation.md)
* [Hook into your shell](docs/hook.md)
* [Develop for direnv](docs/development.md)
* [Manage your rubies with direnv and ruby-install](docs/ruby.md)
* [Community Wiki](https://github.com/direnv/direnv/wiki)

Make sure to take a look at the wiki! It contains all sorts of useful
information such as common recipes, editor integration, tips-and-tricks.

### Man pages

* [direnv(1) man page](man/direnv.1.md)
* [direnv-stdlib(1) man page](man/direnv-stdlib.1.md)
* [direnv.toml(1) man page](man/direnv.toml.1.md)

### FAQ

Based on GitHub issues interactions, here are the top things that have been
confusing for users:

1. direnv has a standard library of functions, a collection of utilities that
   I found useful to have and accumulated over the years. You can find it
   here: https://github.com/direnv/direnv/blob/master/stdlib.sh

2. It's possible to override the stdlib with your own set of function by
   adding a bash file to `~/.config/direnv/direnvrc`. This file is loaded and
   it's content made available to any `.envrc` file.

3. direnv is not loading the `.envrc` into the current shell. It's creating a
   new bash sub-process to load the stdlib, direnvrc and `.envrc`, and only
   exports the environment diff back to the original shell. This allows direnv
   to record the environment changes accurately and also work with all sorts
   of shells. It also means that aliases and functions are not exportable
   right now.

## Contributing

Bug reports, contributions and forks are welcome. All bugs or other forms of
discussion happen on http://github.com/direnv/direnv/issues .

Or drop by on [IRC (#direnv on freenode)](irc://irc.freenode.net/#direnv) to
have a chat. If you ask a question make sure to stay around as not everyone is
active all day.

## Similar projects

Here is a list of other projects found in the same design space. Feel free to
submit new ones.

* [Environment Modules](http://modules.sourceforge.net/) - one of the oldest (in a good way) environment-loading systems
* [autoenv](https://github.com/kennethreitz/autoenv) - lightweight; doesn't support unloads
* [zsh-autoenv](https://github.com/Tarrasch/zsh-autoenv) - a feature-rich mixture of autoenv and [smartcd](https://github.com/cxreg/smartcd): enter/leave events, nesting, stashing (Zsh-only).
* [asdf](https://github.com/asdf-vm/asdf) - a pure bash solution that has a plugin system. The [asdf-direnv](https://github.com/asdf-community/asdf-direnv) plugin allows using asdf managed tools with direnv.
* [ondir](https://github.com/alecthomas/ondir) - OnDir is a small program to automate tasks specific to certain directories

## COPYRIGHT

[MIT licence](LICENSE) - Copyright (C) 2019 @zimbatm and [contributors](https://github.com/direnv/direnv/graphs/contributors)
