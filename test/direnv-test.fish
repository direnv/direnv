#!/usr/bin/env fish
function eq --argument-names a b
	if not test (count $argv) = 2
		echo "Error: " (count $argv) " arguments passed to `eq`: $argv"
		exit 1
	end

	if not test $a = $b
		printf "Error:\n - expected: %s\n -      got: %s\n" "$a" "$b"
		exit 1
	end
end

function ge --argument-names a b
   if not test (count $argv) = 2
    echo "Error: " (count $argv) " arguments passed to `ge`: $argv"
    exit 1
   end

   if not test $a -ge $b
    printf "Error: expected %s > %s\n" "$a" "$b"
    exit 1
   end
end

function lt --argument-names a b
   if not test (count $argv) = 2
    echo "Error: " (count $argv) " arguments passed to `lt`: $argv"
    exit 1
   end

   if not test $a -lt $b
    printf "Error: expected %s < %s\n" "$a" "$b"
    exit 1
   end
end

cd (dirname (status -f))
set TEST_DIR $PWD
set -gx PATH (dirname $TEST_DIR) $PATH

# Reset the direnv loading if any
set -x DIRENV_CONFIG $PWD
set -e DIRENV_BASH
set -e DIRENV_DIR
set -e DIRENV_WATCHES
set -e DIRENV_MTIME
set -e DIRENV_DIFF
set -e DIRENV_ON_UNLOAD_fish

function direnv_eval
  #direnv export fish # for debugging
  direnv export fish | source
end

function test_start -a name
  cd "$TEST_DIR/scenarios/$name"
  direnv allow
  echo "## Testing $name ##"
  pwd
end

function test_stop
  cd $TEST_DIR
  direnv_eval
end

### RUN ###

direnv allow
direnv_eval

test_start dump
	set -e LS_COLORS
	direnv_eval
	eq "$LS_COLORS" "*.ogg=38;5;45:*.wav=38;5;45"
	eq "$LESSOPEN" "||/usr/bin/lesspipe.sh %s"
	eq "$THREE_BACKSLASHES" "\\\\\\"
test_stop

test_start "shell-specific"
  set -e BAR
  set -e FOO
  set -e FOOX
  set -e FOO_OR_NAN
  set -e FOOX_OR_NAN
  set -e BAR_OR_NAN

  set -x TARGET_SHELL fish

  direnv_eval
  # FOO=x; ON UNLOAD WILL RUN FOO=x + 100
  # export FOOX=y; ON UNLOAD WILL RUN FOOX=FOOX + 100

  ge "$FOO" 0
  ge "$FOOX" 0
  ge "$BAR" 0

  set FOO_0 "$FOO"
  set FOOX_0 "$FOOX"
  set BAR_0 "$BAR"

  direnv allow nested/.envrc
  cd nested
  direnv_eval
  # Set by nested/.envrc
  eq "NaN" "$FOO_OR_NAN"
  eq "NaN" "$BAR_OR_NAN"
  eq "$FOOX_0" "$FOOX_OR_NAN"  # unlike BAR

  # Set by unload action of ./.envrc
  eq (math "$FOO_0" + 100) "$FOO"

  # Set by combination of action in nested/..envrc and unload of .envrc
  eq (math \($FOOX_0 + 100\) '%' 3) "$FOOX"

  cd ..
  direnv_eval

  # New random values
  ge "$FOO" 0
  ge "$BAR" 0

  # Random value overrides the FOOX + 1000 on_unload from nested/.envrc
  lt "$FOOX" 100

  set FOO_1 "$FOO"
  set BAR_1 "$BAR"
  set FOOX_1 "$FOOX"

  cd ..
  direnv_eval

  # Unload actions from .envrc
  eq (math "$FOO_1" + 100) "$FOO"
  eq (math "$FOOX_1" + 100) "$FOOX"
test_stop

test_start "alias"
  functions -e foo
  alias bar='original bar'

  function expected_alias_def --argument-names name val
    printf "function %s --description 'alias %s=%s'\n\t%s \$argv;\nend" "$name" "$name" "$val" "$val"
  end

  direnv_eval

  set expected_foo (expected_alias_def foo "ls -l")
  set actual_foo (functions foo | grep -v '^#')
  eq "$expected_foo" "$actual_foo"

  set expected_bar (expected_alias_def bar "ls -t")
  set actual_bar (functions bar | grep -v '^#')
  eq "$expected_bar" "$actual_bar"

  cd nested
  direnv allow
  direnv_eval

  set expected_foo ""
  set actual_foo (functions foo; or true)
  eq "$expected_foo"  "$actual_foo"

  set expected_bar (expected_alias_def bar "df -h")
  set actual_bar (functions bar | grep -v '^#')
  eq "$expected_bar" "$actual_bar"

  cd ..
  direnv_eval

  set expected_foo (expected_alias_def foo "ls -l")
  set actual_foo (functions foo | grep -v '^#')
  eq "$expected_foo" "$actual_foo"

  set expected_bar (expected_alias_def bar "ls -t")
  set actual_bar (functions bar | grep -v '^#')
  eq "$expected_bar" "$actual_bar"

  cd ..
  direnv_eval

  set expected_foo ""
  set actual_foo (functions foo; or true)
  eq "$expected_foo" "$actual_foo"

  set expected_bar (expected_alias_def bar "original bar")
  set actual_bar (functions bar | grep -v '^#')
  eq "$expected_bar" "$actual_bar"

test_stop
