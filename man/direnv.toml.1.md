DIRENV.TOML 1 "2019" direnv "User Manuals"
==========================================

NAME
----

direnv.toml - the direnv configuration file

DESCRIPTION
-----------

A configuration file in [TOML](https://github.com/toml-lang/toml) format to specify a variety of configuration options for direnv. The file is read from `$XDG_CONFIG_HOME/direnv/direnv.toml` (typically ~/.config/direnv/direnv.toml).

> For versions v2.21.0 and below use config.toml instead of direnv.toml

FORMAT
------

See the [TOML GitHub Repository](https://github.com/toml-lang/toml) for details about the syntax of the configuration file.

CONFIG
------

The configuration is specified in sections which each have their own top-level [tables](https://github.com/toml-lang/toml#table), with key/value pairs specified in each section.

Example:

```toml
[section]
key = "value"
```

The following sections are supported:

## [global]

### `bash_path`

This allows one to hard-code the position of bash. It maybe be useful to set this to avoid having direnv to fail when PATH is being mutated.

### `disable_stdin`

If set to `true`, stdin is disabled (redirected to /dev/null) during the `.envrc` evaluation.

### `load_dotenv`

> direnv >= 2.31.0 is required

If set to `true`, also look for and load `.env` files on top of the `.envrc` files. If both `.envrc` and `.env` files exist, the `.envrc` will always be chosen first.

### `strict_env`

If set to `true`, the `.envrc` will be loaded with `set -euo pipefail`. This
option will be the default in the future.

### `warn_timeout`

Specify how long to wait before warning the user that the command is taking
too long to execute. Defaults to "5s".

A duration string is a possibly signed sequence of decimal numbers, each with
optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m".
Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".

This feature is disabled if the duration is lower or equal to zero.
Will be overwritten if the environment variable `DIRENV_WARN_TIMEOUT` is set to any of the above values.

### `hide_env_diff`

Set to `true` to hide the diff of the environment variables when loading the
`.envrc`. Defaults to `false`.

## [whitelist]

Specifying whitelist directives marks specific directory hierarchies or specific directories as "trusted" -- direnv will evaluate any matching .envrc files regardless of whether they have been specifically allowed. **This feature should be used with great care**, as anyone with the ability to write files to that directory (including collaborators on VCS repositories) will be able to execute arbitrary code on your computer.

There are two types of whitelist directives supported:

### `prefix`

Accepts an array of strings. If any of the strings in this list are a prefix of an .envrc file's absolute path, that file will be implicitly allowed, regardless of contents or past usage of `direnv allow` or `direnv deny`.

Example:

```toml
[whitelist]
prefix = [ "/home/user/code/project-a", "~/code/project-b" ]
```

In this example, the following .envrc files will be implicitly allowed:

* `/home/user/code/project-a/.envrc`
* `/home/user/code/project-a/subdir/.envrc`
* `~/code/project-b/.envrc`
* `~/code/project-b/subdir/.envrc`
* and so on

In this example, the following .envrc files will not be implicitly allowed (although they can be explicitly allowed by running `direnv allow`):

* `/home/user/project-c/.envrc`
* `/opt/random/.envrc`

### `exact`

Accepts an array of strings. Each string can be a directory name or the full path to an .envrc file. If a directory name is passed, it will be treated as if it had been passed as itself with `/.envrc` appended. After resolving the filename, each string will be checked for being an exact match with an .envrc file's absolute path. If they match exactly, that .envrc file will be implicitly allowed, regardless of contents or past usage of `direnv allow` or `direnv deny`.

Example:

```toml
[whitelist]
exact = [ "/home/user/project-a/.envrc", "~/project-b/subdir-a" ]
```

In this example, the following .envrc files will be implicitly allowed, and no others:

* `/home/user/code/project-a/.envrc`
* `~/project-b/subdir-a`

In this example, the following .envrc files will not be implicitly allowed (although they can be explicitly allowed by running `direnv allow`):

* `/home/user/code/project-b/subproject-c/.envrc`
* `~/code/.envrc`

COPYRIGHT
---------

MIT licence - Copyright (C) 2019 @zimbatm and contributors

SEE ALSO
--------

direnv(1), direnv-stdlib(1)
