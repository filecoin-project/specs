#!/bin/bash

find_pkgmgr() {
  pkgmgrs=(apt-get brew)
  for pm in ${pkgmgrs[@]}; do
    which "$pm" >/dev/null && printf "$pm" && return 0
  done
  printf ""
}
pkgmgr=$(find_pkgmgr)

die() {
  echo >&2 "error: $@"
  exit 1
}

prun() {
  echo "> $@"
  $@
}

which_v() {
  printf "which $1: "
  which "$1" 2>/dev/null && return 0
  echo "not found" && return 1
}

require() {
  if [ "$1" == "go" ]; then
      gov=$(go version | cut -d" " -f3 | cut -d"." -f2)
      if [[ $gov -lt 12 ]]; then
          die "go v1.12+ required"
      fi
  fi
  which_v "$1" || die "$1 required - install package: $2
$3"
}

tryinstall() {
  # no pkg mgr? bail w/ require msg.
  if [ "" = "$pkgmgr" ]; then
    require "$1" "$2"
  else
    which_v "$1" && return 0 # have it
  fi

  # pkg mgr, try using it
  prun "$pkgmgr" install "$2"
}

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
  # assert we're running from spec root dir
  [ -f "$(pwd)/bin/$(basename $0)" ] || die "please run $(basename $0) from spec root directory"

  # package manager packages
  tryinstall emacs emacs
  tryinstall hugo hugo
  tryinstall dot graphviz

  # other packages
  require go go "recommended install from https://golang.org/dl/ -- we need version 1.12+"

  # git repos
  prun git submodule update --init --recursive

  # orient
  prun bin/install-deps-orient.sh -y
}
main
