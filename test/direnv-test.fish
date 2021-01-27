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

cd (dirname (status -f))
set TEST_DIR $PWD
set XDG_CONFIG_HOME $TEST_DIR/config
set XDG_DATA_HOME $TEST_DIR/data
set -gx PATH (dirname $TEST_DIR) $PATH

# Reset the direnv loading if any
set -x DIRENV_CONFIG $PWD
set -e DIRENV_BASH
set -e DIRENV_DIR
set -e DIRENV_WATCHES
set -e DIRENV_MTIME
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

test_start dump
begin
    set -e LS_COLORS
    direnv_eval
    test_eq "$LS_COLORS" "*.ogg=38;5;45:*.wav=38;5;45"
    test_eq "$LESSOPEN" "||/usr/bin/lesspipe.sh %s"
    test_eq "$THREE_BACKSLASHES" "\\\\\\"
end
test_stop

# Currently broken
# test_start utf-8
# begin
#     direnv_eval
#     test_eq "$UTFSTUFF" "♀♂"
# end
# test_stop
