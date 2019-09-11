#!/bin/bash

die() {
  echo >&2 "error: $@"
  exit 1
}

ipfs swarm peers >/dev/null 2>/dev/null || \
  die "is ipfs daemon not running?"
hash=$(ipfs add -Q -r public)
echo "published /ipfs/$hash"
echo "http://localhost:8080/ipfs/$hash"
echo "https://ipfs.io/ipfs/$hash"
