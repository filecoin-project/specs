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
  tryinstall rsync rsync

  # Github Actions requires this version of hugo.
  snap install hugo --channel=extended

  # submodules required for hugo themes
  prun git submodule update --init --recursive

  # other packages
  require_version "$(go version)" go 1.13 "recommended install from https://golang.org/dl/"

  return 0
}
main
