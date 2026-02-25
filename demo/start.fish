# For debugging
# set fish_trace on
# set -gx DIRENV_DEBUG true
# set -gx COMPLETE_FISH_DEBUG true

make
# Use `direnv exec` to clear an existing direnv environment if there is one.
direnv="$PWD/direnv" ./direnv exec / fish --no-config --init-command '
  eval "$("$direnv" hook fish)"
  "$direnv" allow
'