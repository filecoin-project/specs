
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
	@echo "filecoin spec build toolchain commands"
	@echo ""
	@echo "USAGE"
	@echo "\tmake deps        # run this once, to install dependencies"
	@echo "\tmake build       # run this every time you want to re-build artifacts"
	@echo ""
	@echo "WARNING"
	@echo "\tthis build tool is WIP, so some targets may not work yet"
	@echo "\tthis should stabilize in the next couple of days"
	@echo ""
	@echo "MAIN TARGETS"
	@echo "\tmake help        description of the targets (this message)"
	@echo "\tmake deps        install all dependencies of this tool chain"
	@echo "\tmake build       build all final artifacts (website only for now)"
	@echo "\tmake test        run all test cases (go-test only for now)"
	@echo "\tmake drafts      publish artifacts to ipfs and show an address"
	@echo "\tmake publish     publish final artifacts to spec website (github pages)"
	@echo ""
	@echo "INTERMEDIATE TARGETS"
	@echo "\tmake website     build the website artifact"
	@#echo "\tmake pdf         build the pdf artifact"
	@echo "\tmake hugo-build  run the hugo part of the pipeline"
	@echo "\tmake gen-code    generate code artifacts (eg ipld -> go)"
	@echo "\tmake org2md      run org mode to markdown compilation"
	@echo "\tmake go-test     run test cases in code artifacts"
	@echo ""
	@echo "OTHER TARGETS"
	@echo "\tmake bins        compile some build tools whose source is in this repo"
	@echo "\tmake serve       start hugo in serving mode -- must run make build on changes manually"

# main Targets
build: website

deps:
	bin/install-deps.sh
	@# make bins last, after installing other deps
	@# so we re-invoke make.
	make bins

drafts: website
	bin/publish-to-ipfs.sh

publish: website
	bin/publish-to-gh-pages.sh

# intermediate targets
website: go-test org2md hugo-build
	@echo TODO: add generate-code to this target

pdf: go-test org2md hugo-build
	@echo TODO: add generate-code to this target
	bin/build-pdf.sh

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


# building our tools
bins: bin/codeGen

bin/codeGen: content/codeGen/*.go
	cd content/codeGen && go build -o ../../bin/codeGen

# other

serve: hugo-build
	hugo serve
