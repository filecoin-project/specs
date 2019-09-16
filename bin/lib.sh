# no shbang, not meant to be run directly

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
    prun wget -O "$2" "$1"
  fi
}

ensure_symlink() {
  dst="$1"
  src="$2"

  echo "> linking $src -> $dst"
  if [ -L "$src" ]; then
    # link already exists
    dst2=$(readlink "$src")
    if [ "$dst" = "$dst2" ]; then
      # link is what we expect. done
      return
    else
      # link exists, but is different. dont clobber
      echo "link exists, but has different value:"
      echo "\texpected: $src -> $dst"
      echo "\tactual:   $src -> $dst2"
      echo ""
      echo "either fix link manually or rm it and rerun this script"
      return 1
    fi
  elif [ -e "$src" ]; then
    # some other file is at "$src"
    die "$src exists, but is not a link"
  else
    # link doesn't exist, make itt
    prun ln -s "$dst" "$src"
  fi
}

must_run_from_spec_root() {
  # assert we're running from spec root dir
  err="please run $(basename $0) from spec root directory"
  [ -f "$(pwd)/bin/$(basename $0)" ] || die "$err"
  grep 'filecoin-project/specs' "$(pwd)/.git/config" >/dev/null || die "$err"
  grep -i 'filecoin spec' "$(pwd)/README.md" >/dev/null || die "$err"
}