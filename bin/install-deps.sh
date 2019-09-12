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

# package manager packages
tryinstall emacs emacs
tryinstall hugo hugo
tryinstall dot graphviz

# other packages
require go go "recommended install from https://golang.org/dl/ -- we need version 1.12+"

# git repos
git submodule init
git submodule update

# orient
# TODO: orient (git submodule?)
