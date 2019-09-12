#!/bin/sh

#------- WIP FLAG
echo "Hold up, this isn't ready yet."
echo "if you need it, ping @jbenet to fix it"
exit 1;
#-------

if [[ $(git status -s) ]]
then
    echo "error: the working directory is dirty."
    echo "Please commit any pending changes and test, before publishing."
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
name=$(bin/build-pdf.sh)
echo $name
cp "pdf-build/$name" "static/$name"
msg="You can also download the full spec in [PDF format](.\/$name)."
sed -i "" "s/<\!\-\- REPLACE_ME_PDF_LINK \-\->/$msg/" INTRO.md

echo "Generating site"
make website

echo "Updating gh-pages branch"
cd public && git add --all && git commit -m "Publishing to gh-pages (publish.sh)"

echo "Publishing to github"
git push origin gh-pages
