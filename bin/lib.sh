# no shbang, not meant to be run directly

find_pkgmgr() {
  pkgmgrs=(brew snap apt-get)
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
$4"
}

require_version() {
  require "$2" "$2" $4

  v_actual=$(echo "$1" | get_version)
  v_expect=$3
  compare_versions "$v_actual" "$v_expect"
  case $? in
    0) return ;;
    1) return ;;
    2) die "$2 version $v_expect or greater required. you have $v_actual
$4" ;;
  esac
}

tryinstall() {
  # no pkg mgr? bail w/ require msg.
  if [ "" = "$pkgmgr" ]; then
    require "$1" "$2"
  else
    which_v "$1" && return 0 # have it
  fi
  version="$3"

  # pkg mgr, try using it
  if [ "$pkgmgr" = "apt-get" ]; then
    prun sudo "$pkgmgr" install "$2"
  else
    prun "$pkgmgr" install "$2"
  fi
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


# from https://stackoverflow.com/questions/4023830/how-to-compare-two-strings-in-dot-separated-version-format-in-bash
compare_versions() {
  if [[ $1 == $2 ]]
  then
    return 0
  fi
  local IFS=.
  local i ver1=($1) ver2=($2)
  # fill empty fields in ver1 with zeros
  for ((i=${#ver1[@]}; i<${#ver2[@]}; i++))
  do
    ver1[i]=0
  done
  for ((i=0; i<${#ver1[@]}; i++))
  do
    if [[ -z ${ver2[i]} ]]
    then
      # fill empty fields in ver2 with zeros
      ver2[i]=0
    fi
    if ((10#${ver1[i]} > 10#${ver2[i]}))
    then
      return 1
    fi
    if ((10#${ver1[i]} < 10#${ver2[i]}))
    then
      return 2
    fi
  done
  return 0
}

get_version() {
  grep -o '[0-9]\+\(\.[0-9]\+\)\+'
}
