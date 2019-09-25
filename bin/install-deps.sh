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
  tryinstall node node

  # other packages
  require_version "$(emacs -version)" emacs 26.3 "recommended install from you package manager"
  require_version "$(go version)" go 1.12 "recommended install from https://golang.org/dl/"
  require_version "$(node --version)" node 10.10 "recommended install from https://nodejs.org/en/"
  require_version "$(npm --version)" npm 5.0 "recommended install from https://nodejs.org/en/"

  # git repos
  prun git submodule update --init --recursive

  # orient
  prun bin/install-deps-orient.sh -y

  # npm deps
  # we don't use package.json in spec yet. we may not.
  cwd=$(pwd)
  cd deps
  npm_install phantomjs-prebuilt
  npm_install mermaid.cli
  cd "$cwd"
}
main
