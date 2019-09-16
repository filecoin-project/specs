#!/bin/sh

source "$(dirname $0)/lib.sh"

src=build/website
bdir=build/gh-pages

#------- WIP FLAG
echo "Hold up, this isn't ready yet."
echo "if you need it, ping @jbenet to fix it\n"

echo "WARNING: this will publish to github pages."
# get confirmation from user
if [ "$1" != "-y" ]; then
  get_user_confirmation
fi

#-------

[[ $(git status -s -uno) ]] && die "the working directory is dirty.\n\
Please commit any pending changes and test, before publishing."

[ -d "$src" ] || die "$src not found. did you run: make website ?"

echo "Prepare worktree"
# rm -rf public
# mkdir public
# git worktree prune
# rm -rf .git/worktrees/public/

echo "Checking out gh-pages branch into build/gh-pages"
git worktree add -B gh-pages "$bdir" origin/gh-pages

echo "Removing existing files"
rm -rf "$bdir/*"

# echo "Generating PDF"
# name=$(bin/build-pdf.sh)
# echo $name
# cp "pdf-build/$name" "static/$name"
# msg="You can also download the full spec in [PDF format](.\/$name)."
# sed -i "" "s/<\!\-\- REPLACE_ME_PDF_LINK \-\->/$msg/" INTRO.md

echo "Updating gh-pages branch"
shash=$(git rev-parse --short HEAD)
cd "$bdir" && git add --all && git commit -m "Publishing $shash ($0)"

echo "Publishing to github"
die "cannot publish to github. this is WIP, so only dry runs."
git push origin gh-pages
