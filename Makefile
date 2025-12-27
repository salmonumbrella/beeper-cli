.PHONY: build test lint clean install

BINARY=beeper
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X github.com/salmonumbrella/beeper-cli/internal/cmd.Version=$(VERSION) -X github.com/salmonumbrella/beeper-cli/internal/cmd.BuildDate=$(BUILD_DATE)"

build:
	go build $(LDFLAGS) -o $(BINARY) ./cmd/beeper

test:
	go test -v ./...

lint:
	golangci-lint run

clean:
	rm -f $(BINARY)

install: build
	mv $(BINARY) ~/bin/
