#!/bin/sh

if [[ $(git status -s) ]]
then
    echo "The working directory is dirty. Please commit any pending changes."
    exit 1;
fi

echo "Deleting old publication"
rm -rf public
mkdir public
git worktree prune
rm -rf .git/worktrees/public/

echo "Checking out gh-pages branch into public"
git worktree add -B gh-pages public origin/gh-pages

echo "Removing existing files"
rm -rf public/*

echo "Generating PDF"
name=$(./build-pdf.sh)
echo $name
cp "pdf-build/$name" "static/$name"
msg="You can also download this spec in [PDF format](.\/$name)."
sed -i "" "s/<\!\-\- REPLACE_ME_PDF_LINK \-\->/$msg/" INTRO.md

echo "Generating site"
hugo

echo "Updating gh-pages branch"
cd public && git add --all && git commit -m "Publishing to gh-pages (publish.sh)"

echo "Publishing to github"
git push origin gh-pages
