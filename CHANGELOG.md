
2.3.0 / 2014-02-06
==================

 * NEW: DIRENV_LOG_FORMAT environment variable can be used tocontrol log formatting
 * NEW: `direnv exec [DIR] <COMMAND>` to execute programs with an .envrc context
 * CHANGE: layout_python now tries to make your virtualenv relocatable
 * CHANGE: the export diff is not from the old env, not the current env
 * CHANGE: layout_go now also adds $PWD/bin in the PATH
 * FIX: Hides the DIRENV_ variables in the output diff. Fixes #94
 * FIX: Makes sure the path used in the allow hash is absolute. See #95
 * FIX: Set the executable bit on direnv on install
 * FIX: Some bash installs had a parse error in the hook.

2.2.1 / 2014-01-12
==================

The last release was heavily broken. Ooops !

 * FIX: Refactored the whole export and diff mechanism. Fixes #92 regression.
 * CHANGE: DIRENV_BACKUP has been renamed to DIRENV_DIFF

2.2.0 / 2014-01-11
==================

Restart your shells on upgrade, the format of DIRENV_BACKUP has changed and is
incompatible with previous versions.

 * NEW: `direnv_load <command-that-outputs-a-direnv-dump>` stdlib function
 * CHANGE: Only backup the diff of environments. Fixes #82
 * CHANGE: Renames `$DIRENV_PATH` to `$direnv` in the stdlib.
 * CHANGE: Allow/Deny mechanism now includes the path to make it more secure.
 * CHANGE: `direnv --help` is an alias to `direnv help`
 * CHANGE: more consistent log outputs and error messages
 * CHANGE: `direnv edit` only auto-allows the .envrc if it's mtime has changed.
 * CHANGE: Fixes old bash (OSX) segfault in some cases. See #81
 * CHANGE: The stdlib `dotenv` now supports more .env syntax
 * FIX: Restore the environment properly after loading errors.

2.1.0 / 2013-11-10
==================

 * Added support for the fish shell. See README.md for install instructions.
 * Stop recommending using $0 to detect the shell. Fixes #64.
 * Makes the zsh hook resistant to double-hooking.
 * Makes the bash hook resistant to double-hooking.
 * More precise direnv allow error message. Fixes #72

2.0.1 / 2013-07-27
==================

 * Fixes shell detection corner case

2.0.0 / 2013-06-16
==================

When upgrading from direnv 1.x make sure to restart your shell. The rest is
relatively backward-compatible.

 * changed the execution model. Everything is in a single static executable
 * most of the logic has been rewritten in Go
 * robust shell escaping (supports UTF-8 in env vars)
 * robust eval/export loop, avoids retrys on every prompt if there is an error
 * stdlib: added the `dotenv [PATH]` command to load .env files
 * command: added `direnv reload` to force-reload your environment

