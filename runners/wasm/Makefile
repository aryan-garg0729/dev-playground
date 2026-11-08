# WASI modules for the wasm grading runner. Built from pinned source (go.sum).
PKGS := bufio bytes cmp container/heap container/list errors fmt iter maps math math/bits math/rand os regexp slices sort strconv strings sync time unicode unicode/utf8

dist: dist/go.wasm dist/js.wasm

dist/%.wasm: $(shell find cmd symbols -type f) go.mod go.sum
	GOOS=wasip1 GOARCH=wasm go build -trimpath -ldflags=-s -o $@ ./cmd/$*

symbols: ## Regenerate the Go lesson standard library (after changing PKGS or Go version)
	rm -f symbols/*-*.go $(addprefix symbols/,$(addsuffix .go,$(filter-out symbols,$(PKGS))))
	cd symbols && go run github.com/traefik/yaegi/cmd/yaegi@v0.16.3 extract -name symbols $(PKGS)

.PHONY: dist symbols
