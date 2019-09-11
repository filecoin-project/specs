
build:
	# todo

deps:
	# require emacs
	# require hugo
	# require go 1.12+
	# require cask (emacs deps)
	#
	# make bins last, after installing other deps
	# not all deps can fit nicely into make's dependency graph
	# so we re-invoke make.
	make bins


# main Targets

# this bundles the website
website: hugo-build

publish: website
	bin/publish-to-ipfs

# intermediate targets

hugo-build: $(shell find . | grep .md)
	hugo

generate-code: $(shell find content/ | grep .ipld)
	bin/codeGen <input-files?>

go-test: $(shell find content/ | grep .go)
	# testing should have the side effect that all go is compiled
	go test ./content/... # todo: test?

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


# building our tools

bins: bin/codeGen

bin/codeGen: content/codeGen/*.go
	go build -o bin/codeGen content/codeGen
