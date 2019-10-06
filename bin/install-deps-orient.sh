#!/bin/bash

#------- modes
# set -e

# load lib
source "$(dirname $0)/lib.sh"

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
}

install_slime() {
  echo "> installing slime"
  bin/sbclw --eval "(ql:quickload \"quicklisp-slime-helper\")" --quit
  [ $? -eq 0 ] || die "failed to install slime"
}

link_orient_to_quicklisp() {
  echo "> integrating quicklisp with orient"
  dst="../../orient/orient.asd"
  src="deps/quicklisp/local-projects/orient.asd"
  ensure_symlink "$dst" "$src"

  # this is a symlink. brittle!
  # if we move orient or quicklisp (relative to each other) we have to change this.
}

install_cllaunch() {
  cll="cl-launch.sh"
  dst="deps/bin/cl"
  src="https://common-lisp.net/project/xcvb/cl-launch/$cll"

  download_if_not_present "$src" "$dst"
  [ $? -eq 0 ] || die "failed to install cl-launch"
}

install_emacs_deps() {
  # loading the init file will do
  emacs_init="bin/emacs-init-build.el"
  echo "> installing emacs deps from $emacs_init"
  HOME=$HOME emacs -Q --script "$emacs_init"
  [ $? -eq 0 ] || die "failed to install emacs deps"
}

main() {
  must_run_from_spec_root

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
  require_version "$(emacs -version)" emacs 26.3 "recommended install from you package manager"
  tryinstall sbcl sbcl
  mkdir -p deps/bin

  # lisp / emacs
  install_quicklisp
  install_slime
  install_cllaunch
  install_emacs_deps

  # git repos
  prun git submodule update --init --recursive

  # orient
  link_orient_to_quicklisp
  return 0
}
main $1
