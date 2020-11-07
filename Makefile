SHELL = /bin/bash

NAME := octant-jx
ORG := jenkins-x
ORG_REPO := $(ORG)/$(NAME)
RELEASE_ORG_REPO := $(ORG_REPO)
REV := $(shell git rev-parse --short HEAD 2> /dev/null || echo 'unknown')
ROOT_PACKAGE := github.com/$(ORG_REPO)
BRANCH     := $(shell git rev-parse --abbrev-ref HEAD 2> /dev/null  || echo 'unknown')
BUILD_DATE := $(shell date +%Y%m%d-%H:%M:%S)
#GO_VERSION := $(shell $(GO) version | sed -e 's/^[^0-9.]*\([0-9.]*\).*/\1/')
GO_VERSION := 1.12
PLUGINS_DIR := $(HOME)/.config/octant/plugins

# set dev version unless VERSION is explicitly set via environment
#VERSION ?= $(shell echo "$$(git for-each-ref refs/tags/ --count=1 --sort=-version:refname --format='%(refname:short)' 2>/dev/null)-dev+$(REV)" | sed 's/^v//')


.PHONY: build
build:
	go build -o bin/octant-jx -ldflags " -X main.version=$(VERSION)" cmd/octant-jx/main.go

build-jxo:
	go build -o bin/octant-jxo -ldflags " -X main.version=$(VERSION)" cmd/octant-jxo/main.go

release: build build-jxo
	mkdir -p $(PLUGINS_DIR)
	cp bin/octant-jx $(PLUGINS_DIR)
	cp bin/octant-jxo $(PLUGINS_DIR)

.PHONY: dev
dev:
	scripts/dev.sh

.PHONY: linux
linux:

fmt:
	go fmt ./...

test:
	go test ./...

octant: release
	octant -v --disable-open-browser > octant.log 2>&1

octant-jx: build
	mkdir -p $(PLUGINS_DIR)
	cp bin/octant-jx $(PLUGINS_DIR)
	octant -v --disable-open-browser > octant.log 2>&1

tail:
	tail -f octant.log | grep "octant-jx"

tailo:
	tail -f octant.log | grep "octant-ojx"


.PHONY: goreleaser
goreleaser:
	step-go-releaser --organisation=$(ORG) --revision=$(REV) --branch=$(BRANCH) --build-date=$(BUILD_DATE) --go-version=$(GO_VERSION) --root-package=$(ROOT_PACKAGE) --version=$(VERSION)

.PHONY: clean
clean: ## Clean the generated artifacts
	rm -rf bin release dist

.PHONY: lint
lint: ## Lint the code
	./hack/gofmt.sh
	./hack/linter.sh