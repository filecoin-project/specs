
now=$(shell date -u "+%Y-%m-%d")

daily: cover.jpg
	@mkdir -p build
	gitbook pdf . "build/filecoin-spec.$(now).pdf"

pdf:
	./build-pdf.sh

dev:
	hugo serve

publish:
	git submodule update --init --recursive
	./publish.sh

cover.jpg:
	cover/make-today-cover

.PHONY: cover.jpg publish pdf
