# Test script for Bourne-shell extensions. Set TARGET_SHELL
# to the shell to be tested (bash, zsh, etc) before sourcing it.
if [[ -z "$TARGET_SHELL" ]]; then
  echo "TARGET_SHELL variable not set"
  exit 1
fi

set -e

cd "$(dirname "$0")"
TEST_DIR=$PWD
export XDG_CONFIG_HOME=${TEST_DIR}/config
export XDG_DATA_HOME=${TEST_DIR}/data
PATH=$(dirname "$TEST_DIR"):$PATH
export PATH

# Reset the direnv loading if any
export DIRENV_CONFIG=$PWD
unset DIRENV_BASH
unset DIRENV_DIR
unset DIRENV_MTIME
unset DIRENV_WATCHES
unset DIRENV_DIFF

mkdir -p "${XDG_CONFIG_HOME}/direnv"
touch "${XDG_CONFIG_HOME}/direnv/direnvrc"

has() {
  type -P "$1" &>/dev/null
}

direnv_eval() {
  eval "$(direnv export "$TARGET_SHELL")"
}

test_start() {
  cd "$TEST_DIR/scenarios/$1"
  direnv allow
  if [[ "$DIRENV_DEBUG" == "1" ]]; then
    echo
  fi
  echo "## Testing $1 ##"
  if [[ "$DIRENV_DEBUG" == "1" ]]; then
    echo
  fi
}

test_stop() {
  cd "$TEST_DIR"
  direnv_eval
}

test_eq() {
  if [[ "$1" != "$2" ]]; then
    echo "FAILED: '$1' == '$2'"
    exit 1
  fi
}

test_neq() {
  if [[ "$1" == "$2" ]]; then
    echo "FAILED: '$1' != '$2'"
    exit 1
  fi
}

### RUN ###

direnv allow || true
direnv_eval

test_start base
  echo "Setting up"
  direnv_eval
  test_eq "$HELLO" "world"

  WATCHES=$DIRENV_WATCHES

  echo "Reloading (should be no-op)"
  direnv_eval
  test_eq "$WATCHES" "$DIRENV_WATCHES"

  sleep 1

  echo "Updating envrc and reloading (should reload)"
  touch .envrc
  direnv_eval
  test_neq "$WATCHES" "$DIRENV_WATCHES"

  echo "Leaving dir (should clear env set by dir's envrc)"
  cd ..
  direnv_eval
  echo "${HELLO}"
  test -z "${HELLO}"
test_stop

test_start inherit
  cp ../base/.envrc ../inherited/.envrc
  direnv_eval
  echo "HELLO should be world:" "$HELLO"

  sleep 1
  echo "export HELLO=goodbye" > ../inherited/.envrc
  direnv_eval
  test_eq "$HELLO" "goodbye"
test_stop

if has ruby; then
  test_start "ruby-layout"
    direnv_eval
    test_neq "$GEM_HOME" ""
  test_stop
fi

# Make sure directories with spaces are fine
test_start "space dir"
  direnv_eval
  test_eq "$SPACE_DIR" "true"
test_stop

test_start "child-env"
  direnv_eval
  test_eq "$PARENT_PRE" "1"
  test_eq "$CHILD" "1"
  test_eq "$PARENT_POST" "1"
  test -z "$REMOVE_ME"
test_stop

test_start "special-vars"
  export DIRENV_BASH=$(command -v bash)
  export DIRENV_CONFIG=foobar
  direnv_eval || true
  test -n "$DIRENV_BASH"
  test_eq "$DIRENV_CONFIG" "foobar"
  unset DIRENV_BASH
  unset DIRENV_CONFIG
test_stop

test_start "dump"
  direnv_eval
  test_eq "$LS_COLORS" "*.ogg=38;5;45:*.wav=38;5;45"
  test_eq "$THREE_BACKSLASHES" '\\\'
  test_eq "$LESSOPEN" "||/usr/bin/lesspipe.sh %s"
test_stop

test_start "empty-var"
  direnv_eval
  test_neq "${FOO-unset}" "unset"
  test_eq "${FOO}" ""
test_stop

test_start "empty-var-unset"
  export FOO=""
  direnv_eval
  test_eq "${FOO-unset}" "unset"
  unset FOO
test_stop

test_start "in-envrc"
  direnv_eval
  set +e
  ./test-in-envrc
  es=$?
  set -e
  test_eq "$es" "1"
test_stop

test_start "missing-file-source-env"
  direnv_eval
test_stop

test_start "symlink-changed"
  # when using a symlink, reload if the symlink changes, or if the
  # target file changes.
  ln -fs ./state-A ./symlink
  direnv_eval
  test_eq "${STATE}" "A"
  sleep 1

  ln -fs ./state-B ./symlink
  direnv_eval
  test_eq "${STATE}" "B"
test_stop

test_start "utf-8"
  direnv_eval
  test_eq "${UTFSTUFF}" "♀♂"
test_stop

test_start "hash"
  direnv_eval
  test_eq "$HELLO" "world"
  test_eq "$TEST" "1234"
  test -z "$UNSET"
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
