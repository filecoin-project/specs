#!/bin/bash

source "$(dirname $0)/lib.sh"

main() {
  must_run_from_spec_root

  # package manager packages
  tryinstall dot graphviz
  tryinstall node node

  # other packages
  require_version "$(node --version)" node 10.10 "recommended install from https://nodejs.org/en/"
  require_version "$(npm --version)" npm 5.0 "recommended install from https://nodejs.org/en/"

  # npm deps
  echo "> installing npm deps"
  cwd=$(pwd)
  cd deps
  npm install || die "npm install failed"
  cd "$cwd"
  return 0
}
main
