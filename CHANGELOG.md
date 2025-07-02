2.37.0 / 2025-07-02
==================

  * docs: add github-actions page
  * docs: document sub-commands
  * docs: fix link to guix manual (#1421)
  * docs: re-generate manpages
  * feat(direnv export gha): strengthen export format
  * feat: add windows arm64 target (#1444)
  * fix(powershell): "export pwsh" to resolve PowerShell special character issues (#1448)
  * fix(python): do not include patch level in virtual environment names (#1423)
  * fix(use_nix): always restore special variables (#1424)
  * fix: accept true as valid DIRENV_DEBUG value (#1365)
  * fix: add trailing newline to error messages (#1426)
  * fix: delete duplicate ansi escape code

v2.36.0 / 2025-04-11
==================

  * direnv now requires go 1.24 (#1384)
  * doc: Correct duplicate usage of 'with' in the direnv(1) (#1394)
  * doc: note direnv version for log_{format,filter} (#1369)
  * feat: Add `use_flox` to stdlib.sh (#1372)
  * feat: logging filter (#1336)
  * fix use_nix: unset TMPDIR variables (#1409)
  * fix: A more universal fix for the python 3.14 `find_spec` deprecation warning (#1382)
  * fix: Don't give an error when the current directory doesn't exist (#1395)
  * fix: add support to fully reproducible guix shells (#1392)
  * fix: assert minimum powershell version (#1385)
  * fix: escape newlines in generated vimscript (#1347)
  * fix: fix empty array error in install.sh (#1406)
  * fix: optionally authenticate against github api during install (#1337)
  * fix: use_guix: Enable the watching of Guix related files. (#1353)

2.35.0 / 2024-10-07
==================

  * doc: Add version requirement for load_dotenv option (#1326)
  * doc: change Guix link to its package search. (#1268)
  * doc: fix broken link (#1327)
  * doc: update elvish docs (#1305)
  * feat: add opam support (#1298)
  * fix: add NuShell into list of supported shells (#1260)
  * fix: close tmp file (#1272)
  * fix: direnv edit: use `editor` when EDITOR not found, closes #1246 (#1247)
  * fix: release script
  * fix: stdlib: enable flakes when use flake is used (#1299)
  * fix: stdlib: export GOBIN for layout_go (#1286)
  * fix: stdlib: update layout_python to resolve deprecation warning (#1176)
  * fix: using PWD in .env files (#1052)
  * test: Fix Murex python-layout test (#1293)

2.34.0 / 2024-03-01
==================

  * doc: README.md, man pages: it's typos (#1230)
  * doc: add shell setup instructions for oh-my-zsh (#1070)
  * doc: added fetchurl manpage link to README.md
  * doc: document XDG_DATA_HOME (#1185)
  * doc: update installation.md for Gentoo (#1206)
  * feat: add Murex support (#1242)
  * feat: added systemd shell for export (#1126)
  * feat: allow to disable warn timeouts (#1209)
  * feat: hide env diff (#1223, #1234)
  * feat: made 'direnv export' non private (#1229)
  * fix: `use_julia` should not set LD_LIBRARY_PATH (#900)
  * fix: add missing deps for release in go.mod
  * fix: avoid use of regex in bash hook output (#1043)
  * fix: direnv.toml.1.md: add examples for $HOME expansion
  * fix: stdlib: use_flake: don't keep old generations around (#1089)
  * fix: stdlib: use_node: strip leading v from version (#1071)
  * fix: support Bash 5.1 array PROMPT_COMMAND (#1208)
  * fix: update stdlib.sh to avoid deprecation warning (#1221)
  * fix: update zsh syntax in internal/cmd/shell_zsh.go (#1075)

2.33.0 / 2023-11-29
==================

  * doc: add a Nushell section to `hook.md` by @amtoine in https://github.com/direnv/direnv/pull/1175
  * doc: fix broken links in installation.md by @just1602 in https://github.com/direnv/direnv/pull/1110
  * doc: show how to run tests by @bukzor-sentryio in https://github.com/direnv/direnv/pull/1137
  * doc: update NixOS installation instructions by @Gerg-L in https://github.com/direnv/direnv/pull/1172
  * doc: update direnv.toml.1.md by @Ativerc in https://github.com/direnv/direnv/pull/1099
  * feat: `direnv status --json` by @shivaraj-bh in https://github.com/direnv/direnv/pull/1142
  * feat: add PowerShell Support by @bamsammich in https://github.com/direnv/direnv/pull/1171
  * feat: add mergify configuration by @Mic92 in https://github.com/direnv/direnv/pull/1147
  * feat: add support for armv7l platform in install.sh by @ardje in https://github.com/direnv/direnv/pull/1162
  * feat: add watch print command by @Mic92 in https://github.com/direnv/direnv/pull/1198
  * feat: alias `direnv disallow` to deny by @will in https://github.com/direnv/direnv/pull/1182
  * feat: stdlib: create CACHEDIR.TAG inside .direnv by @Mic92 in https://github.com/direnv/direnv/pull/1148
  * fix: `allowPath` for `LoadedRC` by @shivaraj-bh in https://github.com/direnv/direnv/pull/1157
  * fix: don't prompt to allow if user explicitly denied by @Gabriella439 in https://github.com/direnv/direnv/pull/1158
  * fix: man/direnv-stdlib: fix obsolete opam-env example by @mzacho in https://github.com/direnv/direnv/pull/1170
  * fix: print correct path in source_env log message by @wentasah in https://github.com/direnv/direnv/pull/1144
  * fix: quote tcsh $PATH, to avoid failure on whitespace by @bukzor-sentryio in https://github.com/direnv/direnv/pull/1139
  * fix: remove redundant nil check in `CommandsDispatch` by @Juneezee in https://github.com/direnv/direnv/pull/1166
  * fix: update nixpkgs and shellcheck by @Mic92 in https://github.com/direnv/direnv/pull/1146

2.32.3 / 2023-05-20
==================

  * fix: incorrect escape sequences during Loads under git-bash (Windows) (#1085)
  * fix: skip some tests for IBM Z mainframe's z/OS operating system (#1094)
  * fix: stdlib: use_guix: Switch to guix shell. (#1045)
  * fix: stat the already open rc file rather than another path based one on it (#1044)
  * fix: remove deprecated io/ioutil uses (#1042)
  * fix: spelling fixes (#1041)
  * fix: appease Go 1.19 gofmt (#1040)
  * fix: pass BASH_PATH to make, matches the nixpkgs derivation (#1006)
  * fix: stdlib/layout_python: exclude patchlevel from $python_version (#1033)
  * doc: add Windows installation with winget (#1096)
  * doc: link 12factor webpage for more clarity (#1095)
  * website: add Plausible analytics

2.32.2 / 2022-11-24
==================

  * doc: Add stdlib's layout_pyenv to docs (#969)
  * doc: Fix broken link (#991)
  * doc: Minor typo fix (#1013)
  * doc: `$XDG_CONFIG_HOME/direnv/direnv.toml` => add (typically ~/.config/direnv/direnv.toml) (#985)
  * doc: add quickenv to Related projects (#970)
  * feat: Update layout anaconda to accept a path to a yml file (#962)
  * feat: install.sh: can specify direnv version (#1012)
  * fix: elvish: replace deprecated `except` with `catch` (#987)
  * fix: installer.sh: make direnv executable for all
  * fix: path escaping (#975)
  * fix: stdlib: only use ANSI escape on TTY (#958)
  * fix: test: remove mentions of DIRENV_MTIME (#1009)
  * fix: test: use lowercase -d flag for base64 decoding of DIRENV_DIFF (#996)
  * update: build(deps): bump github.com/BurntSushi/toml from 1.1.0 to 1.2.0 (#974)

2.32.1 / 2022-06-21
==================

  * feat: Support custom VIRTUAL_ENV for layout_python (#876)
  * fix: vendor go-dotenv (#955)

2.32.0 / 2022-06-13
==================

  * feat: Add gha shell for GitHub Actions (#910)
  * feat: Enable ppc64le builds (#947)
  * feat: allow conda environment names to be detected from environment.yml (#909)
  * feat: source_up_if_exists: A strict_env compatible version of source_up (#921)
  * feat: Expand ~/ in whitelist paths (#931)
  * feat: Add "block" and "revoke" as aliases of the "deny" command (#935)
  * feat: Add "permit" and "grant" as aliases of the "allow" command (#935)
  * fix: update go-dotenv
  * fix: fetchurl: store files as hex (#930)
  * fix: fetchurl: only store 200 responses (#944)
  * fix: Ensure status log messages are printed with normal color (#884)
  * fix: Clarify handling of .env files (#941)
  * fix: Update shell_elvish.go (#896)
  * fix: stdlib.sh: remove dependency on tput (#932)
  * fix: Use setenv in vim to allow non alphanumeric vars (#901)
  * fix: install.sh: add information about bin_path (#920)
  * fix: Treat `mingw*` as windows (direnv/direnv#918) (#919)
  * fix: man: clarify paths (#929)
  * fix: installation.md: Fix Fedora package link (#915)
  * Merge pull request #874 from direnv/refactor
  * chore: rc: stop using --noprofile --norc
  * chore: rc: prepare stdin earlier
  * chore: rc: install interrupt handler earlier
  * chore: stdlib: factor out stdlib preparation
  * chore: fix CI
  * chore: source_env: show full path (#870)
  * chore: Sort shells in DetectShell
  * chore: Enable codeql action (#938)
  * chore: Set permissions for GitHub actions (#937)
  * go: bump golang.org/x/sys for linux/loong64 support (#946)
  * build(deps): bump actions/checkout from 2.4.0 to 3.0.0 (#922)
  * build(deps): bump actions/checkout from 3.0.0 to 3.0.1 (#933)
  * build(deps): bump actions/checkout from 3.0.1 to 3.0.2 (#936)
  * build(deps): bump actions/setup-go from 2.1.5 to 3.0.0 (#923)
  * build(deps): bump actions/setup-go from 3.0.0 to 3.1.0 (#943)
  * build(deps): bump actions/setup-go from 3.1.0 to 3.2.0 (#950)
  * build(deps): bump cachix/install-nix-action from 16 to 17 (#925)
  * build(deps): bump github.com/BurntSushi/toml from 0.4.1 to 1.1.0 (#924)

2.31.0 / 2022-03-26
==================

  * Don't load .env files by default (#911)
  * doc: `~/.config/direnv/direnvrc` is the default
  * doc: fix the broken link to arch linux (#892)
  * Re-add accidentally deleted comment line (#881)
  * fix version test

2.30.3 / 2022-01-05
==================

  * Allow skipping `.env` autoload (#878)
  * stdlib: add `env_vars_required` (#872) (#872)
  * Test whether version.txt contains semantic version (#871)

2.30.2 / 2021-12-28
==================

  * FIX: version: trim surrounding spaces (#869)
  * build(deps): bump actions/setup-go from 2.1.4 to 2.1.5 (#866)
  * move most code under internal/cmd (#865)

2.30.1 / 2021-12-24
==================

  * FIX: ignore .envrc and .env if they are not files (#864)

2.30.0 / 2021-12-23
==================

  * Add automatic `.env` load (#845)
  * Resolve symlinks during `direnv deny` (#851)
  * update installer for Apple Silicon (#849)
  * stdlib: use_flake handle no layout dir (#861)
  * embed stdlib.sh (#782)
  * embed version.txt
  * go mod update
  * make dist: remove references to Go

2.29.0 / 2021-11-28
==================

  * stdlib: add use_flake function (#847)
  * docs(direnv.toml) Add config.toml clarification (#831)
  * docs(install): fix macos links (#841)
  * Corrects stdlib link in Ruby docs (#837)
  * stdlib.sh: Fix removal of temp file (#830)
  * install.sh: add aarch64 support
  * Updated conditional for zsh hook to be more forgiving (#808)
  * Add -r flag for matching Git branches with a regexp (#800)
  * Add docs about pipenv (#797)
  * Enable syntax highlights to the quick demo code (#752)
  * Fixed extra quotes for lower alpha characters (#783)
  * Remove noisy warning about PS1 again (#781)

2.28.0 / 2021-03-12
==================

  * Merge pull request #779 from wingrunr21/go_1_16
  * Build for darwin/arm64. Resolves #738
  * Update to go 1.16
  * test: Fix errors for elvish test (#767)
  * tcsh: fix variable escaping (#778)
  * Change DESTDIR to PREFIX in development.md (#774)
  * go: use the /v2 prefix (#765)
  * Relax README's recommendation for nix-direnv (#763)
  * man/direnv.1.md: add FILES section (fix #758) (#759)
  * Add/update fish tests (#754)
  * build(deps): bump golang.org/x/mod from 0.4.0 to 0.4.1 (#749)
  * Fix typo "avaible" in install.sh (#750)
  * docs: improve the use_node documentation

2.27.0 / 2021-01-01
==================

  * fixed fish shell hook to work with eval (#743)
  * dist: remove darwin/386
  * nix: update to nixpkgs@nixos-20.09
  * packaging: stop vendoring the Go code (#739)
  * packaging: change packaging. DESTDIR -> PREFIX, fish hook (#741)

2.26.0 / 2020-12-27
==================

  * updated fish hook support issue (#732)
  * ci: add basic windows CI (#737)
  * test: fix shellcheck usage in ./test/stdlib.bash
  * test: fix use_julia test for NixOS
  * remove dead code: rootDir
  * fix: create temp dir in current working dir for one test (#735)
  * Add `dotenv_if_exists` (#734)
  * stdlib: add watch_dir command (#697)

2.25.2 / 2020-12-12
==================

There was a generation issue in 2.25.1. This release only bumps the version
to do another release.

2.25.1 / 2020-12-11
==================

  * stdlib.go: re-generate (fixes #707)
  * README: remove old Azure badge
  * build(deps): bump golang.org/x/mod from 0.3.0 to 0.4.0 (#730)

2.25.0 / 2020-12-03
==================

  * dist: add linux/arm64 and linux/ppc64
  * Added use_nodenv to stdlib (#727)
  * Fix proposal for  #707, broken direnv compatibility under Windows (#723)
  * fix: layout anaconda <env_name_or_prefix> (#717)
  * Add on_git_branch command to detect whether a specific git branch is checked out (#702)

2.24.0 / 2020-11-15
==================

  * direnv_load: avoid leaking DIRENV_DUMP_FILE_PATH (#715)
  * Add strict_env and unstrict_env (#572)
  * stdlib: add `use_vim` to source local vimrc (#497)
  * stdlib: add source_env_if_exists (#714)
  * Wording (#713)
  * build(deps): bump actions/checkout from v2.3.3 to v2.3.4 (#709)
  * build(deps): bump cachix/install-nix-action from v11 to v12 (#710)
  * Fix XDG_CACHE_HOME path (#711)
  * rc: make file existence check more robust (#706)

2.23.1 / 2020-10-22
==================

  * fix: handle links on Mac when using `allow` (#696)
  * fix: use restored env in exec (#695)
  * stdlib: add basename and dirname from realpath (#693)
  * stdlib.sh: remove tabs
  * dist: compile all the binaries statically

2.23.0 / 2020-10-10
==================

  * stdlib: add source_url function (#562)
  * direnv: add fetchurl command (#686)
  * shell: Update Elvish hook to replace deprecated `explode` (#685)

2.22.1 / 2020-10-06
==================

  * Look for stdlib in DIRENV_CONFIG (#679)
  * stdlib: use Bash 3.0-compatible array expansion (#676)
  * Clarify path to direnv.toml (#678)
  * stdlib/use_julia: fix a bug in parameter substitution for empty or (#667)
  * man: update the layout_go documentation
  * stdlib:  adds GOPATH/bin to PATH (#670)

2.22.0 / 2020-09-01
==================

  * stdlib: use_julia <version> (#666)
  * stdlib: semver_search (#665)
  * direnv-stdlib.1: add layout julia (#661)
  * README: spelling correction (#660)
  * README.md: add shadowenv to similar projects (#659)
  * docs: remove Snap from the installations
  * OSX -> macOS (#655)
  * Update shell_fish.go to use \X for UTF encoding (#584)
  * Change XDG_CONFIG_DIR to XDG_CONFIG_HOME (#641)
  * Streamline core algorithm of export and exec (#636)
  * test: add failure test-case (#637)

2.21.3 / 2020-05-08
==================

  * Replace `direnv expand_path` with pure bash (#631)
  * Fix #594 - write error to fd 3 on Windows (#634)
  * Make direnv hook output work on Windows (#632)
  * Update hook.md to remove ">" typo in Fish instructions (#624)
  * stdlib: `layout go` adds layout dir to GOPATH (#622)
  * direnv-stdlib.1: add layout php (#619)
  * stdlib: add PATH_rm <pattern> [<pattern> ...] (#615)
  * Error handling tuples (#610)
  * Merge pull request #607 from punitagrawal/master
  * test: elvish: Fix evaluation function
  * stdlib.sh: Re-write grep pattern to avoid shell escape
  * man: Escape '.' at the beginning of line to remove manpage warning
  * stdlib: fix direnv_config_dir usage (#601)
  * direnv version: improve error message (#599)
  * README: fix NixOS link in installation.md (#589)
  * stdlib: add direnv_apply_dump <file> (#587)
  * Simplify direnv_load and make it work even when the command crashes. (#568)
  * docs: fix fish installation instruction
  * test: test for utf-8 compatibility
  * config: add [global] section
  * config: add strict_env option
  * config: fix warn_timeout parsing (#582)
  * Github action for releases
  * config: fix the configuration file selection
  * stdlib: fix shellcheck warnings

2.21.2 / 2020-01-28
==================

Making things stable again.

  * stdlib: revert the `set -euo pipefail` change. It was causing too many
    issues for users.
  * direnv allow: fix the allow migration by also creating the parent target
    directory.

2.21.1 / 2020-01-26
==================

Fix release

  * stdlib: fix unused variable in `use node`
  * stdlib: fix unused variable in `source_up`
  * test: add stdlib test skeleton
  * add dist release utility

2.21.0 / 2020-01-25
==================

This is a massive release!

## Highlights

You can now hit Ctrl-C during a long reload in bash and zsh and it will not
loop anymore.

Commands that use `direnv_load` won't fail when there is an output to stdout
anymore (eg: `use_nix`).

Direnv now also loads files from `.config/direnv/lib/*.sh`. This is intended
to be used by third-party tools to augment direnv with their own stdlib
functions.

The `.envrc` is now loaded with `set -euo pipefail`. This will more likely
expose issues with existing `.envrc` files.

## docs

  * Update README.md (#536)
  * Add link to asdf-direnv. (#535)
  * docs: fix invalid link (#533)
  * adds experimental curl based installer (#539)

## commands

  * change where the allow files are being stored
  * direnv status: also show the config
  * direnv exec: improve the error message
  * warn if PS1 is being exported
  * handle SIGINT during export in bash
  * export: display the full RC path instead of a relative one
  * direnv exec: the DIR argument is always required (#493)

## build

  * ci: use GitHub Actions instead of Azure Pipelines
  * staticcheck (#543)
  * use go modules
  * make: handle when /dev/stderr doesn't exist (#491)
  * site: use jekyll to render the website
  * Pin nixpkgs to current NixOS 19.09 channel (#526)

## shells

  * fix elvish hook
  * Use `fish_preexec` hook instead of `fish_prompt` (#512)
  * Use `fish_postexec` to make sure direnv hook executed 'after' the directory has changed when using `cd`.
  * improve zsh hook (#514)

## config.toml

  * rename the configuration from config.toml to direnv.toml (#498)
  * add warn_timeout option. DIRENV_WARN_TIMEOUT is now deprecated.

## stdlib

  * `direnv_load` can now handle stdout outputs
  * stdlib: add layout_julia
  * Handle failing pipenv on empty file and avoid an extra pipenv execution (#510)
  * fix `source_env` behaviour when the file doesn't exists (#487)
  * `watch_file` can now watch multiple files in a single invocation (#524)
  * `layout_python`: prefer venv over virtualenv. Do not export VIRTUAL_ENV if $python_version is unavailable or a virtual environment does not exist/can't be created
  * Adds layout_pyenv (#505)
  * Fix `source_up` docs to explain that search starts in parent directory (#518)
  * fix `path_add` to not leak local variables
  * `layout_pyenv`: support multiple python versions (#525)
  * Add a `direnv_version <version_at_least>` command to check the direnv
    version.
  * `dotenv`: handle undefined variables
  * source files from `.config/direnv/lib/*.sh`
  * stdlib: set `-euo pipefail`

2.20.1 / 2019-03-31
==================

  * ci: try to fix releases

2.20.0 / 2019-03-31
==================

  * CHANGE: Use source instead of eval on fish hook
  * DOC: Remove duplicate build badge (#465)
  * DOC: add note about auth (#463)
  * DOC: change nixos link (#460)
  * FIX: Corrects reverse patching when using exec cmd. (#466)
  * FIX: Perform stricter search for existing Anaconda environments (#462)
  * FIX: arity mismatch for elvish (#482)
  * FIX: avoid reloading on each prompt after error (#468)
  * FIX: improve bash hook handling of empty PROMPT_COMMAND (#473)
  * FIX: improved the tests for bash, zsh, fish and tcsh (#469)
  * MISC: migrated from Travis CI to Azure Pipelines (#484)

2.19.2 / 2019-02-09
==================

  * FIX: file_times: check Stat and Lstat (#457)

2.19.1 / 2019-01-31
==================

  * FIX: watched files now handle symlinks properly. Thanks @grahamc! #452

2.19.0 / 2019-01-11
==================

  * NEW: add support for .env variable expansion. Thanks to @hakamadare!

2.18.2 / 2018-11-23
==================

  * make: generate direnv.exe on windows (#417)

2.18.1 / 2018-11-22
==================

  * travis: fix the release process

2.18.0 / 2018-11-22
==================

A lot of changes!

  * stdlib: add DIRENV_IN_ENVRC (#414)
  * Fix typo in readme. (#412)
  * Merge pull request #407 from zimbatm/direnv-dump-shell
  * direnv dump can now dump to arbitrary shells
  * add a new "gzenv" shell
  * move gzenv into new package
  * shell: introduce a dump capability
  * cleanup the shells
  * Add alias '--version' to version command. Closes #377. (#404)
  * Correctes spelling of openSUSE (#403)
  * testing: elvish 0.12 is released now (#402)
  * Merge pull request #397 from zimbatm/readme-packaging-status
  * README: add packaging status badge
  * README: remove equinox installation
  * direnv show_dump: new command to debug encoded env (#395)
  * Document possibility to unset vars (#392)
  * stdlib: fix typo
  * go dep: update Gopkg.lock
  * make: don't make shfmt a dependency
  * Avoid to add unnecessary trailing semicolon character (#384)
  * add asdf to the list of known projects
  * stdlib.go: re-generate
  * Add PHP layout to stdlib (#346)
  * make: fix formatting
  * README: add build status badge
  * Overhaul the build system (#375)
  * stdlib, layout_pipenv: handle `$PIPENV_PIPFILE` (#371)
  * README: improve the source build instructions

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
 * FIX: man changes, replaced abandoned ronn by md2man
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

