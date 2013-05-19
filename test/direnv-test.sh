#!/usr/bin/env bash

set -e

cd `dirname $0`
TEST_DIR=$PWD
export PATH=`dirname $TEST_DIR`:$PATH

# Reset the direnv loading if any
unset DIRENV_BACKUP
unset DIRENV_DIR
unset DIRENV_MTIME

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

# Pending: test that the mtime is looked on the original file
# test_start "utils"
#   LINK_TIME=`direnv file-mtime link-to-somefile`
#   touch somefile
#   NEW_LINK_TIME=`direnv file-mtime link-to-somefile`
#   test "$LINK_TIME" = "$NEW_LINK_TIME"
# test_stop
