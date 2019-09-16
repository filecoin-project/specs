#!/bin/bash

# common code
source "$(dirname $0)/lib.sh"

sbclr() {
  sbcl --noinform --non-interactive "$@" --quit
}

is_quicklisp_installed() {
  echo "> check whether quicklisp is already installed."

  lisp_qlcheck="\
    (if\
      (find-package '#:ql)\
      (write-line \"quicklisp is installed\")\
      (write-line \"quicklisp is not installed\"))\
  "

  qli=$(sbclr --eval "$lisp_qlcheck")
  [ $? -eq 0 ] || die "failed to check"

  # show whether quicklisp is installed to the user reading the trace
  echo "$qli"

  # check if quicklisp is installed
  return $(echo "$qli" | grep 'quicklisp is installed' >/dev/null)
}

install_quicklisp() {
  echo "> installing quicklisp and asdf"

  is_quicklisp_installed
  if [ $? -ne 0 ]; then
    qlif="$(pwd)/deps/quicklisp-install.lisp"
    download_if_not_present "https://beta.quicklisp.org/quicklisp.lisp" "$qlif"

    sbclr --load "$qlif" \
      --eval "(quicklisp-quickstart:install :path \"$HOME/quicklisp\")" \
      --eval "(ql-util:without-prompting (ql:add-to-init-file))"
    [ $? -eq 0 ] || die "failed to run quicklisp-install.lisp"

    is_quicklisp_installed || die "failed to install quicklisp"
  fi
}

install_slime() {
  echo "> installing slime"
  sbclr --eval "(ql:quickload \"quicklisp-slime-helper\")"
  [ $? -eq 0 ] || die "failed to install slime"
}

link_orient_to_quicklisp() {
  echo "> integrating quicklisp with orient"

  dst="$(pwd)/orient/orient.asd"
  src="$HOME/quicklisp/local-projects/orient.asd"
  ensure_symlink "$dst" "$src"

  # this is a symlink. brittle!
  # if we move orient or quicklisp (relative to each other) we have to change this.
}

install_cllaunch() {
  cll="cl-launch.sh"
  dst="/usr/local/bin/cl"
  src="https://common-lisp.net/project/xcvb/cl-launch/$cll"

  download_if_not_present "$src" "$dst"
  [ $? -eq 0 ] || die "failed to install cl-launch"
}

install_emacs_deps() {
  # loading the init file will do
  emacs_init="bin/emacs-init-user.el"
  echo "> installing emacs deps from $emacs_init"
  emacs --script "$emacs_init"
  [ $? -eq 0 ] || die "failed to install emacs deps"
}

main() {
  # warn the user
  echo "WARNING: $0 is work in progress."
  echo "         This warning will be removed once it is done."

  echo "WARNING: this may mess up your user environment, especially if you use:"
  echo "         emacs"
  echo "         sbcl"

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
  install_emacs_deps
  link_orient_to_quicklisp

  # TODO: make this automated.
  emacs_init="$(pwd)/bin/emacs-init-user.el"
  cat <<- EOF

  =============== ORIENT DEPS INSTALL FINISHED ===============

  We need to link emacs, slime, orgmode, and so on.
  Do this by loading this file: $emacs_init

  To load it, use this elisp code in an emacs session, or put it in your emacs init file:
    (load "$emacs_init")

  Your emacs user init file is usually ~/.emacs, or ~/.emacs.el, or ~/emacs.d/init.el
  (TODO automate this)
EOF

}
main $1
