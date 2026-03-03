// Package initproject handles project initialization (go, node).
package initproject

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// GoMode defines the type of Go project.
type GoMode string

const (
	GoModeCLI     GoMode = "cli"
	GoModeService GoMode = "service"
	GoModeLibrary GoMode = "library"
)

// GoOptions holds options for Go project init.
type GoOptions struct {
	ModuleName string
	Mode       GoMode
	Docker     bool
	DryRun     bool
}

// FileEntry represents a file to create.
type FileEntry struct {
	Path    string
	Content string
}

// GoFiles returns all files to create for a Go project.
func GoFiles(opts GoOptions) []FileEntry {
	year := strconv.Itoa(time.Now().Year())
	name := filepath.Base(opts.ModuleName)
	if idx := strings.LastIndex(opts.ModuleName, "/"); idx >= 0 {
		name = opts.ModuleName[idx+1:]
	}

	files := []FileEntry{
		{
			Path: "cmd/" + name + "/main.go",
			Content: fmt.Sprintf(`package main

import "fmt"

func main() {
	fmt.Println("Hello from %s!")
}
`, name),
		},
		{
			Path:    "internal/.gitkeep",
			Content: "",
		},
		{
			Path:    "pkg/.gitkeep",
			Content: "",
		},
		{
			Path:    "configs/.gitkeep",
			Content: "",
		},
		{
			Path:    "scripts/.gitkeep",
			Content: "",
		},
		{
			Path: "Makefile",
			Content: fmt.Sprintf(`.PHONY: build test lint fmt clean

BINARY_NAME=%s
MODULE=%s

build:
	go build -o bin/$(BINARY_NAME) ./cmd/$(BINARY_NAME)/...

test:
	go test ./... -v -cover

lint:
	golangci-lint run

fmt:
	gofmt -w .
	goimports -w .

clean:
	rm -rf bin/

run:
	go run ./cmd/$(BINARY_NAME)/...

tidy:
	go mod tidy
`, name, opts.ModuleName),
		},
		{
			Path: "README.md",
			Content: fmt.Sprintf(`# %s

A Go project.

## Requirements

- Go 1.21+

## Build

`+"```"+`sh
make build
`+"```"+`

## Test

`+"```"+`sh
make test
`+"```"+`

## License

MIT
`, name),
		},
		{
			Path: "LICENSE",
			Content: fmt.Sprintf(`MIT License

Copyright (c) %s

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
`, year),
		},
		{
			Path: ".gitignore",
			Content: `# Go
/bin/
/dist/
/tmp/
*.test
*.out
coverage.txt
coverage.html

# OS
.DS_Store
Thumbs.db

# IDE
.idea/
.vscode/
`,
		},
		{
			Path: ".editorconfig",
			Content: `root = true

[*]
charset = utf-8
end_of_line = lf
insert_final_newline = true
trim_trailing_whitespace = true
indent_style = tab
indent_size = 4

[*.md]
trim_trailing_whitespace = false
`,
		},
	}

	if opts.Mode == GoModeCLI {
		files[0] = FileEntry{
			Path: "cmd/" + name + "/main.go",
			Content: fmt.Sprintf(`package main

import (
	"fmt"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	fmt.Println("Hello from %s!")
	return nil
}
`, name),
		}
	}

	if opts.Docker {
		files = append(files, FileEntry{
			Path: "Dockerfile",
			Content: fmt.Sprintf(`# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/%s ./cmd/%s/...

# Run stage
FROM gcr.io/distroless/static-debian12
COPY --from=builder /app/bin/%s /%s
ENTRYPOINT ["/%s"]
`, name, name, name, name, name),
		})
		files = append(files, FileEntry{
			Path: ".dockerignore",
			Content: `.git
*.md
/tmp/
*.test
`,
		})
	}

	return files
}

// RunGoModInit runs go mod init for the module.
func RunGoModInit(moduleName string) error {
	cmd := exec.Command("go", "mod", "init", moduleName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
