#!/usr/bin/env bash

set -e

export DIRENV_DEBUG=1

cd `dirname $0`
TEST_DIR=$PWD
export PATH=$PWD/../bin:$PATH

direnv_eval() {
  eval `direnv export`
}

test_start() {
  cd "$TEST_DIR/scenarios/$1"
  echo "## Testing $1 ##"
  if test "${2-noparam}" == "noparam"; then
    # no second argument, auto-enable environment
    direnv enable
  fi
}

test_stop() {
  direnv disable
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

test_start "evil" --no-auto-enable
  # Environment is disabled
  test -z "${DIRENV_HAVE_EVIL-}"
  
  sleep 1; direnv enable; direnv_eval
  
  test -n "${DIRENV_HAVE_EVIL}"
test_stop
