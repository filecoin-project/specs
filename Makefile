
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
	@echo "	make deps-basic  run this once, to install & build basic dependencies"
	@echo "	make build       run this every time you want to re-build artifacts"
	@echo ""
	@echo "MAIN TARGETS"
	@echo "	make help        description of the targets (this message)"
	@echo "	make build       build all final artifacts (website only for now)"
	@echo "	make test        run all test cases (test-code only for now)"
	@echo "	make drafts      publish artifacts to ipfs and show an address"
	@echo "	make publish     publish final artifacts to spec website (github pages)"
	@echo "	make clean       removes all build artifacts. you shouldn't need this"
	@echo "	make serve       start hugo in serving mode -- must run 'make build' on changes manually"
	@echo ""
	@echo "INSTALL DEPENDENCIES"
	@echo "	make deps        install ALL dependencies of this tool chain"
	@echo "	make deps-basic  install minimal dependencies of this tool chain"
	@echo "	make deps-diag   install dependencies for rendering diagrams"
	@echo "	make deps-orient install dependencies for running orient"
	@echo "	make deps-ouser  install dependencies for orient user-environment tooling"
	@echo "	make bins        compile some build tools whose source is in this repo"
	@echo ""
	@echo "INTERMEDIATE TARGETS"
	@echo "	make website     build the website artifact"
	@#echo "	make pdf         build the pdf artifact"
	@echo "	make diagrams    build diagram artifacts ({dot, mmd} -> svg)"
	@echo "	make org2md      run org mode to markdown compilation"
	@echo ""
	@echo "HUGO TARGETS"
	@echo "	make hugo-src    copy sources into hugo dir"
	@echo "	make build-hugo  run the hugo part of the pipeline"
	@echo "	make watch-hugo  watch and rebuild hugo"
	@echo ""
	@echo "CODE TARGETS"
	@echo "	make gen-code    generate code artifacts (eg id -> go)"
	@echo "	make test-code   run test cases in code artifacts"
	@echo "	make build-code  build all src go code (test it)"
	@echo "	make clean-code  remove build code artifacts"
	@echo "	make watch-code  watch and rebuild code"
	@echo ""
	@echo "CLEAN TARGETS"
	@echo "	make clean       remove all build artifacts"
	@echo "	make clean-deps  remove (some of) the dependencies installed in this repo"
	@echo "	make clean-hugo  remove intermediate hugo artifacts"
	@echo "	make clean-code  remove build code artifacts"
	@echo ""
	@echo "WATCH TARGETS"
	@echo "	make serve-and-watch -j2  serve, watch, and rebuild all - works for live edit"
	@echo "	make watch-code           watch and rebuild code"
	@echo "	make watch-hugo           watch and rebuild hugo"
	@echo ""

# main Targets
build: diagrams build-code website

test: test-code test-codeGen

drafts: website
	bin/publish-to-ipfs.sh

publish: website
	bin/publish-to-gh-pages.sh

clean: .PHONY
	rm -rf build

# install dependencies

deps: deps-basic deps-diag deps-orient
	@# make bins last, after installing other deps
	@# so we re-invoke make.
	make bins

deps-basic:
	bin/install-deps-basic.sh -y

deps-ouser:
	bin/install-deps-orient-user.sh -y

deps-orient: submodules
	bin/install-deps-orient.sh -y

deps-diag:
	bin/install-deps-diagrams.sh -y

submodules:
	git submodule update --init --recursive

clean-deps: .PHONY
	@echo "WARNING: this does not uninstall global packages, sorry."
	@echo "         If you would like to remove them, see bin/install-deps.sh"
	-rm -r deps
	-git checkout ./deps/package.json
	-rm -r .slime
	-rm -r bin/.emacs

# intermediate targets
# NOTE: For now, disable org2md — must manually generate until batch-mode build issues are resolved.
website: diagrams build-hugo # org2md
	mkdir -p build/website
	-rm -rf build/website/*
	mv hugo/public/* build/website
	@echo TODO: add generate-code to this target

pdf: diagrams build-hugo # org2md
	@echo TODO: add generate-code to this target
	bin/build-pdf.sh

build-hugo: hugo-src $(shell find hugo/content | grep '.md')
	cd hugo && hugo

hugo-src: $(shell find src | grep '.md')
	rm -rf hugo/content/docs
	cp -r src hugo/content/docs
	# ox-hugo exports to src/content, so we need to copy that also.
	cp -r src/content/ hugo/content/docs
	mkdir -p hugo/content/ox-hugo
	cp src/static/ox-hugo/* hugo/content/ox-hugo

# this is used to get "serve-and-watch" working. trick is to use
hugo-src-rsync: $(shell find src | grep '.md') gen-code diagrams
	@mkdir -p hugo/content/docs
	rsync -av --inplace src/ hugo/content/docs
	printf " " >> hugo/content/_index.md # force reload
	printf " " >> hugo/content/menu/index.md # force reload

watch-hugo: .PHONY
	bin/watcher --cmd="make hugo-src-rsync" --startcmd src 2>/dev/null

clean-hugo: .PHONY
	rm -rf hugo/content/docs

orient: .PHONY
	bin/build-spec-orient.sh

# convert orgmode to markdown
ORG_FILES=$(shell find src | grep .org)
ORG_MD_FILES=$(patsubst %.org, %.md, $(ORG_FILES))
org2md: $(ORG_MD_FILES)
%.md: %.org
	# use emacs to compile.
	# cd to each target's directory, run there
	# this should invoke orient
	# this should produce hugo markdown output

        # Skip this until batch-mode build issues are resolve.
	# bin/org2hugomd.el <$< >$@


# building our tools
bins: bin/codeGen

bin/codeGen: $(shell find tools/codeGen | grep .go)
	cd tools/codeGen && go build -o ../../bin/codeGen

bin/watcher:
	go get -u github.com/radovskyb/watcher/...
	go build -o $@ github.com/radovskyb/watcher/cmd/watcher

test-codeGen: bin/codeGen
	cd tools/codeGen && go build && go test ./...

# other

serve: build-hugo .PHONY
	echo "run 'make website' and refresh to update"
	cd hugo && hugo serve --noHTTPCache

serve-website: website .PHONY
	# use this if `make serve` breaks
	echo "run 'make website' and refresh to update"
	cd build/website && python -m SimpleHTTPServer 1313

serve-and-watch: serve watch-hugo
	echo "make sure you run this with 'make -j2'"

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
	-bin/codeGen gen $< $@

gen-code: bin/codeGen build/code/go.mod $(GEN_GO_FILES) $(GO_OUTPUT_FILES) $(GO_UTIL_OUTPUT_FILE)

build/code/go.mod: src/build_go.mod
	mkdir -p $(dir $@)
	cp $< $@

build-code: gen-code
	@cd build/code && go build -gcflags="-e" ./...

test-code: build-code
	cd build/code && go test ./...

clean-code:
	rm -rf build/code

fmt-code:
	go fmt ./src/...
	go fmt ./build/code/...
	bin/codeGen fmt ./src/...

watch-code: .PHONY
	bin/watcher --cmd="make gen-code" --startcmd src 2>/dev/null

## diagrams

DOTs=$(shell find src -name '*.dot')
MMDs=$(shell find src -name '*.mmd')
SVGs=$(DOTs:%=%.svg) $(MMDs:%=%.svg)

diagrams: ${SVGs}

watch-diagrams: diagrams
	bin/watcher --cmd="make diagrams" --startcmd src 2>/dev/null

%.dot.svg: %.dot
	@which dot >/dev/null || echo "requires dot (graphviz) -- run make deps" && exit
	dot -Tsvg $< >$@

%.mmd.svg: %.mmd %.mmd.css
	deps/node_modules/.bin/mmdc -i $< -o $@ --cssFile $(word 2,$^)
