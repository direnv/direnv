# Setup

For direnv to work properly it needs to be hooked into the shell. Each shell
has its own extension mechanism.

Once the hook is configured, restart your shell for direnv to be activated.

## BASH

Add the following line at the end of the `~/.bashrc` file:

```sh
eval "$(direnv hook bash)"
```

Make sure it appears even after rvm, git-prompt and other shell extensions
that manipulate the prompt.

## ZSH

Add the following line at the end of the `~/.zshrc` file:

```sh
eval "$(direnv hook zsh)"
```

## FISH

Add the following line at the end of the `~/.config/fish/config.fish` file:

```fish
direnv hook fish | source
```

Fish supports 3 modes you can set with with the global environment variable `direnv_fish_mode`:

```fish
set -g direnv_fish_mode eval_on_arrow    # trigger direnv at prompt, and on every arrow-based directory change (default)
set -g direnv_fish_mode eval_after_arrow # trigger direnv at prompt, and only after arrow-based directory changes before executing command
set -g direnv_fish_mode disable_arrow    # trigger direnv at prompt only, this is similar functionality to the original behavior
```

## TCSH

Add the following line at the end of the `~/.cshrc` file:

```sh
eval `direnv hook tcsh`
```

## Elvish (0.12+)

Run:

```
$> direnv hook elvish > ~/.elvish/lib/direnv.elv
```

and add the following line to your `~/.elvish/rc.elv` file:

```
use direnv
```
