
2.17.0 / 2018-06-17
==================

  * CHANGE: hook expands the direnv path. Ensures that direnv can be executed even if the PATH is changed #369.
  * CHANGE: stdlib: direnv_load: disallow watching in child
    Allows the `use nix --pure` scenario in #368
  * README: add OpenSuSE to the list of distros
  * Revert "use_nix: unset IN_NIX_SHELL"

2.16.0 / 2018-05-09
==================

  * NEW: add support for elvish (#356)
  * NEW: config: allow to disable stdin on eval (#351)
  * DOC: Add the usage of source_up to the README (#347)
  * FIX: default.nix: fix compilation

2.15.2 / 2018-02-25
==================

  * FIX: lintian warnings (#340)
  * FIX: release process (#342)

2.15.1 / 2018-02-24
==================

  * FIX: support for go 1.10 (#339)

2.15.0 / 2018-02-23
==================

  * NEW: TOML configuration file! (#332, #337)
  * NEW: support for allow folder whitelist (#332)
  * NEW: add anaconda support (#312)
  * CHANGE: use_nix: unset IN_NIX_SHELL

2.14.0 / 2017-12-13
==================

  * NEW: Add support for Pipenv layout (#314)
  * CHANGE: direnv version: make public
  * FIX: direnv edit: run the command through bash
  * FIX: website: update ditto to v0.15

2.13.3 / 2017-11-30
==================

  * FIX: fixes dotenv loading issue on macOS `''=''`

2.13.2 / 2017-11-28
==================

  * FIX: direnv edit: fix path escaping
  * FIX: stdlib: fix find_up
  * FIX: stdlib: use absolute path in source_up
  * FIX: remove ruby as a build dependency
  * FIX: go-dotenv: update to latest master to fix a parsing error

2.13.1 / 2017-09-27
==================

  * FIX: stdlib: make direnv_layout_dir lazy (#298)

2.13.0 / 2017-09-24
==================

  * NEW: stdlib: configurable direnv_layout_dir
  * CHANGE: stdlib: source the direnvrc directly
  * FIX: permit empty NODE_VERSION_PREFIX variable
  * FIX: pwd: Don't use -P to remove symlinks (#295)
  * FIX: also reload when mtime goes back in time
  * FIX: Prevent `$HOME` path from being striked (#287)
  * BUILD: use the new `dep` tool to manage dependencies
  * BUILD: dotenv: move to vendor folder

2.12.2 / 2017-07-05
==================

  * stdlib layout_python: fixes on no arg

2.12.1 / 2017-07-01
==================

  * FIX: stdlib path_add(), see #278
  * FIX: install from source instructions

2.12.0 / 2017-06-30
==================

  * NEW: support multiple items in path_add and PATH_add (#276)
  * NEW: add a configurable DIRENV_WARN_TIMEOUT option (#273)
  * CHANGE: rewrite the dotenv parsing, now supports commented lines
  * CHANGE: pass additional args to virtualenv (#261)
  * FIX: stdlib watch_file(): escaping fix
  * FIX: only output color if $TERM is not dumb (#264)
  * FIX: the watch_file documentation

2.11.3 / 2017-03-02
==================

  * FIX: node version sorting (#255)

2.11.2 / 2017-03-01
==================

  * FIX: Typo in MANPATH_add always generates "PATH missing" error. (#256)

2.11.1 / 2017-02-20
==================

  * FIX: only deploy the go 1.8 version

2.11.0 / 2017-02-20
==================

  * NEW: stdlib.sh: introduce MANPATH_add <path> (#248)
  * NEW: provide packages using the equinox service
  * CHANGE: test direnv with go 1.8 (#254)
  * FIX: Add warning about source_env/up
  * FIX: go-md2man install instruction

2.10.0 / 2016-12-10
==================

  * NEW: `use guix` (#242)
  * CHANGE: use go-md2man to generate the man pages
  * FIX: tcsh escaping (#241)
  * FIX: doc typos and rewords (#226)

2.9.0 / 2016-07-03
==================

  * NEW: use_nix() is now watching default.nix and shell.nix
  * NEW: Allow to fix the bash path at built time
  * FIX: Panic on `direnv current` with no argument
  * FIX: Permit empty NODE_VERSION_PREFIX variable
  * FIX: layout_python: fail properly when python is not found

2.8.1 / 2016-04-04
==================

  * FIX: travis dist release

2.8.0 / 2016-03-27
==================

  * NEW: `direnv export json` to facilitate IDE integration
  * NEW: watch functionality thanks to @avnik
    Now direnv also reload on associate .env and .envrc changes.
  * NEW: stdlib `watch_file` function thanks to @avnik
    Allows to monitor more files for change.
  * NEW: stdlib `use node` function thanks to @wilmoore
  * NEW: `direnv prune` to remove old allowed files thanks to @punitagrawal
    Only works with newly-generated files since we're not storing the path
    inside of them.

2.7.0 / 2015-08-08
==================

  * NEW: use_nix() helper to stdlib. Thanks @gfxmonk
  * FIX: Added SHELLOPTS to ignored vars. Thanks @fernandomora
  * FIX: Removed shellcheck offenses in the stdlib, better escaping
  * FIX: typos. Thanks @camelpunch, @oppegard

2.6.1 / 2015-06-23
==================

  * FIX: source_env handles missing .envrc gracefully. Thanks @gerhard
  * FIX: Empty variable as unloading in Vim. Thanks @p0deje
  * FIX: Corrected spelling mistake in deny command. Thanks @neanias

2.6.0 / 2015-02-15
==================

  * NEW: tcsh is now supported ! Thanks @bbense
  * CHANGE: `direnv dump` now ignores `BASH_FUNC_` exports. Thanks @gfxmonk
  * CHANGE: Interactive input during load is now possible. Thanks @toao
  * FIX: allow workaround for tmux users: `alias tmux='direnv exec / tmux'`
  * FIX: hardened fish shell escaping thanks to @gfxmonk

Thanks @bbense @vially and @dadooda for corrections in the docs

2.5.0 / 2014-11-04
==================

  * NEW: Use a different virtualenv per python versions for easier version
    switching. Eg: ./.direnv/python-${python_version}
  * NEW: Makes `layout python3` a shortcut for `layout python python3`. Thanks
    @ghickman !
  * NEW: Allows to specify which executable of python to use in `layout_python`
  * CHANGE: `layout python` now unsets $PYTHONHOME to better mimic virtualenv
  * CHANGE: Don't make virtualenvs relocatable. Fixes #137
  * OTHER: Use Travis to push release builds to github

2.4.0 / 2014-06-15
==================

 * NEW: Try to detect an editor in the PATH if EDITOR is not set.
 * NEW: Preliminary support for vim
 * NEW: New site: put the doc inside the project so it stays in sync
 * NEW: Support for Cygwin - Thanks @CMCDragonkai !
 * NEW: Allow to disable logging by setting an empty `DIRENV_LOG_FORMAT`
 * NEW: stdlib `layout perl`. Thanks @halkeye !
 * CHANGE: layout ruby: share the gem home starting from rubygems v2.2.0
 * CHANGE: Allow arbitrary number of args in `log_status`
 * CHANGE: Bump command timeout to 5 seconds
 * FIX: Adds selected bash executable in `direnv status`
 * FIX: man changes, replaced abandonned ronn by md2man
 * FIX: `make install` was creating a ./bin directory
 * FIX: issue #114 - work for blank envs. Thanks @pwaller !
 * FIX: man pages warning. Thanks @punitagrawal !
 * FIX: Multi-arg EDITOR was broken #108
 * FIX: typos in doc. Thanks @HeroicEric and @lmarlow !
 * FIX: If two paths don't have a common ancestors, don't make them relative.
 * FIX: missing doc on layered .envrc. Thanks @take !

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

