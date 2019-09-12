
build: website

# main Targets

# this bundles the website
website: hugo-build

publish: website
	bin/publish-to-ipfs

# intermediate targets

hugo-build: $(shell find . | grep .md)
	hugo

# todo
generate-code: $(shell find content/ | grep .ipld)
	echo TODO: use codeGen && exit 1
	# bin/codeGen <input-files?>

go-test: $(shell find content/ | grep .go)
	# testing should have the side effect that all go is compiled
	cd content/codeGen && go build && go test ./...
	cd content/code && go build && go test ./...

# convert orgmode to markdown
ORG_FILES=$(shell find content/ | grep .org)
ORG_MD_FILES=$(patsubt %.md, %.org, $(ORG_FILES))
org2md: $(ORG_MD_FILES)
%.md: %.org
	# use emacs to compile.
	# cd to each target's directory, run there
	# this should invoke orient
	# this should produce hugo markdown output
	bin/org2hugomd.el <$< >$@


# installing deps

deps:
	bin/install-deps.sh
	@# make bins last, after installing other deps
	@# so we re-invoke make.
	make bins

# building our tools

bins: bin/codeGen

bin/codeGen: content/codeGen/*.go
	cd content/codeGen && go build -o ../../bin/codeGen

# other

serve: hugo-build
	hugo serve
