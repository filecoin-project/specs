#!/bin/bash

source "$(dirname $0)/lib.sh"

# usage() {
#   cat << USAGE
# SYNOPSIS
#     install dependencies for the filecoin spec buildsys
#     usage: $0 [-y]
#
# OPTIONS
#     -h,--help      show usage
#     -y,--yes       pass yes to installers and confirmation prompts
# USAGE
# }
#
# parse_args() {
#   while [ $# -gt 0 ]; do
#     case "$1" in
#     -y|--yes) y=y ;;
#     -h|--help) usage ; exit 0 ;;
#     *) die "unrecognized argument: $1 (-h shows usage)" ;;
#     esac
#     shift
#   done
# }

main() {
  must_run_from_spec_root

  # package manager packages
  tryinstall hugo hugo
  tryinstall rsync rsync

  # submodules required for hugo themes
  prun git submodule update --init --recursive

  # other packages
  require_version "$(hugo version)" hugo 0.54 "recommended install from package manager"
  require_version "$(go version)" go 1.12 "recommended install from https://golang.org/dl/"
  return 0
}
main
