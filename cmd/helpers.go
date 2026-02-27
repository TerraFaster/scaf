package cmd

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
)

func runCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return out.String(), nil
}

// autoDetectProject checks files/dirs in current directory to detect project types.
func autoDetectProject() ([]string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	var detected []string

	checks := []struct {
		pattern string
		result  string
	}{
		{"package.json", "node"},
		{"Cargo.toml", "rust"},
		{"go.mod", "go"},
		{"requirements.txt", "python"},
	}

	for _, c := range checks {
		if _, err := os.Stat(filepath.Join(cwd, c.pattern)); err == nil {
			detected = append(detected, c.result)
		}
	}

	// Check for *.csproj
	matches, _ := filepath.Glob(filepath.Join(cwd, "*.csproj"))
	if len(matches) > 0 {
		detected = append(detected, "dotnet")
	}

	// Check for Assets/ directory (Unity)
	if info, err := os.Stat(filepath.Join(cwd, "Assets")); err == nil && info.IsDir() {
		detected = append(detected, "unity")
	}

	return detected, nil
}
