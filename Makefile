BINARY=scaf
VERSION=0.1.0
BUILD_DIR=dist

.PHONY: build test clean install cross

build:
	go build -ldflags="-X main.version=$(VERSION)" -o $(BINARY) .

test:
	go test ./...

test-verbose:
	go test -v ./...

clean:
	rm -f $(BINARY)
	rm -rf $(BUILD_DIR)

install: build
	cp $(BINARY) /usr/local/bin/$(BINARY)

# Cross-platform builds
cross:
	mkdir -p $(BUILD_DIR)
	GOOS=linux   GOARCH=amd64  go build -o $(BUILD_DIR)/$(BINARY)-linux-amd64   .
	GOOS=linux   GOARCH=arm64  go build -o $(BUILD_DIR)/$(BINARY)-linux-arm64   .
	GOOS=darwin  GOARCH=amd64  go build -o $(BUILD_DIR)/$(BINARY)-darwin-amd64  .
	GOOS=darwin  GOARCH=arm64  go build -o $(BUILD_DIR)/$(BINARY)-darwin-arm64  .
	GOOS=windows GOARCH=amd64  go build -o $(BUILD_DIR)/$(BINARY)-windows-amd64.exe .
	GOOS=windows GOARCH=arm64  go build -o $(BUILD_DIR)/$(BINARY)-windows-arm64.exe .

lint:
	golangci-lint run ./...

tidy:
	go mod tidy

.DEFAULT_GOAL := build
