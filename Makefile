BINARY_NAME=slackcli
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-s -w -X github.com/jackchuka/slackcli/internal/version.Version=$(VERSION) -X github.com/jackchuka/slackcli/internal/version.Commit=$(COMMIT) -X github.com/jackchuka/slackcli/internal/version.BuildDate=$(BUILD_DATE)"

.PHONY: build test lint clean install

build:
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/slackcli

install:
	go install $(LDFLAGS) ./cmd/slackcli

generate:
	go generate ./...

test:
	go test ./... -race -count=1

test-coverage:
	go test ./... -race -coverprofile=coverage.out -covermode=atomic
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run ./...

clean:
	rm -rf bin/ coverage.out coverage.html

fmt:
	gofmt -s -w .

vet:
	go vet ./...
