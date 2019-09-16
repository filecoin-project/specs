#!/bin/bash

#------- modes
# set -e

#------- install-deps.sh

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

get_user_confirmation() {
  while : ; do
    read -p "Continue (y/n)? " choice
    case "$choice" in
      y|yes ) break ;;
      n|no ) die "aborting" ;;
    esac
  done
}

download_if_not_present() {
  src="$1"
  dst="$2"
  echo "> installing $dst"
  if [ -f "$dst" ]; then
    echo "$dst already exists. skipping download. to force redownload, remove the file."
  else
    echo "downloading $src to $dst"
    wget -O "$2" "$1"
  fi
}

is_quicklisp_installed() {
  echo "> check whether quicklisp is already installed."
  bin/sbclw --script >deps/quicklisp.installed-or-not << LISP
  (load "bin/sbcl-userinit.lisp") ; --script means --no-userinit
  (if
    (find-package '#:ql)
    (write-line "quicklisp is installed")
    (write-line "quicklisp is not installed"))
LISP
  [ $? -eq 0 ] || die "failed to check"

  # show whether quicklisp is installed to the user reading the trace
  grep 'quicklisp' deps/quicklisp.installed-or-not

  # check if quicklisp is installed
  return $(grep 'quicklisp is installed' deps/quicklisp.installed-or-not >/dev/null)
}

install_quicklisp() {

  # install quicklisp & asdf
  echo "> installing quicklisp and asdf"

  is_quicklisp_installed
  if [ $? -ne 0 ]; then
    download_if_not_present "https://beta.quicklisp.org/quicklisp.lisp" "deps/quicklisp-install.lisp"

    # load quicklisp into sbcl
    # bin/sbclw is a wrapper that sets some options first (for isolation)
    bin/sbclw --script << LISP
    (load "deps/quicklisp-install.lisp")
    (quicklisp-quickstart:install :path "$(pwd)/deps/quicklisp")
LISP
    [ $? -eq 0 ] || die "failed to run quicklisp-install.lisp"

    is_quicklisp_installed || die "failed to install quicklisp"
  fi

  # the below should also work, but harder time getting it to work:
  # bin/sbclw --script << LISP
  # (if
  #   (find-package '#:ql)
  #   (princ "quicklisp already installed")
  #   (progn
  #     (load "deps/quicklisp.lisp")
  #     (eval (quicklisp-quickstart:install :path "$(pwd)/.quicklisp")
  #     (ql:add-to-init-file)))
# LISP

}

install_slime() {
  # add slime to emacs + user startup file
  echo "> adding slime to emacs"
  bin/sbclw --eval "(ql:quickload \"quicklisp-slime-helper\")" --quit
  [ $? -eq 0 ] || die "failed to install slime"
}

link_orient_to_quicklisp() {
  # integrate quicklisp with orient
  echo "> integrating quicklisp with orient"
  ln -s ../../orient/orient.asd deps/quicklisp/local-projects/orient.asd
  # this is a symlink. brittle!
  # if we move orient or quicklisp (relative to each other) we have to change this.
}

install_cllaunch() {
  # installing cl-launch
  cll="cl-launch.sh"
  dst="deps/bin/cl"
  src="https://common-lisp.net/project/xcvb/cl-launch/$cll"

  download_if_not_present "$src" "$dst"
  [ $? -eq 0 ] || die "dailed to install cl-launch"
}

main() {
  # warn the user
  echo "WARNING: $0 is work in progress."
  echo "         This warning will be removed once it is done."

  # get confirmation from user
  if [ "$1" != "-y" ]; then
    get_user_confirmation
  fi

  # package manager packages
  tryinstall wget wget
  tryinstall emacs emacs
  tryinstall sbcl sbcl
  mkdir -p deps/bin
  install_quicklisp
  install_slime
  install_cllaunch
}
main $1

