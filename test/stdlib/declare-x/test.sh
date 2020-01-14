#!/usr/bin/env bash
set -eo pipefail

cd "$(dirname "$0")"

#source ../../../stdlib.sh

please_source() {
  source "$1"
  echo "foo inside: $foo"
  env | grep foo
}

please_source .envrc

#source_env .
#source .envrc

echo foo is $foo
echo bar is $bar
echo baz is $baz

[[ -n $foo ]]
[[ -n $bar ]]

echo OK
