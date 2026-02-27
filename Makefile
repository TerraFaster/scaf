.PHONY: build test lint fmt clean install

BINARY_NAME=scaf
MODULE=github.com/TerraFaster/scaf
VERSION=2.0.0

# Build
build:
	go build -ldflags="-X '$(MODULE)/cmd.Version=$(VERSION)'" -o bin/$(BINARY_NAME) .

# Build for all platforms
build-all:
	GOOS=linux   GOARCH=amd64 go build -o bin/$(BINARY_NAME)-linux-amd64 .
	GOOS=linux   GOARCH=arm64 go build -o bin/$(BINARY_NAME)-linux-arm64 .
	GOOS=darwin  GOARCH=amd64 go build -o bin/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin  GOARCH=arm64 go build -o bin/$(BINARY_NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build -o bin/$(BINARY_NAME)-windows-amd64.exe .

# Test
test:
	go test ./... -v

# Test with coverage
test-coverage:
	go test ./... -v -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func=coverage.out

# Lint
lint:
	golangci-lint run

# Format
fmt:
	gofmt -w .

# Clean
clean:
	rm -rf bin/ coverage.out coverage.html

# Install to $GOPATH/bin
install:
	go install -ldflags="-X '$(MODULE)/cmd.Version=$(VERSION)'" .

# Tidy
tidy:
	go mod tidy

# Verify
vet:
	go vet ./...
