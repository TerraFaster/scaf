package initproject

import (
	"testing"
)

func TestGoFilesDefaultMode(t *testing.T) {
	opts := GoOptions{
		ModuleName: "github.com/test/myapp",
		Mode:       "default",
	}
	files := GoFiles(opts)
	if len(files) == 0 {
		t.Error("should return files")
	}

	// Check for expected files
	expectedPaths := map[string]bool{
		"Makefile":     false,
		"README.md":    false,
		"LICENSE":      false,
		".gitignore":   false,
		".editorconfig": false,
	}
	for _, f := range files {
		expectedPaths[f.Path] = true
	}
	for path, found := range expectedPaths {
		if !found {
			t.Errorf("expected file %q in GoFiles output", path)
		}
	}
}

func TestGoFilesWithDocker(t *testing.T) {
	opts := GoOptions{
		ModuleName: "github.com/test/myapp",
		Mode:       GoModeCLI,
		Docker:     true,
	}
	files := GoFiles(opts)

	hasDockerfile := false
	for _, f := range files {
		if f.Path == "Dockerfile" {
			hasDockerfile = true
			break
		}
	}
	if !hasDockerfile {
		t.Error("Docker flag should add Dockerfile")
	}
}

func TestNodeFilesDefault(t *testing.T) {
	opts := NodeOptions{
		ProjectName: "myproject",
		Mode:        NodeModeDefault,
	}
	files := NodeFiles(opts)
	if len(files) == 0 {
		t.Error("should return files")
	}
}

func TestNodeFilesTypeScript(t *testing.T) {
	opts := NodeOptions{
		ProjectName: "myproject",
		Mode:        NodeModeTS,
	}
	files := NodeFiles(opts)

	hasTsconfig := false
	hasTsEntry := false
	for _, f := range files {
		if f.Path == "tsconfig.json" {
			hasTsconfig = true
		}
		if f.Path == "src/index.ts" {
			hasTsEntry = true
		}
	}
	if !hasTsconfig {
		t.Error("TS mode should include tsconfig.json")
	}
	if !hasTsEntry {
		t.Error("TS mode should include src/index.ts")
	}
}
