.PHONY: build install clean test fmt lint

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X github.com/CromulentConsulting/ygm-cli/internal/cmd.Version=$(VERSION)"

build:
	go build $(LDFLAGS) -o bin/ygm ./cmd/ygm

install:
	go install $(LDFLAGS) ./cmd/ygm

clean:
	rm -rf bin/

test:
	go test ./...

fmt:
	go fmt ./...

lint:
	golangci-lint run

# Development: build and run
run: build
	./bin/ygm $(ARGS)

# Cross-compilation targets
build-all: build-darwin-amd64 build-darwin-arm64 build-linux-amd64 build-linux-arm64 build-windows-amd64

build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/ygm-darwin-amd64 ./cmd/ygm

build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/ygm-darwin-arm64 ./cmd/ygm

build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/ygm-linux-amd64 ./cmd/ygm

build-linux-arm64:
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o bin/ygm-linux-arm64 ./cmd/ygm

build-windows-amd64:
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/ygm-windows-amd64.exe ./cmd/ygm
