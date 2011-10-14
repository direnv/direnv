#!/usr/bin/env bash

set -e
set -u

cd `dirname $0`
export PATH=$PWD/../bin:$PATH
eval `direnv export`

test_start() {
  pushd "scenarios/$1" > /dev/null
  echo "## Testing $1 ##"
  time direnv dump > /dev/null
  foo=`direnv dump`
  time direnv diff "$foo"
  time direnv export > /dev/null 2>&1
  eval `direnv export` 2>/dev/null
}

test_stop() {
  popd > /dev/null
  eval `direnv export` 2>/dev/null
}

test_start base
test "$HELLO" = "world"
test_stop

test_start inherit
test "$HELLO" = "world"
test_stop

test_start "ruby-layout"
test "$GEM_HOME" != ""
test_stop

# Make sure directories with spaces are fine
test_start "space dir"
test_stop
