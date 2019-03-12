
now=$(shell date -u "+%Y-%m-%d")

daily: cover.jpg
	@mkdir -p build
	gitbook pdf . "build/filecoin-spec.$(now).pdf"

bikeshed.prepare:
	rm -f spec.bs
	for f in _pre.md INTRO.md process.md operation.md sync.md validation.md storage-market.md retrieval-market.md payments.md mining.md expected-consensus.md state-machine.md actors.md faults.md signatures.md proofs.md network-protocols.md data-propagation.md data-structures.md local-storage.md definitions.md; \
		do (cat "$${f}"; echo) >> spec.bs; \
	done

bikeshed.local: bikeshed.prepare
	bikeshed spec -f spec.bs
	mv spec.html docs/index.html

bikeshed: bikeshed.prepare
	curl https://api.csswg.org/bikeshed/ -F file=@spec.bs -F force=1 > docs/index.html

cover.jpg:
	cover/make-today-cover

.PHONY: cover.jpg bikeshed bikeshed.prepare bikeshed.local
