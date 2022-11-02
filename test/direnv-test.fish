#!/usr/bin/env fish
function test_eq --argument-names a b
    if not test (count $argv) = 2
        echo "Error: " (count $argv) " arguments passed to `eq`: $argv"
        exit 1
    end

    if not test $a = $b
        printf "Error:\n - expected: %s\n -      got: %s\n" "$a" "$b"
        exit 1
    end
end

function test_neq --argument-names a b
    if not test (count $argv) = 2
        echo "Error: " (count $argv) " arguments passed to `neq`: $argv"
        exit 1
    end

    if test $a = $b
        printf "Error:\n - expected: %s\n -      got: %s\n" "$a" "$b"
        exit 1
    end
end

function has
    type -q $argv[1]
end

cd (dirname (status -f))
set TEST_DIR $PWD
set XDG_CONFIG_HOME $TEST_DIR/config
set XDG_DATA_HOME $TEST_DIR/data
set -gx PATH (dirname $TEST_DIR) $PATH

# Reset the direnv loading if any
set -x DIRENV_CONFIG $PWD
set -e DIRENV_BASH
set -e DIRENV_DIR
set -e DIRENV_FILE
set -e DIRENV_WATCHES
set -e DIRENV_DIFF

function direnv_eval
    #direnv export fish # for debugging
    direnv export fish | source
end

function test_start -a name
    cd "$TEST_DIR/scenarios/$name"
    direnv allow
    echo "## Testing $name ##"
end

function test_stop
    cd /
    direnv_eval
end

### RUN ###

direnv allow
direnv_eval

test_start base
begin
    echo "Setting up"
    direnv_eval
    test_eq "$HELLO" world

    set WATCHES $DIRENV_WATCHES

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
    echo $HELLO
    test -z "$HELLO" || exit 1

    set -e WATCHES
end
test_stop

test_start inherit
begin
    cp ../base/.envrc ../inherited/.envrc
    direnv_eval
    echo "HELLO should be world:" "$HELLO"

    sleep 1
    echo "export HELLO=goodbye" >../inherited/.envrc
    direnv_eval
    test_eq "$HELLO" goodbye
end
test_stop

if has ruby
    test_start ruby-layout
    begin
        direnv_eval
        test_neq "$GEM_HOME" ""
    end
    test_stop
end

# Make sure directories with spaces are fine
test_start "space dir"
begin
    direnv_eval
    test_eq "$SPACE_DIR" true
end
test_stop

test_start child-env
begin
    direnv_eval
    test_eq "$PARENT_PRE" 1
    test_eq "$CHILD" 1
    test_eq "$PARENT_POST" 1
    test -z "$REMOVE_ME" || exit 1
end
test_stop

test_start special-vars
begin
    set -x DIRENV_BASH (command -s bash)
    set -x DIRENV_CONFIG foobar
    direnv_eval || true
    test -n "$DIRENV_BASH" || exit 1
    test_eq "$DIRENV_CONFIG" foobar
    set -e DIRENV_BASH
    set -e DIRENV_CONFIG
end
test_stop

test_start dump
begin
    set -e LS_COLORS
    direnv_eval
    test_eq "$LS_COLORS" "*.ogg=38;5;45:*.wav=38;5;45"
    test_eq "$LESSOPEN" "||/usr/bin/lesspipe.sh %s"
    test_eq "$THREE_BACKSLASHES" "\\\\\\"
end
test_stop

test_start empty-var
begin
    direnv_eval
    set -q FOO || exit 1
    test_eq "$FOO" ""
end
test_stop

test_start empty-var-unset
begin
    set -x FOO ""
    direnv_eval
    set -q FOO && exit 1
    set -e FOO
end
test_stop

test_start in-envrc
begin
    direnv_eval
    ./test-in-envrc
    test_eq $status 1
end
test_stop

test_start missing-file-source-env
begin
    direnv_eval
end
test_stop

test_start symlink-changed
begin
    # when using a symlink, reload if the symlink changes, or if the
    # target file changes.
    ln -fs ./state-A ./symlink
    direnv_eval
    test_eq "$STATE" A
    sleep 1

    ln -fs ./state-B ./symlink
    direnv_eval
    test_eq "$STATE" B
end
test_stop

# Currently broken
# test_start utf-8
# begin
#     direnv_eval
#     test_eq "$UTFSTUFF" "♀♂"
# end
# test_stop

test_start failure
begin
    # Test that DIRENV_DIFF and DIRENV_WATCHES are set even after a failure.
    #
    # This is needed so that direnv doesn't go into a loop when the loading
    # fails.

    test_eq "$DIRENV_DIFF" ""
    test_eq "$DIRENV_WATCHES" ""

    direnv_eval

    test_neq "$DIRENV_DIFF" ""
    test_neq "$DIRENV_WATCHES" ""

end
test_stop

test_start watch-dir
begin
    echo "No watches by default"
    test_eq "$DIRENV_WATCHES" "$WATCHES"

    direnv_eval

    if ! direnv show_dump $DIRENV_WATCHES | grep -q testfile
        echo "FAILED: testfile not added to DIRENV_WATCHES"
        exit 1
    end

    echo "After eval, watches have changed"
    test_neq "$DIRENV_WATCHES" "$WATCHES"
end
test_stop
