#!/usr/bin/env bash
set -eo pipefail

cd "$(dirname "$0")"

#source ../../../stdlib.sh

please_source() {
  source "$1"
}

please_source .envrc

#source_env .
#source .envrc

echo foo is $foo
echo bar is $bar

[[ -n $foo ]]
[[ -n $bar ]]

echo OK
