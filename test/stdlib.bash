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

# test find_up
(
  load_stdlib
  path=$(find_up "README.md")
  assert_eq "$path" "$root/README.md"
)

# test source_up
(
  load_stdlib
  cd scenarios/inherited
  source_up
)
