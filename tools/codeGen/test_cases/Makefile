ID_FILES=$(shell find . -name '*.id')
GEN_GO_FILES=$(patsubst %.id, %.gen.go, $(ID_FILES))

test: build

build: $(GEN_GO_FILES)
	go build -gcflags="-e" ./...

codeGen: $(shell find .. -name '*.go' -path test_cases -prune)
	cd .. && go build -o test_cases/codeGen

%.gen.go: %.id codeGen
	-./codeGen $< $@
