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

test_name direnv_apply_dump
(
  tmpfile=$(mktemp)
  # shellcheck disable=SC2317
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
    echo "#!/usr/bin/env bash
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

test_name uv
(
  load_stdlib

  if ! has uv; then
    echo "WARN: uv not found, skipping..."
    return
  fi

  tmpdir=$(mktemp -d)
  trap 'rm -rf $tmpdir' EXIT

  cd "$tmpdir"

  # Create a virtual environment
  layout uv 3.12

  # Check if VIRTUAL_ENV is set correctly
  [[ $VIRTUAL_ENV = "$tmpdir/.venv" ]]

  # Check if Python is available and its version
  [[ $(python --version) = "Python 3.12"* ]]

  # Check if UV_ACTIVE is set
  [[ $UV_ACTIVE = "1" ]]

  # Check if UV_PROJECT_ENVIRONMENT is set correctly
  [[ $UV_PROJECT_ENVIRONMENT = "$tmpdir/.venv" ]]
)

test_name uvp
(
  load_stdlib

  if ! has uv; then
    echo "WARN: uv not found, skipping..."
    return
  fi

  tmpdir=$(mktemp -d)
  trap 'rm -rf $tmpdir' EXIT

  cd "$tmpdir"

  # Create a virtual environment
  layout uvp 3.12

  # Check if VIRTUAL_ENV is set correctly
  [[ $VIRTUAL_ENV = "$tmpdir/.venv" ]]

  # Check if Python is available and its version
  [[ $(python --version) = "Python 3.12"* ]]

  # Check if pyproject.toml exists
  [[ -f pyproject.toml ]]

  # README should exist
  [[ -f README.md ]]

  # Check if UV_ACTIVE is set
  [[ $UV_ACTIVE = "1" ]]

  # Check if UV_PROJECT_ENVIRONMENT is set correctly
  [[ $UV_PROJECT_ENVIRONMENT = "$tmpdir/.venv" ]]
)

test_name uv_no_version_twice
(
  load_stdlib

  if ! has uv; then
    echo "WARN: uv not found, skipping..."
    return
  fi

  tmpdir=$(mktemp -d)
  trap 'rm -rf $tmpdir' EXIT

  cd "$tmpdir"

  # Create a virtual environment with no version specified
  layout uv

  # Create a virtual environment with no version specified again
  second_run_output=$(layout uv 2>&1)
  [[ "${second_run_output}" != *"No virtual environment exists"* ]]
)

test_name uvp_no_version_twice
(
  load_stdlib

  tmpdir=$(mktemp -d)
  trap 'rm -rf $tmpdir' EXIT

  cd "$tmpdir"

  layout uvp

  second_run_output=$(layout uvp 2>&1)
  [[ "${second_run_output}" != *"wrong python version"* ]]
  [[ "${second_run_output}" != *"No uv project exists"* ]]
)

test_name uv_version_switch
(
  load_stdlib

  tmpdir=$(mktemp -d)
  trap 'rm -rf $tmpdir' EXIT

  cd "$tmpdir"

  layout uv 3.12

  same_version_output_message=$(layout uv 3.12 2>&1)
  [[ "${same_version_output_message}" != *"No virtual environment exists"* ]]

  different_version_output_message=$(layout uv 3.11 2>&1)
  [[ "${different_version_output_message}" = *"No virtual environment exists"* ]]
)

test_name uvp_version_mismatch
(
  load_stdlib

  tmpdir=$(mktemp -d)
  trap 'rm -rf $tmpdir' EXIT

  cd "$tmpdir"

  layout uvp 3.12

  error_output=$(layout uvp 3.11 2>&1 >/dev/null)

  # Make sure error message contains the string "wrong python versionk"
  [[ "${error_output#*'wrong python version'}" != "$error_output" ]]
)

test_name uvp_additional_args
(
  load_stdlib

  tmpdir=$(mktemp -d)
  trap 'rm -rf $tmpdir' EXIT

  cd "$tmpdir"

  layout uvp -- --no-readme
  # Make sure README.md does not exist
  [[ ! -f README.md ]]
)

# test strict_env and unstrict_env
./strict_env_test.bash

echo OK
