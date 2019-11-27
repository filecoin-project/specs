#!/usr/bin/env bash

sc=$(pwd)/src/orient
tmp=$(mktemp)
for json in $(ls $sc/*json)
do
    small=$(basename $json)
    echo "[+] jq formatting src/orient/$small "
    jq . $json > $tmp
    cp $tmp $json
done
rm $tmp
