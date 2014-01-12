#!/bin/sh
#
# Small util to display the content of the current DIRENV_DIFF env var.
#

GZIP_HEADER="\x1f\x8b\x08\x00\x00\x00\x00\x00"

(printf $GZIP_HEADER; echo $DIRENV_DIFF | base64 -D) | gzip -dc | python -mjson.tool
