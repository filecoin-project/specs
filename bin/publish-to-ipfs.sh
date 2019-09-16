#!/bin/bash

source "$(dirname $0)/lib.sh"
must_run_from_spec_root

dir=build/website

ipfs swarm peers >/dev/null 2>/dev/null || \
  die "is ipfs daemon not running?"

[ -d "$dir" ] || die "$dir not found. did you run: make website ?"

hash=$(ipfs add -Q -r "$dir")
echo "published /ipfs/$hash"
echo "http://localhost:8080/ipfs/$hash"
echo "https://ipfs.io/ipfs/$hash"
