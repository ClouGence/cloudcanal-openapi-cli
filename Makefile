GO ?= go
BIN ?= bin/cloudcanal
PKG ?= ./...
TEST_PKG ?= ./test/...
COVER_PROFILE ?= coverage.out
DIST_DIR ?= dist
VERSION ?= dev
COMMIT ?= $(shell git rev-parse HEAD 2>/dev/null || echo unknown)
BUILD_TIME ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GO_BUILD_FLAGS ?=
EXTRA_LDFLAGS ?=
BUILDINFO_PKG := github.com/ClouGence/cloudcanal-openapi-cli/internal/buildinfo
LDFLAGS ?= -X $(BUILDINFO_PKG).Version=$(VERSION) -X $(BUILDINFO_PKG).Commit=$(COMMIT) -X $(BUILDINFO_PKG).BuildTime=$(BUILD_TIME) $(EXTRA_LDFLAGS)

.PHONY: build test vet test-race cover ci release-assets clean

build:
	mkdir -p $(dir $(BIN))
	$(GO) build $(GO_BUILD_FLAGS) -ldflags "$(LDFLAGS)" -o $(BIN) ./cmd/cloudcanal

test:
	$(GO) test $(PKG)

vet:
	$(GO) vet $(PKG)

test-race:
	$(GO) test -race $(TEST_PKG)

cover:
	$(GO) test -coverpkg=$(PKG) -coverprofile=$(COVER_PROFILE) $(TEST_PKG)
	$(GO) tool cover -func=$(COVER_PROFILE)

ci: test vet test-race cover build

release-assets:
	./scripts/build_release_assets.sh

clean:
	rm -rf $(dir $(BIN)) $(DIST_DIR) $(COVER_PROFILE)
