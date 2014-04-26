#!/usr/bin/env bash

set -e

cd `dirname $0`
TEST_DIR=$PWD
export PATH=`dirname $TEST_DIR`:$PATH

# Reset the direnv loading if any
export DIRENV_CONFIG=$PWD
unset DIRENV_BASH
unset DIRENV_DIR
unset DIRENV_MTIME
unset DIRENV_DIFF

direnv_eval() {
  eval `direnv export bash`
}

test_start() {
  cd "$TEST_DIR/scenarios/$1"
  direnv allow
  echo "## Testing $1 ##"
}

test_stop() {
  cd $TEST_DIR
  direnv_eval
}

### RUN ###

direnv allow || true
direnv_eval

test_start base
  direnv_eval
  test "$HELLO" = "world"

  MTIME=$DIRENV_MTIME
  direnv_eval
  test "$MTIME" = "$DIRENV_MTIME"

  sleep 1

  touch .envrc
  direnv_eval
  test "$MTIME" != "$DIRENV_MTIME"

  cd ..
  direnv_eval
  echo "${HELLO}"
  test -z "${HELLO}"
test_stop

test_start inherit
  direnv_eval
  test "$HELLO" = "world"
test_stop

test_start "ruby-layout"
  direnv_eval
  test "$GEM_HOME" != ""
test_stop

# Make sure directories with spaces are fine
test_start "space dir"
  direnv_eval
  test "$SPACE_DIR" = "true"
test_stop

test_start "child-env"
  direnv_eval
  test "$PARENT_PRE" = "1"
  test "$CHILD" = "1"
  test "$PARENT_POST" = "1"
  test -z "$REMOVE_ME"
test_stop

test_start "special-vars"
  export DIRENV_BASH=`which bash`
  export DIRENV_CONFIG=foobar
  direnv_eval || true
  test -n "$DIRENV_BASH"
  test "$DIRENV_CONFIG" = "foobar"
  unset DIRENV_BASH
  unset DIRENV_CONFIG
test_stop

test_start "empty-var"
  direnv_eval
  test "${FOO-unset}" != "unset"
  test "${FOO}" = ""
test_stop

test_start "empty-var-unset"
  export FOO=""
  direnv_eval
  test "${FOO-unset}" == "unset"
  unset FOO
test_stop

# Context: foo/bar is a symlink to ../baz. foo/ contains and .envrc file
# BUG: foo/bar is resolved in the .envrc execution context and so can't find
#      the .envrc file.
#
# Apparently, the CHDIR syscall does that so I don't know how to work around
# the issue.
#
# test_start "symlink-bug"
#   cd foo/bar
#   direnv_eval
# test_stop

# Pending: test that the mtime is looked on the original file
# test_start "utils"
#   LINK_TIME=`direnv file-mtime link-to-somefile`
#   touch somefile
#   NEW_LINK_TIME=`direnv file-mtime link-to-somefile`
#   test "$LINK_TIME" = "$NEW_LINK_TIME"
# test_stop
