ID_FILES=$(shell find . -name '*.id')
GEN_GO_FILES=$(patsubst %.id, %.gen.go, $(ID_FILES))

test: build

build: $(GEN_GO_FILES)
	go build -gcflags="-e" ./...

codeGen: $(shell find .. | grep .go | grep -v .gen.go)
	cd .. && go build -o test_cases/codeGen

%.gen.go: %.id codeGen
	-./codeGen $< $@

clean:
	find . | grep .gen.go | xargs rm
	rm codeGen
