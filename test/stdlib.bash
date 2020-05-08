#!/usr/bin/env bash
set -euo pipefail

# always execute relative to here
cd "$(dirname "$0")"

# add the built direnv to the path
root=$(cd .. && pwd -P)
export PATH=$root:$PATH

load_stdlib() {
  # shellcheck disable=SC1090
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

echo OK
