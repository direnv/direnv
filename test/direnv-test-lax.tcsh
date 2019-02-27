#!/usr/bin/env tcsh

# NB. This script is ran without the -e option. That means that it will not
# stop automatically on errors. We need this when testing tcsh-specific code
# that need to perform tests at runtime (e.g. using grep), since any non-0 exit
# code will cause the script to fail (can't be escaped, disabled, etc.).
#
# If possible, place new tests in direnv-test.tcsh

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
unsetenv DIRENV_ON_UNLOAD_tcsh

alias direnv_eval 'eval `direnv export tcsh` || exit 1'

cd $TEST_DIR/scenarios/"alias"
  direnv allow
  echo "Testing alias"

  unalias foo
  alias bar 'original bar'

  direnv_eval

  test "`alias foo`" = "ls -l" || exit 1
  test "`alias bar`" = "ls -t" || exit 1

  cd nested
  direnv allow
  direnv_eval

  test -z "`alias foo`" || exit 1
  test "`alias bar`" = "df -h" || exit 1

  cd ..
  direnv_eval

  test "`alias foo`" = "ls -l" || exit 1
  test "`alias bar`" = "ls -t" || exit 1

  cd ..
  direnv_eval

  test -z "`alias foo`" || exit 1
  test "`alias bar`" = "original bar" || exit 1
cd $TEST_DIR ; direnv_eval
