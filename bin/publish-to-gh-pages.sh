#!/bin/sh

source "$(dirname $0)/lib.sh"

website=build/website
ghpages=build/gh-pages

## -- warnings
# echo "WARNING: this will publish to github pages."
# get confirmation from user
# if [ "$1" != "-y" ]; then
#  get_user_confirmation
# fi

## -- checks
must_run_from_spec_root
# [[ $(git status -s -uno) ]] && die "the working directory is dirty.\n\
# Please commit any pending changes and test, before publishing."
[ -d "$website" ] || die "$website not found. did you run: make website ?"


## -- work tree setup
echo "Setting up worktree in build/gh-pages"
git fetch origin gh-pages
git worktree remove gh-pages
rm -rf "$ghpages"
git worktree add -B gh-pages "$ghpages" origin/gh-pages

## -- updating website content
echo "Removing existing files"
dotgit=$(cat "$ghpages/.git") # preserve .git
rm -rf "$ghpages"
cp -r "$website" "$ghpages"
echo "$dotgit" >"$ghpages/.git" # paste .git

# echo "Generating PDF"
# name=$(bin/build-pdf.sh)
# echo $name
# cp "pdf-build/$name" "static/$name"
# msg="You can also download the full spec in [PDF format](.\/$name)."
# sed -i "" "s/<\!\-\- REPLACE_ME_PDF_LINK \-\->/$msg/" INTRO.md

echo "Updating gh-pages branch"
shash=$(git rev-parse --short HEAD)
cd "$ghpages" && git add --all && git commit -m "Publishing $shash ($0)"

## -- publishing
echo "Publishing to github"
die "cannot publish to github. this is WIP, so only dry runs."
# git push origin gh-pages
