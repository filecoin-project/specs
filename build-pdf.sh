#!/bin/sh
set -e

die() {
    echo >&2 "error: $@"
    exit 1
}

require() {
    which "$1" >/dev/null || die "please install $1"
}

require_tex() {
    tlmgr info --only-installed "$1" | grep installed | grep Yes >/dev/null || \
        die "please install $1:\n  sudo tlmgr install $1"
}

require pandoc
require xelatex

require_tex adjustbox
require_tex collectbox
require_tex enumitem
require_tex pagecolor
require_tex csquotes
require_tex mdframed
require_tex needspace
require_tex titling
require_tex lm-math

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
    -V fontsize=9pt \
    --metadata title="Filecoin Specification $tag" \
    --metadata subtitle="$(date -u '+%Y-%m-%d %H:%M:%S Z')"

rm -rf .pdfbuild

echo "open pdf-build/filecoin-spec-$tag.pdf"
