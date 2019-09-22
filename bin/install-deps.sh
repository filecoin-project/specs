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
  tryinstall emacs emacs
  tryinstall hugo hugo
  tryinstall dot graphviz
  tryinstall rsync rsync

  # other packages
  require go go "recommended install from https://golang.org/dl/ -- we need version 1.12+"

  # git repos
  prun git submodule update --init --recursive

  # orient
  prun bin/install-deps-orient.sh -y
}
main
