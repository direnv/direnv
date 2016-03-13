#!/usr/bin/env tcsh -e -x

# set -e

cd `dirname $0`
setenv TEST_DIR $PWD
setenv PATH `dirname $TEST_DIR`:$PATH

# Reset the direnv loading if any
setenv DIRENV_CONFIG $PWD
unsetenv DIRENV_BASH
unsetenv DIRENV_DIR
unsetenv DIRENV_MTIME
unsetenv DIRENV_WATCHES
unsetenv DIRENV_DIFF

# direnv_eval() {
#   eval `direnv export bash`
# }
alias direnv_eval 'eval `direnv export tcsh`'

# test_start() {
#   cd "$TEST_DIR/scenarios/$1"
#   direnv allow
#   echo "## Testing $1 ##"
# }


# test_stop {
#   cd $TEST_DIR
#   direnv_eval
# }

### RUN ###

direnv allow || true
direnv_eval

cd $TEST_DIR/scenarios/base
  echo "Testing base"
  direnv_eval
  test "$HELLO" = "world"

  setenv WATCHES $DIRENV_WATCHES
  direnv_eval
  test "$WATCHES" = "$DIRENV_WATCHES"

  sleep 1

  touch .envrc
  direnv_eval
  test "$WATCHES" != "$DIRENV_WATCHES"

  cd ..
  direnv_eval
  echo "$?HELLO"
  test 0 -eq "$?HELLO"
cd $TEST_DIR ; direnv_eval

cd $TEST_DIR/scenarios/inherit
  direnv allow
  echo "Testing inherit"
  direnv_eval
  test "$HELLO" = "world"
cd $TEST_DIR ; direnv_eval

cd $TEST_DIR/scenarios/ruby-layout
  direnv allow
  echo "Testing ruby-layout"
  direnv_eval
  test "$GEM_HOME" != ""
cd $TEST_DIR ; direnv_eval

# Make sure directories with spaces are fine
cd $TEST_DIR/scenarios/"space dir"
  direnv allow
  echo "Testing space dir"
  direnv_eval
  test "$SPACE_DIR" = "true"
cd $TEST_DIR ; direnv_eval

cd $TEST_DIR/scenarios/child-env
  direnv allow
  echo "Testing child-env"
  direnv_eval
  test "$PARENT_PRE" = "1"
  test "$CHILD" = "1"
  test "$PARENT_POST" = "1"
  test 0 -eq "$?REMOVE_ME"
cd $TEST_DIR ; direnv_eval

# cd $TEST_DIR/scenarios/special-vars
#   direnv allow
#   echo "Testing special-vars"
#   setenv DIRENV_BASH `which bash`
#   setenv DIRENV_CONFIG foobar
#   direnv_eval || true
#   test -n "$DIRENV_BASH"
#   test "$DIRENV_CONFIG" = "foobar"
#   unsetenv DIRENV_BASH
#   unsetenv DIRENV_CONFIG
# cd $TEST_DIR ; direnv_eval

cd $TEST_DIR/scenarios/"empty-var"
  direnv allow
  echo "Testing empty-var"
  direnv_eval
  test "$?FOO" -eq 1
  test "$FOO" = ""
cd $TEST_DIR ; direnv_eval

cd $TEST_DIR/scenarios/"empty-var-unset"
  direnv allow
  echo "Testing empty-var-unset"
  setenv FOO ""
  direnv_eval
  test "$?FOO" -eq '0'
  unsetenv FOO
cd $TEST_DIR ; direnv_eval

# Context: foo/bar is a symlink to ../baz. foo/ contains and .envrc file
# BUG: foo/bar is resolved in the .envrc execution context and so can't find
#      the .envrc file.
#
# Apparently, the CHDIR syscall does that so I don't know how to work around
# the issue.
#
# cd $TEST_DIR/scenarios/"symlink-bug"
#   cd foo/bar
#   direnv_eval
# cd $TEST_DIR ; direnv_eval

# Pending: test that the mtime is looked on the original file
# cd $TEST_DIR/scenarios/"utils"
#   LINK_TIME=`direnv file-mtime link-to-somefile`
#   touch somefile
#   NEW_LINK_TIME=`direnv file-mtime link-to-somefile`
#   test "$LINK_TIME" = "$NEW_LINK_TIME"
# cd $TEST_DIR ; direnv_eval
