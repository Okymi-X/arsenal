# arsenal Makefile
#
# Targets: build, test, lint, fmt, install, clean, release.

BINARY      := arsenal
PKG         := ./cmd/arsenal
BIN_DIR     := bin
PREFIX      ?= /usr/local
VERSION     := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS     := -s -w -X main.version=$(VERSION)

.PHONY: all build test lint fmt install clean release tidy registry registry-check verify-registry

all: build

## build: compile a static binary into bin/
build: registry-check
	mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY) $(PKG)

## registry: reassemble registry/registry.toml from registry/segments/*.toml
registry:
	go run ./tools/regbuild

## registry-check: fail if registry.toml is out of date with its segments
registry-check:
	go run ./tools/regbuild -verify

## test: run the unit and integration tests
test:
	go test ./...

## lint: run golangci-lint using the committed config
lint:
	golangci-lint run ./...

## fmt: format the codebase with gofumpt
fmt:
	gofumpt -w .

## tidy: tidy and verify module dependencies
tidy:
	go mod tidy
	go mod verify

## verify-registry: check every registry version exists at its official source
verify-registry:
	go run ./tools/regcheck -registry registry/registry.toml

## install: install the binary into PREFIX/bin
install: build
	install -d $(DESTDIR)$(PREFIX)/bin
	install -m 0755 $(BIN_DIR)/$(BINARY) $(DESTDIR)$(PREFIX)/bin/$(BINARY)

## clean: remove build artifacts
clean:
	rm -rf $(BIN_DIR)

## release: cross-compile static binaries for supported platforms
release: clean
	mkdir -p $(BIN_DIR)
	GOOS=linux  GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY)-linux-amd64  $(PKG)
	GOOS=linux  GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY)-linux-arm64  $(PKG)
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY)-darwin-amd64 $(PKG)
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY)-darwin-arm64 $(PKG)
