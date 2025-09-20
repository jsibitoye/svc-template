APP_NAME := svc-template
PKG := github.com/jsibitoye/svc-template
BIN := bin/$(APP_NAME)
VERSION ?= 0.1.0
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo dev)
DATE    ?= $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

GOFLAGS := -trimpath
LDFLAGS := -s -w -X $(PKG)/internal/version.Version=$(VERSION) -X $(PKG)/internal/version.GitCommit=$(COMMIT) -X $(PKG)/internal/version.BuildDate=$(DATE)

.PHONY: all run test lint build clean

all: test build

run:
	@PORT=8080 go run ./cmd/api

test:
	go test ./...

lint:
	@golangci-lint run

build:
	CGO_ENABLED=0 go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BIN) ./cmd/api

clean:
	rm -rf bin dist coverage.out
