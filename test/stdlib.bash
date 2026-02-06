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
  # shellcheck disable=SC2030
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

test_name load_netrc
(
  load_stdlib

  workdir=$(mktemp -d)
  trap 'rm -rf "$workdir"' EXIT

  cd "$workdir"

  # Comprehensive netrc file covering: comments, macdef blocks, multi-line,
  # single-line, multiple machines, and empty lines
  cat > .netrc << 'EOF'
# comment at top
macdef init
login fake_from_macro
password fake_from_macro

machine api.example.com
  # inline comment
  login testuser
  password testpass123

machine single.example.com login singleuser password singlepass

machine reversed.example.com
password reversedpass
login reverseduser

machine withaccount.example.com
login acctuser
password acctpass
account acctvalue

default
login defaultuser
password defaultpass
EOF

  # Test: multi-line format
  load_netrc --file .netrc api.example.com MY_USER MY_PASS 2>/dev/null
  assert_eq "$MY_USER" "testuser"
  assert_eq "$MY_PASS" "testpass123"

  # Test: single-line format
  unset MY_USER MY_PASS
  load_netrc --file .netrc single.example.com MY_USER MY_PASS 2>/dev/null
  assert_eq "$MY_USER" "singleuser"
  assert_eq "$MY_PASS" "singlepass"

  # Test: password before login (reversed order)
  unset MY_USER MY_PASS
  load_netrc --file .netrc reversed.example.com MY_USER MY_PASS 2>/dev/null
  assert_eq "$MY_USER" "reverseduser"
  assert_eq "$MY_PASS" "reversedpass"

  # Test: account variable
  unset MY_USER MY_PASS MY_ACCOUNT
  load_netrc --file .netrc withaccount.example.com MY_USER MY_PASS MY_ACCOUNT 2>/dev/null
  assert_eq "$MY_USER" "acctuser"
  assert_eq "$MY_PASS" "acctpass"
  assert_eq "$MY_ACCOUNT" "acctvalue"

  # Test: account variable not required when not requested
  unset MY_USER MY_PASS
  load_netrc --file .netrc api.example.com MY_USER MY_PASS 2>/dev/null
  assert_eq "$MY_USER" "testuser"
  assert_eq "$MY_PASS" "testpass123"

  # Test: account variable requested but not present should fail
  unset MY_USER MY_PASS MY_ACCOUNT
  if load_netrc --file .netrc api.example.com MY_USER MY_PASS MY_ACCOUNT 2>/dev/null; then
    echo "Expected load_netrc to fail when account_var requested but no account in netrc"
    return 1
  fi

  # Test: --file=<path> form
  unset MY_USER MY_PASS
  load_netrc --file=.netrc api.example.com MY_USER MY_PASS 2>/dev/null
  assert_eq "$MY_USER" "testuser"
  assert_eq "$MY_PASS" "testpass123"

  # Test: load default entry explicitly
  unset MY_USER MY_PASS
  load_netrc --file .netrc default MY_USER MY_PASS 2>/dev/null
  assert_eq "$MY_USER" "defaultuser"
  assert_eq "$MY_PASS" "defaultpass"

  # Test: default entry should not be used for unknown machines
  unset MY_USER MY_PASS
  if load_netrc --file .netrc unknown.example.com MY_USER MY_PASS 2>/dev/null; then
    echo "Expected load_netrc to fail for machine only in default block"
    return 1
  fi

  # Test: default ~/.netrc path
  # shellcheck disable=SC2031
  cp .netrc "$HOME/.netrc"
  unset MY_USER MY_PASS
  load_netrc api.example.com MY_USER MY_PASS 2>/dev/null
  assert_eq "$MY_USER" "testuser"
  assert_eq "$MY_PASS" "testpass123"

  # Test: missing required argument should fail
  if load_netrc --file .netrc "" MY_USER MY_PASS 2>/dev/null; then
    echo "Expected load_netrc to fail for missing machine argument"
    return 1
  fi

  # Test: missing netrc file should fail
  if load_netrc --file missing.netrc api.example.com MY_USER MY_PASS 2>/dev/null; then
    echo "Expected load_netrc to fail for missing netrc file"
    return 1
  fi

  # Test: --file without path should fail
  if load_netrc --file 2>/dev/null; then
    echo "Expected load_netrc to fail for --file without path"
    return 1
  fi

  # Test: missing machine in file should fail
  unset MY_USER MY_PASS
  if load_netrc --file .netrc nonexistent.com MY_USER MY_PASS 2>/dev/null; then
    echo "Expected load_netrc to fail for nonexistent machine"
    return 1
  fi

  # Test: only login (no password) should fail
  cat > .netrc_nopw << 'EOF'
machine nopw.example.com
  login nopwuser
EOF
  if load_netrc --file .netrc_nopw nopw.example.com MY_USER MY_PASS 2>/dev/null; then
    echo "Expected load_netrc to fail for machine with no password"
    return 1
  fi

  # Test: macdef inside a machine block should not leak into parsing
  cat > .netrc_macdef << 'EOF'
machine other.example.com
login otheruser
password otherpass
macdef upload
put myfile

machine target.example.com
login targetuser
password targetpass
EOF
  unset MY_USER MY_PASS
  load_netrc --file .netrc_macdef target.example.com MY_USER MY_PASS 2>/dev/null
  assert_eq "$MY_USER" "targetuser"
  assert_eq "$MY_PASS" "targetpass"

  # Test: keyword-like values (login is "machine", password is "login")
  cat > .netrc_tricky << 'EOF'
machine tricky.example.com login machine password login
EOF
  unset MY_USER MY_PASS
  load_netrc --file .netrc_tricky tricky.example.com MY_USER MY_PASS 2>/dev/null
  assert_eq "$MY_USER" "machine"
  assert_eq "$MY_PASS" "login"

  # Test: keyword and value split across lines
  cat > .netrc_split << 'EOF'
machine
split.example.com
login
splituser
password
splitpass
EOF
  unset MY_USER MY_PASS
  load_netrc --file .netrc_split split.example.com MY_USER MY_PASS 2>/dev/null
  assert_eq "$MY_USER" "splituser"
  assert_eq "$MY_PASS" "splitpass"
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

# test strict_env and unstrict_env
./strict_env_test.bash

echo OK
