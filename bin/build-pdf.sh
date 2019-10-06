#!/bin/sh

#------- WIP FLAG
echo "Hold up, this isn't ready yet."
echo "if you need it, ping @jbenet to fix it"
exit 1;
#-------

set -e

short=$(git rev-parse --short HEAD)
tag=${1-$short}

files=(
    INTRO.md
    data-structures.md
    address.md
    signatures.md
    proofs.md
    validation.md
    network-protocols.md
    bootstrap.md
    data-propagation.md
    sync.md
    expected-consensus.md
    state-machine.md
    local-storage.md
    operation.md
    actors.md
    mining.md
    storage-market.md
    retrieval-market.md
    payments.md
    faults.md
    zigzag-circuit.md
    zigzag-porep.md
    definitions.md
    style.md
    process.md
)

rm -rf .pdfworking
mkdir -p .pdfworking
mkdir -p pdf-build

i=0
for f in ${files[@]}; do
    printf -v n "%03d" $i
    cp "./${f}" "./.pdfworking/${n}_${f}"
    i=$((i + 1))
done

find ./.pdfworking -type f -exec sed -i "" 's/{{%.*%}}//g' {} \;

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
