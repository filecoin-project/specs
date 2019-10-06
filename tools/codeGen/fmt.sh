#!/bin/bash
set -e
set -u

make bin/codeGen

FILES="$(find src -name '*.id')"
for F in $FILES; do
    echo "Formatting in-place: $F"
    cp "$F" build/code/tmp.id
    bin/codeGen fmt build/code/tmp.id "$F"
done
