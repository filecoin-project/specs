
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
	@echo "MAIN TARGETS"
	@echo "	make help        description of the targets (this message)"
	@echo "	make deps        install all dependencies of this tool chain"
	@echo "	make deps-user   install dependencies for user tooling"
	@echo "	make build       build all final artifacts (website only for now)"
	@echo "	make test        run all test cases (go-test only for now)"
	@echo "	make drafts      publish artifacts to ipfs and show an address"
	@echo "	make publish     publish final artifacts to spec website (github pages)"
	@echo "	make clean       removes all build artifacts. you shouldn't need this"
	@echo ""
	@echo "INTERMEDIATE TARGETS"
	@echo "	make website     build the website artifact"
	@#echo "	make pdf         build the pdf artifact"
	@echo "	make hugo-build  run the hugo part of the pipeline"
	@echo "	make gen-code    generate code artifacts (eg id -> go)"
	@echo "	make build-code  build all src go code (test it)"
	@echo "	make org2md      run org mode to markdown compilation"
	@echo "	make go-test     run test cases in code artifacts"
	@echo ""
	@echo "OTHER TARGETS"
	@echo "	make bins        compile some build tools whose source is in this repo"
	@echo "	make serve       start hugo in serving mode -- must run make build on changes manually"

# main Targets
build: build-code website

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

clean: .PHONY
	rm -rf build

clean-deps: .PHONY
	@echo "WARNING: this does not uninstall global packages, sorry."
	@echo "         If you would like to remove them, see bin/install-deps.sh"
	-rm -r deps
	-rm -r .slime
	-rm -r bin/.emacs

# intermediate targets
website: org2md hugo-build
	mkdir -p build/website
	-rm -rf build/website/*
	mv hugo/public/* build/website
	@echo TODO: add generate-code to this target

pdf: org2md hugo-build
	@echo TODO: add generate-code to this target
	bin/build-pdf.sh

hugo-build: hugo-src $(shell find hugo/content | grep '.md')
	cd hugo && hugo

hugo-src: $(shell find src | grep '.md')
	rm -rf hugo/content/docs
	cp -r src hugo/content/docs

hugo-src-rsync: $(shell find src | grep '.md')
	@mkdir -p hugo/content/docs
	rsync -av --inplace src/ hugo/content/docs
	echo "" >> hugo/content/_index.md # force reload
	echo "" >> hugo/content/menu/index.md # force reload

hugo-watch: .PHONY
	bin/watcher --cmd="make hugo-src-rsync" --startcmd src 2>/dev/null

orient: .PHONY
	bin/build-spec-orient.sh

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

bin/codeGen: $(shell find tools/codeGen | grep .go)
	cd tools/codeGen && go build -o ../../bin/codeGen

bin/watcher:
	go get -u github.com/radovskyb/watcher/...
	go build -o $@ github.com/radovskyb/watcher/cmd/watcher

# other

serve: hugo-build .PHONY
	echo "run `make website` and refresh to update"
	cd hugo && hugo serve --noHTTPCache

serve-website: website .PHONY
	# use this if `make serve` breaks
	echo "run `make website` and refresh to update"
	cd build/website && python -m SimpleHTTPServer 1313

serve-and-watch: serve hugo-watch
	echo "make sure you run this with `make -j2`"

.PHONY:

# code generation and building targets

GO_INPUT_FILES=$(shell find src -iname '*.go')
GO_OUTPUT_FILES=$(patsubst src/%.go, build/code/%.go, $(GO_INPUT_FILES))

GO_UTIL_INPUT_FILE=tools/codeGen/util/util.go
GO_UTIL_OUTPUT_FILE=build/code/util/util.go

$(GO_UTIL_OUTPUT_FILE): $(GO_UTIL_INPUT_FILE)
	mkdir -p $(dir $@)
	cp $< $@

build/code/%.go: src/%.go
	mkdir -p $(dir $@)
	cp $< $@

ID_FILES=$(shell find src -name '*.id')
GEN_GO_FILES=$(patsubst src/%.id, build/code/%.gen.go, $(ID_FILES))
build/code/%.gen.go: src/%.id bin/codeGen
	mkdir -p $(dir $@)
	-bin/codeGen $< $@

gen-code: bin/codeGen $(GEN_GO_FILES) $(GO_OUTPUT_FILES) $(GO_UTIL_OUTPUT_FILE)

build/code/go.mod: src/build_go.mod
	mkdir -p $(dir $@)
	cp $< $@

build-code: gen-code build/code/go.mod
	cd build/code && go build -gcflags="-e" ./...

test-code: build-code
	# testing should have the side effect that all go is compiled
	cd tools/codeGen && go build && go test ./...
	cd build/code && go build && go test ./...
