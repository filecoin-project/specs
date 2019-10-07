#!/bin/bash

source "$(dirname $0)/lib.sh"

main() {
  must_run_from_spec_root

  hash=$(git rev-parse HEAD)
  shorthash=$(git rev-parse --short HEAD)
  date=$(date -u '+%Y-%m-%d_%H:%M:%SZ')

  while read line; do
      eval echo "$line"
  done
}
main
