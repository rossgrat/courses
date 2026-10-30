GO_VERSION := $(shell cat system/cli/.go-version)
LATEXMK_VERSION := 4.88
PDFTEX_VERSION := 3.141592653-2.6-1.40.29 (TeX Live 2026)

.PHONY: build vet validate pdfs check-tex clean-pdfs

build:
	@test "$$(GOTOOLCHAIN=local go env GOVERSION)" = "go$(GO_VERSION)" || { printf '%s\n' 'Install Go $(GO_VERSION) before building.' >&2; exit 1; }
	cd system/cli && GOTOOLCHAIN=local CGO_ENABLED=0 go build -trimpath -o ../../build/bin/course .

vet:
	cd system/cli && GOTOOLCHAIN=local go vet .

validate: build
	bash course validate

pdfs: build
	@sources="$$(bash course pdf-sources)" || exit; \
	if [ -n "$$sources" ]; then \
		printf '%s\n' "$$sources" | while IFS= read -r source; do \
			$(MAKE) "$${source%.tex}.pdf" || exit; \
		done; \
	fi

check-tex:
	@test "$$(latexmk -v | sed -n 's/.*Version //p')" = '$(LATEXMK_VERSION)' || { printf '%s\n' 'Install latexmk $(LATEXMK_VERSION) before compiling PDFs.' >&2; exit 1; }
	@pdflatex --version | head -n 1 | grep -Fxq 'pdfTeX $(PDFTEX_VERSION)' || { printf '%s\n' 'Install pdfTeX $(PDFTEX_VERSION) before compiling PDFs.' >&2; exit 1; }

%.pdf: %.tex system/latex/curriculum.sty Makefile | check-tex
	mkdir -p build/pdf/$(*D)
	latexmk -norc -pdf -interaction=nonstopmode -halt-on-error -outdir=build/pdf/$(*D) $<
	cp build/pdf/$*.pdf $@

clean-pdfs:
	rm -rf build/pdf
