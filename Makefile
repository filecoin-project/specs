
# guidelines for editing this makefile:
#
# - keep it simple -- put complicated commands into scripts inside bin/ (eg install-deps.sh)
# - document targets in the 'help' target
# - distinguish main targets (meant for users) from intermediate targets
# - if you write a new tool that requires compilation:
#      add a compilation target here and move the binary into bin/
# - if you add a dependency on another tool:
#      make sure you edit install-deps.sh to install or prompt to install it
# - keep diagrams/builsys/buildsys.dot in sync with the targets here
#      that is a diagram that is meant to make it easy to understand everything here.

help:
	@echo "SYNOPSIS"
	@echo "	make -- filecoin spec build toolchain commands"
	@echo ""
	@echo "USAGE"
	@echo "	make deps        # run this once, to install & build dependencies"
	@echo "	make build       # run this every time you want to re-build artifacts"
	@echo ""
	@echo "WARNING !"
	@echo "	this build tool is WIP, so some targets may not work yet"
	@echo "	this should stabilize in the next couple of days"
	@echo ""
	@echo "MAIN TARGETS"
	@echo "	make help        description of the targets (this message)"
	@echo "	make deps        install all dependencies of this tool chain"
	@echo "	make deps-user   install dependencies for user tooling"
	@echo "	make build       build all final artifacts (website only for now)"
	@echo "	make test        run all test cases (go-test only for now)"
	@echo "	make drafts      publish artifacts to ipfs and show an address"
	@echo "	make publish     publish final artifacts to spec website (github pages)"
	@echo ""
	@echo "INTERMEDIATE TARGETS"
	@echo "	make website     build the website artifact"
	@#echo "	make pdf         build the pdf artifact"
	@echo "	make hugo-build  run the hugo part of the pipeline"
	@echo "	make gen-code    generate code artifacts (eg id -> go)"
	@echo "	make org2md      run org mode to markdown compilation"
	@echo "	make go-test     run test cases in code artifacts"
	@echo ""
	@echo "OTHER TARGETS"
	@echo "	make bins        compile some build tools whose source is in this repo"
	@echo "	make serve       start hugo in serving mode -- must run make build on changes manually"

# main Targets
build: website

deps: submodules
	bin/install-deps.sh
	@# make bins last, after installing other deps
	@# so we re-invoke make.
	make bins

submodules:
	git submodule update --init --recursive

deps-user: deps
	bin/install-deps-orient-user.sh

drafts: website
	bin/publish-to-ipfs.sh

publish: website
	bin/publish-to-gh-pages.sh

# intermediate targets
website: go-test org2md hugo-build
	mkdir -p build/website
	-rm -rf build/website/*
	mv hugo/public/* build/website
	@echo TODO: add generate-code to this target

pdf: go-test org2md hugo-build
	@echo TODO: add generate-code to this target
	bin/build-pdf.sh

hugo-build: hugo-src $(shell find hugo/content | grep '.md')
	cd hugo && hugo

hugo-src: $(shell find src)
	rm -rf hugo/content/docs
	cp -r src hugo/content/docs

orient: src/orient/*
	bin/build-spec-orient.sh

ID_FILES=$(shell find src/ -name '*.id')
GEN_GO_FILES=$(patsubst %.id, %.gen.go, $(ID_FILES))
%.gen.go: %.id bin/codeGen
	bin/codeGen $<

gen-code: $(GEN_GO_FILES)

go-test: $(shell find hugo/content/ | grep .go)
	# testing should have the side effect that all go is compiled
	cd hugo/content/codeGen && go build && go test ./...
	# cd hugo/content/code && go build && go test ./...

# convert orgmode to markdown
ORG_FILES=$(shell find hugo/content | grep .org)
ORG_MD_FILES=$(patsubt %.md, %.org, $(ORG_FILES))
org2md: $(ORG_MD_FILES)
%.md: %.org
	# use emacs to compile.
	# cd to each target's directory, run there
	# this should invoke orient
	# this should produce hugo markdown output
	bin/org2hugomd.el <$< >$@


# building our tools
bins: bin/codeGen

bin/codeGen: hugo/content/codeGen/*.go
	cd hugo/content/codeGen && go build -o ../../../bin/codeGen

# other

serve: hugo-build .PHONY
	echo "run `make website` and refresh to update"
	cd hugo && hugo serve

serve-website: website .PHONY
	# use this if `make serve` breaks
	echo "run `make website` and refresh to update"
	cd build/website && python -m SimpleHTTPServer 1313


.PHONY:
