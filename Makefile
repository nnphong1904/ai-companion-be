.PHONY: build test lint clean run deps

BINARY_NAME=ai-companion-be
BUILD_DIR=bin
GO=go

build:
	$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server

test:
	$(GO) test -v -race -coverprofile=coverage.out ./...

test-coverage: test
	$(GO) tool cover -html=coverage.out

lint:
	golangci-lint run ./...

fmt:
	$(GO) fmt ./...

run:
	$(GO) run ./cmd/server

clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out

deps:
	$(GO) mod download
	$(GO) mod tidy
