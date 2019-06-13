#!/bin/sh
set -e

short=$(git rev-parse --short HEAD)
tag=${1-$short}

files=$(awk -F "[(/)]" '{ if($2) print $2; }' docs/menu/index.md)

rm -rf .pdfworking
mkdir -p .pdfworking
mkdir -p pdf-build

cp "./docs/_index.md" "./.pdfworking/0000_start.md"

i=1
for f in ${files[@]}; do
    printf -v n "%03d" $i
    cp "./docs/${f}" "./.pdfworking/${n}_${f}"
    i=$((i + 1))
done

pandoc ./.pdfworking/*.md \
    --pdf-engine xelatex \
    -o "pdf-build/filecoin-spec-$tag".pdf \
    --from markdown+grid_tables \
    --template ./pdf/eisvogel.tex \
    -H ./pdf/helpers.tex \
    --listings \
    --toc \
    -V titlepage=true \
    -V titlepage-color=FFFFFF \
    -V titlepage-rule-color=FFFFFF \
    -V titlepage-rule-height=0 \
    -V logo="./pdf/cover-src.jpg" \
    -V logo-width="120" \
    -V listings-disable-line-numbers=true \
    -V block-headings=true \
    -V mainfont="Georgia" \
    --metadata title="Filecoin Specification $tag"

rm -rf .pdfbuild

echo "filecoin-spec-$tag.pdf"
