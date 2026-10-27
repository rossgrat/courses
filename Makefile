GO_VERSION := $(shell cat .go-version)

.PHONY: build vet

build:
	@test "$$(GOTOOLCHAIN=local go env GOVERSION)" = "go$(GO_VERSION)" || { printf '%s\n' 'Install Go $(GO_VERSION) before building.' >&2; exit 1; }
	GOTOOLCHAIN=local CGO_ENABLED=0 go build -trimpath -o bin/course ./cmd/course

vet:
	GOTOOLCHAIN=local go vet ./cmd/course
