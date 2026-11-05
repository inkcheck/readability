.DEFAULT_GOAL := help

BINARY := readability
TAG := $(shell git describe --tags --exact-match HEAD 2>/dev/null)

.PHONY: help build install test vet fmt check snapshot release clean

help: ## Show this help
	@echo "Usage: make <target>"
	@echo
	@awk 'BEGIN {FS = ":.*## "} /^[a-zA-Z_-]+:.*## / {printf "  %-10s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the CLI into ./readability
	go build -o $(BINARY) ./cmd/readability

install: ## Install the CLI into $GOPATH/bin
	go install ./cmd/readability

test: ## Run the tests
	go test ./...

vet: ## Run go vet
	go vet ./...

fmt: ## Format the code with gofmt
	gofmt -w .

check: vet test ## Run vet and tests

snapshot: check ## Build release archives locally without signing or publishing
	goreleaser release --snapshot --clean

release: check ## Push the tag at HEAD, then sign, notarise and publish with GoReleaser
	@test -z "$$(git status --porcelain)" || { echo "Working tree is not clean."; exit 1; }
	@test "$$(git rev-parse --abbrev-ref HEAD)" = main || { echo "Not on main."; exit 1; }
	@test -n "$(TAG)" || { echo "HEAD has no tag. Run: git tag -a vX.Y.Z -m 'vX.Y.Z'"; exit 1; }
	git push origin main $(TAG)
	GITHUB_TOKEN=$${GITHUB_TOKEN:-$$(gh auth token)} goreleaser release --clean

clean: ## Remove build output
	rm -rf $(BINARY) dist/
