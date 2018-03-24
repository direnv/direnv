DIRENV.TOML 1 "FEBRUARY 2018" direnv "User Manuals"
==================================================

NAME
----

direnv.toml - the direnv configuration file

DESCRIPTION
-----------

A configuration file in [TOML](https://github.com/toml-lang/toml) format to specify a variety of configuration options for direnv. Resides at CONFIGURATION_DIR/config.toml. For many users, this will be located at $HOME/.config/direnv/config.toml.

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

## `whitelist`

Specifying whitelist directives marks specific directory hierarchies or specific directories as "trusted" -- direnv will evaluate any matching .envrc files regardless of whether they have been specifically allowed. **This feature should be used with great care**, as anyone with the ability to write files to that directory (including collaborators on VCS repositories) will be able to execute arbitrary code on your computer.

There are two types of whitelist directives supported:

### `prefix`
Accepts an array of strings. If any of the strings in this list are a prefix of an .envrc file's absolute path, that file will be implicitly allowed, regardless of contents or past usage of `direnv allow` or `direnv deny`.

Example:

```toml
[whitelist]
prefix = [ "/home/user/code/project-a" ]
```

In this example, the following .envrc files will be implicitly allowed:

* `/home/user/code/project-a/.envrc`
* `/home/user/code/project-a/subdir/.envrc`
* and so on

In this example, the following .envrc files will not be implicitly allowed (although they can be explicitly allowed by running `direnv allow`):

* `/home/user/project-b/.envrc`
* `/opt/random/.envrc`

### `exact`

Accepts an array of strings. Each string can be a directory name or the full path to an .envrc file. If a directory name is passed, it will be treated as if it had been passed as itself with `/.envrc` appended. After resolving the filename, each string will be checked for being an exact match with an .envrc file's absolute path. If they match exactly, that .envrc file will be implicitly allowed, regardless of contents or past usage of `direnv allow` or `direnv deny`.

Example:

```toml
[whitelist]
exact = [ "/home/user/project-b/.envrc", "/home/user/project-b/subdir-a" ]
```

In this example, the following .envrc files will be implicitly allowed, and no others:

* `/home/user/code/project-b/.envrc`
* `/home/user/code/project-b/subdir-a`

In this example, the following .envrc files will not be implicitly allowed (although they can be explicitly allowed by running `direnv allow`):

* `/home/user/code/project-b/subproject-c/.envrc`
* `/home/user/code/.envrc`

## `bash_path`

This allows one to hard-code the position of bash. It maybe be useful to set this to avoid having direnv to fail when PATH is being mutated.

## `disable_stdin`

If set to true, stdin is disabled (redirected to /dev/null) during the `.envrc` evaluation.

COPYRIGHT
---------

Copyright (C) 2018 zimbatm <http://zimbatm.com> and contributors under the MIT licence.

SEE ALSO
--------

direnv(1)
