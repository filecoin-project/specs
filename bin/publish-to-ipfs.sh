#!/bin/bash

die() {
  echo >&2 "error: $@"
  exit 1
}

dir=build/website

ipfs swarm peers >/dev/null 2>/dev/null || \
  die "is ipfs daemon not running?"

[ -d "$dir" ] || die "$dir not found. did you run: make website ?"

hash=$(ipfs add -Q -r "$dir")
echo "published /ipfs/$hash"
echo "http://localhost:8080/ipfs/$hash"
echo "https://ipfs.io/ipfs/$hash"
