# Don't run these tests if we don't have cygpath.

if ! has cygpath; then
  return 0
fi

# These tests apply to systems with cygpath (MSYS2/Cygwin). Some of these tests
# may be moved to the common test where it makes sense.

# This shows that MSYS configuration variables don't seep into the exported
# environment.
test_start 'cygpath-conf'
  direnv_eval
  test_eq "${HELLO}" "world"
  test_eq "${MSYS_NO_PATHCONV-unset}" "unset"
  test_eq "${MSYS_NO_PATHCONV:-empty}" "empty"
  test_eq "${MSYS_NO_PATHCONV}" ""
  test -z "$MSYS_NO_PATHCONV"
  test_eq "${MSYS2_ENV_CONV_EXCL-unset}" "unset"
  test_eq "${MSYS2_ENV_CONV_EXCL:-empty}" "empty"
  test_eq "${MSYS2_ENV_CONV_EXCL}" ""
  test -z "$MSYS2_ENV_CONV_EXCL"
  test_eq "${MSYS2_ARG_CONV_EXCL-unset}" "unset"
  test_eq "${MSYS2_ARG_CONV_EXCL:-empty}" "empty"
  test_eq "${MSYS2_ARG_CONV_EXCL}" ""
  test -z "$MSYS2_ARG_CONV_EXCL"
test_stop

# This shows that MSYS configuration variables don't seep into the exported
# environment.
test_start 'cygpath-conf'
  direnv_eval_msys2
  test_eq "${HELLO}" "world"
  test_eq "${MSYS_NO_PATHCONV-unset}" "unset"
  test_eq "${MSYS_NO_PATHCONV:-empty}" "empty"
  test_eq "${MSYS_NO_PATHCONV}" ""
  test -z "$MSYS_NO_PATHCONV"
  test_eq "${MSYS2_ENV_CONV_EXCL-unset}" "unset"
  test_eq "${MSYS2_ENV_CONV_EXCL:-empty}" "empty"
  test_eq "${MSYS2_ENV_CONV_EXCL}" ""
  test -z "$MSYS2_ENV_CONV_EXCL"
  test_eq "${MSYS2_ARG_CONV_EXCL-unset}" "unset"
  test_eq "${MSYS2_ARG_CONV_EXCL:-empty}" "empty"
  test_eq "${MSYS2_ARG_CONV_EXCL}" ""
  test -z "$MSYS2_ARG_CONV_EXCL"
test_stop

# This shows what happens to path and path list variables (part 1).
test_start 'cygpath-perf'
  direnv_eval

  test_eq "${HELLO}" "world"

  [[ "$PATH" == *";"* ]] && { echo "FAILED: PATH contains semicolons."; exit 1; }
  [[ "$PATH" != *"/bin"* ]] && { echo "FAILED: PATH does not contain /bin."; exit 1; }
  [[ "$PATH" != *"/foo"* ]] && { echo "FAILED: PATH does not contain /foo."; exit 1; }

  test_eq "${FILE1}" "C:\\Windows\\System32\\cmd.exe"
  test_eq "${FILE2}" "C:/Windows/System32/cmd.exe"

  if is_msys2; then
    # MSYS2 changed our Unix path to a Windows path
    test_neq "${FILE3}" "/c/Windows/System32/cmd.exe"
  fi

  if is_cygwin; then
    # Cygwin kept our Unix path
    test_eq "${FILE3}" "/c/Windows/System32/cmd.exe"
  fi

  test_eq "${LIST1}" "C:\\Windows\\System32;C:\\Windows\\System32"
  test_eq "${LIST2}" "C:/Windows/System32;C:/Windows/System32"

  if is_msys2; then
    # MSYS2 changed our Unix path list to a Windows path list
    test_neq "${LIST3}" "/c/Windows/System32:/c/Windows/System32"
  fi

  if is_cygwin; then
    test_eq "${LIST3}" "/c/Windows/System32:/c/Windows/System32"
  fi

test_stop

# This shows what happens to path and path list variables (part 2).
test_start 'cygpath-perf'
  direnv_eval_msys2

  test_eq "${HELLO}" "world"

  [[ "$PATH" == *";"* ]] && { echo "FAILED: PATH contains semicolons."; exit 1; }
  [[ "$PATH" != *"/bin"* ]] && { echo "FAILED: PATH does not contain /bin."; exit 1; }
  [[ "$PATH" != *"/foo"* ]] && { echo "FAILED: PATH does not contain /foo."; exit 1; }

  test_eq "${FILE1}" "C:\\Windows\\System32\\cmd.exe"
  test_eq "${FILE2}" "C:/Windows/System32/cmd.exe"

  # MSYS2 kept our Unix path
  test_eq "${FILE3}" "/c/Windows/System32/cmd.exe"

  test_eq "${LIST1}" "C:\\Windows\\System32;C:\\Windows\\System32"
  test_eq "${LIST2}" "C:/Windows/System32;C:/Windows/System32"

  # MSYS2 kept our Unix path
  test_eq "${LIST3}" "/c/Windows/System32:/c/Windows/System32"

test_stop

# This shows what happens to path and path list variables (part 3).
# After a direnv_eval_msys2, we just make sure that direnv_eval works as
# expected.
test_start 'cygpath-perf'
  direnv_eval

  if is_msys2; then
    test_neq "${FILE3}" "/c/Windows/System32/cmd.exe"
    test_neq "${LIST3}" "/c/Windows/System32:/c/Windows/System32"
  fi

  if is_cygwin; then
    test_eq "${FILE3}" "/c/Windows/System32/cmd.exe"
    test_eq "${LIST3}" "/c/Windows/System32:/c/Windows/System32"
  fi

test_stop

# This highlights differences in cygpath environments and the lookPath function.
# Here we check that exec path resolution works. Part 1.
# NOTE: Windows does not have exec like Unix, so in this example the called
# programs will not replace the direnv process.
test_start 'cygpath-perf'
  direnv_exec . hostname
  direnv_exec . "$(which hostname)"
  direnv_exec . env hostname
  direnv_exec . "$(which env)" hostname
  direnv_exec . "$(which env)" "$(which hostname)"
  direnv_exec . env env hostname
  direnv_exec . which direnv
  direnv_exec . direnv version
  direnv_exec . direnv exec . which direnv
  direnv_exec . direnv exec . direnv version
  grep -Ff - <(direnv_exec . env) <<EOF
FILE1=C:\\Windows\\System32\\cmd.exe
FILE2=C:/Windows/System32/cmd.exe
FILE3=C:/Windows/System32/cmd.exe
EOF
  grep -Ff - <(direnv_exec . direnv exec . env) <<EOF
FILE1=C:\\Windows\\System32\\cmd.exe
FILE2=C:/Windows/System32/cmd.exe
FILE3=C:/Windows/System32/cmd.exe
EOF
test_stop

# This highlights differences in cygpath environments and the lookPath function.
# Here we check that exec path resolution works. Part 2.
test_start 'cygpath-perf'
  direnv_exec_msys2 . hostname
  direnv_exec_msys2 . "$(which hostname)"
  direnv_exec_msys2 . env hostname
  direnv_exec_msys2 . "$(which env)" hostname
  direnv_exec_msys2 . "$(which env)" "$(which hostname)"
  direnv_exec_msys2 . env env hostname
  direnv_exec_msys2 . which direnv
  direnv_exec_msys2 . direnv version
  direnv_exec_msys2 . direnv exec . which direnv
  direnv_exec_msys2 . direnv exec . direnv version
  grep -Ff - <(direnv_exec_msys2 . env) <<EOF
FILE1=C:\\Windows\\System32\\cmd.exe
FILE2=C:/Windows/System32/cmd.exe
FILE3=/c/Windows/System32/cmd.exe
EOF
  grep -Ff - <(direnv_exec_msys2 . direnv exec . env) <<EOF
FILE1=C:\\Windows\\System32\\cmd.exe
FILE2=C:/Windows/System32/cmd.exe
FILE3=/c/Windows/System32/cmd.exe
EOF
test_stop

# This shows that the direnv allow and direnv reload commands work as expected
# in cygpath environments.
# This test is similar in some respects to the inherit, symlink-dir, load-env
# and aliases tests.
# Part 1.
test_start 'cygpath-path'

  direnv_eval && [[ "$PATH" == *";"* ]] && { echo "FAILED: PATH contains semicolons."; exit 1; }

  sleep 1
  symlink ./state-A ./symlink
  direnv allow && direnv_eval && [[ "$PATH" == *";"* ]] && { echo "FAILED: PATH contains semicolons."; exit 1; }
  test_eq "${STATE}" "A"

  sleep 1
  symlink ./state-B ./symlink
  direnv allow && direnv_eval && [[ "$PATH" == *";"* ]] && { echo "FAILED: PATH contains semicolons."; exit 1; }
  test_eq "${STATE}" "B"

  sleep 1
  symlink ./empty ./symlink
  direnv allow && direnv_eval && [[ "$PATH" == *";"* ]] && { echo "FAILED: PATH contains semicolons."; exit 1; }
  test_eq "${STATE-}" ""
  test -z "$STATE"

  sleep 1
  symlink ./state-A ./symlink
  direnv reload && direnv_eval && [[ "$PATH" == *";"* ]] && { echo "FAILED: PATH contains semicolons."; exit 1; }
  test_eq "${STATE}" "A"

  sleep 1
  symlink ./state-B ./symlink
  direnv reload && direnv_eval && [[ "$PATH" == *";"* ]] && { echo "FAILED: PATH contains semicolons."; exit 1; }
  test_eq "${STATE}" "B"

  sleep 1
  symlink ./empty ./symlink
  direnv reload && direnv_eval && [[ "$PATH" == *";"* ]] && { echo "FAILED: PATH contains semicolons."; exit 1; }
  test_eq "${STATE-}" ""
  test -z "$STATE"

test_stop

# This shows that the direnv allow and direnv reload commands work as expected
# in cygpath environments.
# This test is similar in some respects to the inherit, symlink-dir, load-env
# and aliases tests.
# Part 2.
test_start 'cygpath-path'

  direnv_eval_msys2 && [[ "$PATH" == *";"* ]] && { echo "FAILED: PATH contains semicolons."; exit 1; }

  sleep 1
  symlink ./state-A ./symlink
  direnv allow && direnv_eval_msys2 && [[ "$PATH" == *";"* ]] && { echo "FAILED: PATH contains semicolons."; exit 1; }
  test_eq "${STATE}" "A"

  sleep 1
  symlink ./state-B ./symlink
  direnv allow && direnv_eval_msys2 && [[ "$PATH" == *";"* ]] && { echo "FAILED: PATH contains semicolons."; exit 1; }
  test_eq "${STATE}" "B"

  sleep 1
  symlink ./empty ./symlink
  direnv allow && direnv_eval_msys2 && [[ "$PATH" == *";"* ]] && { echo "FAILED: PATH contains semicolons."; exit 1; }
  test_eq "${STATE-}" ""
  test -z "$STATE"

  sleep 1
  symlink ./state-A ./symlink
  direnv reload && direnv_eval_msys2 && [[ "$PATH" == *";"* ]] && { echo "FAILED: PATH contains semicolons."; exit 1; }
  test_eq "${STATE}" "A"

  sleep 1
  symlink ./state-B ./symlink
  direnv reload && direnv_eval_msys2 && [[ "$PATH" == *";"* ]] && { echo "FAILED: PATH contains semicolons."; exit 1; }
  test_eq "${STATE}" "B"

  sleep 1
  symlink ./empty ./symlink
  direnv reload && direnv_eval_msys2 && [[ "$PATH" == *";"* ]] && { echo "FAILED: PATH contains semicolons."; exit 1; }
  test_eq "${STATE-}" ""
  test -z "$STATE"

test_stop
