GO ?= go
BIN ?= bin/cloudcanal
PKG ?= ./...
TEST_PKG ?= ./test/...
COVER_PROFILE ?= coverage.out

.PHONY: build test vet test-race cover ci clean

build:
	mkdir -p $(dir $(BIN))
	$(GO) build -o $(BIN) ./cmd/cloudcanal

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

clean:
	rm -rf $(dir $(BIN)) $(COVER_PROFILE)
