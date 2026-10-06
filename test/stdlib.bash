#!/usr/bin/env bash
set -euo pipefail

# always execute relative to here
cd "$(dirname "$0")"

# add the built direnv to the path
root=$(cd .. && pwd -P)
export PATH=$root:$PATH

load_stdlib() {
  # shellcheck disable=SC1090,SC1091
  source "$root/stdlib.sh"
}

assert_eq() {
  if [[ $1 != "$2" ]]; then
    echo "expected '$1' to equal '$2'"
    return 1
  fi
}

test_name() {
  echo "--- $*"
}

test_name dotenv
(
  load_stdlib

  workdir=$(mktemp -d)
  trap 'rm -rf "$workdir"' EXIT

  cd "$workdir"

  # Try to source a file that doesn't exist - should not succeed
  dotenv .env.non_existing_file && return 1

  # Try to source a file that exists
  echo "export FOO=bar" > .env
  dotenv .env
  [[ $FOO = bar ]]
)

test_name dotenv_if_exists
(
  load_stdlib

  workdir=$(mktemp -d)
  trap 'rm -rf "$workdir"' EXIT

  cd "$workdir"

  # Try to source a file that doesn't exist - should succeed
  dotenv_if_exists .env.non_existing_file  || return 1

  # Try to source a file that exists
  echo "export FOO=bar" > .env
  dotenv_if_exists .env
  [[ $FOO = bar ]]
)

test_name "dotenv / dotenv_if_exists (named pipe / FIFO)"
(
  load_stdlib

  workdir=$(mktemp -d)
  writer_pid=""
  # shellcheck disable=SC2064
  trap '[[ -n $writer_pid ]] && kill "$writer_pid" 2>/dev/null; rm -rf "$workdir"' EXIT

  # Reap the background writer without ever hanging the suite: a writer blocked on
  # opening the FIFO for write (e.g. if the reader never opened it) is bounded-waited,
  # then killed. The EXIT trap is the final backstop.
  reap_writer() {
    for _ in $(seq 1 50); do kill -0 "$writer_pid" 2>/dev/null || break; sleep 0.1; done
    kill "$writer_pid" 2>/dev/null || true
    wait "$writer_pid" 2>/dev/null || true
    writer_pid=""
  }

  cd "$workdir"

  # Positive control: watch_file is functional in this harness, so the no-watch
  # assertions below are meaningful (watching a regular file changes DIRENV_WATCHES).
  before=${DIRENV_WATCHES:-}
  : > regular.txt
  watch_file regular.txt
  [[ ${DIRENV_WATCHES:-} != "$before" ]] || { echo "watch_file did not record a regular file"; return 1; }

  # A ".env" provided as a named pipe (FIFO) - as mounted by secrets managers
  # like 1Password Environments to inject secrets on read, without writing the
  # secret contents to disk. The writer blocks until dotenv opens the pipe.
  mkfifo fifo.env

  # dotenv loads the FIFO ...
  before=${DIRENV_WATCHES:-}
  ( echo "export FOO=bar" > fifo.env ) &
  writer_pid=$!
  dotenv fifo.env
  reap_writer
  [[ $FOO = bar ]]
  # ... and must NOT add the FIFO to the watch list: a FIFO's mtime changes on
  # every read, which would otherwise force a reload on every prompt.
  assert_eq "${DIRENV_WATCHES:-}" "$before"

  # dotenv_if_exists shares the same gate, so it must load a FIFO too.
  unset FOO
  before=${DIRENV_WATCHES:-}
  ( echo "export FOO=baz" > fifo.env ) &
  writer_pid=$!
  dotenv_if_exists fifo.env
  reap_writer
  [[ $FOO = baz ]]
  assert_eq "${DIRENV_WATCHES:-}" "$before"
)

test_name find_up
(
  load_stdlib
  path=$(find_up "README.md")
  assert_eq "$path" "$root/README.md"
)

test_name source_up
(
  load_stdlib
  cd scenarios/inherited
  source_up
)

test_name eval_propagates_failure
(
  workdir=$(mktemp -d)
  trap 'rm -rf "$workdir"' EXIT
  cd "$workdir"

  # shellcheck disable=SC2016
  printf '#!/usr/bin/env bash\n[[ $1 == "$FAIL_CMD" ]] && exit 1\nexec %q "$@"\n' \
    "$(command -v direnv)" >fake-direnv
  chmod +x fake-direnv
  echo "FOO=bar" >.env
  direnv dump >env.dump

  # a separate process, because errexit is ignored inside `if` and `&&`
  run_failing() {
    FAIL_CMD=$1 bash -euo pipefail -c 'source "$1"; direnv=$2; shift 2; "$@"' \
      _ "$root/stdlib.sh" "$workdir/fake-direnv" "${@:2}"
  }

  for t in "dotenv dotenv .env" "dotenv dotenv_if_exists .env" \
    "check-required require_allowed .env" "watch watch_file .env" \
    "watch-dir watch_dir ." "apply_dump direnv_apply_dump env.dump"; do
    # shellcheck disable=SC2086
    if run_failing $t 2>/dev/null; then
      echo "failure of 'direnv ${t%% *}' was ignored"
      return 1
    fi
  done

  load_stdlib
  echo "result=kept" >.env
  dotenv .env
  assert_eq "$result" kept
)

test_name direnv_apply_dump
(
  tmpfile=$(mktemp)
  # shellcheck disable=SC2329
  cleanup() { rm "$tmpfile"; }
  trap cleanup EXIT

  load_stdlib
  FOO=bar direnv dump > "$tmpfile"
  direnv_apply_dump "$tmpfile"
  assert_eq "$FOO" bar
)

test_name PATH_rm
(
  load_stdlib

  export PATH=/usr/local/bin:/home/foo/bin:/usr/bin:/home/foo/.local/bin
  PATH_rm '/home/foo/*'

  assert_eq "$PATH" /usr/local/bin:/usr/bin
)

test_name path_rm
(
  load_stdlib

  somevar=/usr/local/bin:/usr/bin:/home/foo/.local/bin
  path_rm somevar '/home/foo/*'

  assert_eq "$somevar" /usr/local/bin:/usr/bin
)

test_name expand_path
(
  load_stdlib
  tmpdir=$(mktemp -d)
  trap 'rm -rf $tmpdir' EXIT

  cd "$tmpdir"
  ret=$(expand_path ./bar)

  assert_eq "$ret" "$tmpdir/bar"
)

test_name semver_search
(
  load_stdlib
  versions=$(mktemp -d)
  trap 'rm -rf $versions' EXIT

  mkdir "$versions/program-1.4.0"
  mkdir "$versions/program-1.4.1"
  mkdir "$versions/program-1.5.0"
  mkdir "$versions/1.6.0"

  assert_eq "$(semver_search "$versions" "program-" "1.4.0")" "1.4.0"
  assert_eq "$(semver_search "$versions" "program-" "1.4")"   "1.4.1"
  assert_eq "$(semver_search "$versions" "program-" "1")"     "1.5.0"
  assert_eq "$(semver_search "$versions" "program-" "1.8")"   ""
  assert_eq "$(semver_search "$versions" "" "1.6")"           "1.6.0"
  assert_eq "$(semver_search "$versions" "program-" "")"      "1.5.0"
  assert_eq "$(semver_search "$versions" "" "")"              "1.6.0"
)

test_name use_julia
(
  load_stdlib
  JULIA_VERSIONS=$(TMPDIR=. mktemp -d -t tmp.XXXXXXXXXX)
  trap 'rm -rf $JULIA_VERSIONS' EXIT

  test_julia() {
    version_prefix="$1"
    version="$2"
    # Fake the existence of a julia binary
    julia=$JULIA_VERSIONS/$version_prefix$version/bin/julia
    mkdir -p "$(dirname "$julia")"
    echo "#!$(command -v bash)
    echo \"test-julia $version\"" > "$julia"
    chmod +x "$julia"
    # Locally disable set -u (see https://github.com/direnv/direnv/pull/667)
    if ! [[ "$(set +u; use julia "$version" 2>&1)" =~ Successfully\ loaded\ test-julia\ $version ]]; then
      return 1
    fi
  }

  # Default JULIA_VERSION_PREFIX
  unset JULIA_VERSION_PREFIX
  test_julia "julia-" "1.0.0"
  test_julia "julia-" "1.1"
  # Custom JULIA_VERSION_PREFIX
  JULIA_VERSION_PREFIX="jl-"
  test_julia "jl-"    "1.2.0"
  test_julia "jl-"    "1.3"
  # Empty JULIA_VERSION_PREFIX
  # shellcheck disable=SC2034
  JULIA_VERSION_PREFIX=
  test_julia ""    "1.4.0"
  test_julia ""    "1.5"
)

test_name use_guix
(
  load_stdlib
  workdir=$(mktemp -d)
  trap 'rm -rf "$workdir"' EXIT
  cd "$workdir"

  # shellcheck disable=SC2329
  guix() { echo "mode=search"; }
  # shellcheck disable=SC2329
  direnv_load() { mode=load; }
  # shellcheck disable=SC2329
  watch_file() { :; }

  for args in --container -C -NC --emulate-fhs -F "-m m.scm -CF"; do
    mode=
    # shellcheck disable=SC2086
    use_guix $args
    assert_eq "$mode" search
  done
  for args in --file=Config.scm "-f Foo.scm" --development hello; do
    mode=
    # shellcheck disable=SC2086
    use_guix $args
    assert_eq "$mode" load
  done
)

test_name source_env_if_exists
(
  load_stdlib

  workdir=$(mktemp -d)
  trap 'rm -rf "$workdir"' EXIT

  cd "$workdir"

  # Try to source a file that doesn't exist
  source_env_if_exists non_existing_file

  # Try to source a file that exists
  echo "export FOO=bar" > existing_file
  source_env_if_exists existing_file
  [[ $FOO = bar ]]

  # Expect correct path being logged
  export HOME=$workdir
  output="$(source_env_if_exists existing_file 2>&1 > /dev/null)"
  [[ "${output#*'loading ~/existing_file'}" != "$output" ]]
)

test_name env_vars_required
(
  load_stdlib

  export FOO=1
  env_vars_required FOO

  # these should all fail
  # shellcheck disable=SC2034
  BAR=1
  export BAZ=
  output="$(env_vars_required BAR BAZ MISSING 2>&1 > /dev/null || echo "--- result: $?")"

  [[ "${output#*'--- result: 1'}" != "$output" ]]
  [[ "${output#*'BAR is required'}" != "$output" ]]
  [[ "${output#*'BAZ is required'}" != "$output" ]]
  [[ "${output#*'MISSING is required'}" != "$output" ]]
)


test_name require_allowed_security
(
  load_stdlib
  set +e

  # Test that absolute paths are rejected
  output="$(require_allowed /etc/passwd 2>&1)"
  result=$?
  [[ $result -eq 1 ]]
  [[ "${output#*'path must be relative'}" != "$output" ]]

  # Test that parent traversal paths are rejected
  output="$(require_allowed ../etc/passwd 2>&1)"
  result=$?
  [[ $result -eq 1 ]]
  [[ "${output#*'must not contain'}" != "$output" ]]

  # Test that paths with .. in the middle are rejected
  output="$(require_allowed foo/../bar 2>&1)"
  result=$?
  [[ $result -eq 1 ]]
  [[ "${output#*'must not contain'}" != "$output" ]]
)

test_name global_lib_noglob
(
  # Regression test for https://github.com/direnv/direnv/issues/1610
  # With pathname expansion disabled (set -f / noglob) and no global library
  # files present, direnv must not source the literal "*.sh" glob and emit a
  # spurious missing-file diagnostic for the optional lib directory.
  workdir=$(mktemp -d)
  trap 'rm -rf "$workdir"' EXIT

  mkdir -p "$workdir/home" "$workdir/project" "$workdir/config/direnv"
  echo "export DIRENV_NOGLOB_TEST=1" > "$workdir/project/.envrc"

  HOME="$workdir/home" XDG_CONFIG_HOME="$workdir/config" \
    DIRENV_CONFIG="$workdir/config/direnv" \
    direnv allow "$workdir/project/.envrc" >/dev/null 2>&1

  output="$(
    cd "$workdir/project"
    env HOME="$workdir/home" XDG_CONFIG_HOME="$workdir/config" \
      DIRENV_CONFIG="$workdir/config/direnv" \
      SHELLOPTS=braceexpand:hashall:interactive-comments:noglob \
      direnv export bash 2>&1 >/dev/null
  )"

  if [[ "$output" == *"lib/*.sh"* ]]; then
    echo "unexpected missing-library diagnostic under noglob: $output"
    return 1
  fi
)

# test strict_env and unstrict_env
./strict_env_test.bash

echo OK
