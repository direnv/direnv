#! /usr/bin/env bash

# Note: this is _explicitly_ not setting `set -euo pipefail`
# because we are testing functions that configure that.

declare base expected actual
base="${TMPDIR:-/tmp}/$(basename "$0").$$"
expected="$base".expected
actual="$base".actual

declare -i success
success=0

# always execute relative to here
cd "$(dirname "$0")" || exit 1

# add the built direnv to the path
root=$(cd .. && pwd -P)
export PATH=$root:$PATH

load_stdlib() {
  # shellcheck disable=SC1090
  source "$root/stdlib.sh"
}

test_fail() {
  echo "FAILED $*: expected is not actual"
  exit 1
}

test_strictness() {
  local step
  step="$1"
  echo "$2" > "$expected"

  set -o | grep 'errexit\|nounset\|pipefail' > "$actual"
  diff -u "$expected" "$actual" || test_fail "$step"
  (( success += 1 ))
}

load_stdlib

test_strictness 'after source' $'errexit        	off
nounset        	off
pipefail       	off'

strict_env
test_strictness 'after strict_env' $'errexit        	on
nounset        	on
pipefail       	on'

unstrict_env echo HELLO > /dev/null
test_strictness "after unstrict_env with command" $'errexit        	on
nounset        	on
pipefail       	on'

strict_env echo HELLO > /dev/null
test_strictness "after strict_env with command" $'errexit        	on
nounset        	on
pipefail       	on'

unstrict_env
test_strictness 'after unstrict_env' $'errexit        	off
nounset        	off
pipefail       	off'

echo "OK ($success tests)"
