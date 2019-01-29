#!/usr/bin/env bash

cd "$(dirname "$0")"

source ../../../stdlib.sh

source_env .
#source .envrc

echo foo is $foo
echo bar is $bar

[[ -n $bar ]]
[[ -n $foo ]]

echo OK
