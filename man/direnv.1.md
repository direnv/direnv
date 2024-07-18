DIRENV 1 "2019" direnv "User Manuals"
===========================================

NAME
----

direnv - unclutter your .profile

SYNOPSIS
--------

`direnv` *command* ...

DESCRIPTION
-----------

`direnv` is an environment variable manager for your shell. It knows how to
hook into bash, zsh and fish shell to load or unload environment variables
depending on your current directory. This allows you to have project-specific
environment variables and not clutter the "~/.profile" file.

Before each prompt it checks for the existence of an `.envrc` file in the
current and parent directories. If the file exists, it is loaded into a bash
sub-shell and all exported variables are then captured by direnv and then made
available to your current shell, while unset variables are removed.

Because direnv is compiled into a single static executable it is fast enough
to be unnoticeable on each prompt. It is also language agnostic and can be
used to build solutions similar to rbenv, pyenv, phpenv, ...

EXAMPLE
-------

```
$ cd ~/my_project
$ echo ${FOO-nope}
nope
$ echo export FOO=foo > .envrc
\.envrc is not allowed
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

SETUP
-----

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

Add the following line at the end of the `$XDG_CONFIG_HOME/fish/config.fish` file:

```fish
direnv hook fish | source
```

Fish supports 3 modes you can set with with the global environment variable `direnv_fish_mode`:

```fish
set -g direnv_fish_mode eval_on_arrow    # trigger direnv at prompt, and on every arrow-based directory change (default)
set -g direnv_fish_mode eval_after_arrow # trigger direnv at prompt, and only after arrow-based directory changes before executing command
set -g direnv_fish_mode disable_arrow    # trigger direnv at prompt only, this is similar functionality to the original behavior
```


### TCSH

Add the following line at the end of the `~/.cshrc` file:

```sh
eval `direnv hook tcsh`
```

### Elvish

Run:

```
~> mkdir -p ~/.config/elvish/lib
~> direnv hook elvish > ~/.config/elvish/lib/direnv.elv
```

and add the following line to your `~/.config/elvish/rc.elv` file:

```
use direnv
```

### PowerShell

Add the following line to your `$PROFILE`:

```powershell
Invoke-Expression "$(direnv hook pwsh)"
```

### Murex

Add the following line to your `~/.murex_profile`:

```
direnv hook murex -> source
```

USAGE
-----

In some target folder, create an `.envrc` file and add some export(1)
and unset(1) directives in it.

On the next prompt you will notice that direnv complains about the `.envrc`
being blocked. This is the security mechanism to avoid loading new files
automatically. Otherwise any git repo that you pull, or tar archive that you
unpack, would be able to wipe your hard drive once you `cd` into it.

So here we are pretty sure that it won't do anything bad. Type `direnv allow .`
and watch direnv loading your new environment. Note that `direnv edit .` is a
handy shortcut that opens the file in your $EDITOR and automatically reloads it
if the file's modification time has changed.

Now that the environment is loaded you can notice that once you `cd` out
of the directory it automatically gets unloaded. If you `cd` back into it it's
loaded again. That's the base of the mechanism that allows you to build cool
things.

Exporting variables by hand is a bit repetitive so direnv provides a set of
utility functions that are made available in the context of the `.envrc` file.
Check the direnv-stdlib(1) man page for more details. You can also define your
own extensions inside `$XDG_CONFIG_HOME/direnv/direnvrc` or
`$XDG_CONFIG_HOME/direnv/lib/*.sh` files.

Hopefully this is enough to get you started.

ENVIRONMENT
-----------

`XDG_CONFIG_HOME`
: Defaults to `$HOME/.config`.

`XDG_DATA_HOME`
: Defaults to `$HOME/.local/share`.

FILES
-----

`$XDG_CONFIG_HOME/direnv/direnv.toml`
: Direnv configuration. See direnv.toml(1).

`$XDG_CONFIG_HOME/direnv/direnvrc`
: Bash code loaded before every `.envrc`. Good for personal extensions.

`$XDG_CONFIG_HOME/direnv/lib/*.sh`
: Bash code loaded before every `.envrc`. Good for third-party extensions.

`$XDG_DATA_HOME/direnv/allow`
: Records which `.envrc` files have been `direnv allow`ed.

CONTRIBUTE
----------

Bug reports, contributions and forks are welcome.

All bugs or other forms of discussion happen on
<http://github.com/direnv/direnv/issues>

There is also a wiki available where you can share your usage patterns or
other tips and tricks <https://github.com/direnv/direnv/wiki>

Or drop by on the [#direnv channel on FreeNode](irc://#direnv@FreeNode) to
have a chat.

COPYRIGHT
---------

MIT licence - Copyright (C) 2019 @zimbatm and contributors

SEE ALSO
--------

direnv-stdlib(1), direnv.toml(1), direnv-fetchurl(1)
