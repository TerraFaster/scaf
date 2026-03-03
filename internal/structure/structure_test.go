package structure

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetKnownTemplate(t *testing.T) {
	tmpl, ok := Get("minimal")
	if !ok {
		t.Fatal("minimal template should exist")
	}
	if tmpl.Name != "minimal" {
		t.Error("name mismatch")
	}
}

func TestGetUnknownTemplate(t *testing.T) {
	_, ok := Get("nonexistent")
	if ok {
		t.Error("should return false for unknown template")
	}
}

func TestAvailableTemplates(t *testing.T) {
	available := AvailableTemplates()
	if len(available) == 0 {
		t.Error("should have templates")
	}
	for _, expected := range []string{"minimal", "cli", "layered", "ddd"} {
		found := false
		for _, a := range available {
			if a == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected %q in available templates", expected)
		}
	}
}

func TestApplyDryRun(t *testing.T) {
	tmpl, _ := Get("minimal")
	// Dry run should not create files
	err := Apply(tmpl, true)
	if err != nil {
		t.Fatalf("dry run failed: %v", err)
	}
}

func TestApplyCreatesFiles(t *testing.T) {
	tmpDir := t.TempDir()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working dir: %v", err)
	}
	defer func() {
		if err := os.Chdir(orig); err != nil {
			t.Fatalf("failed to restore working dir: %v", err)
		}
	}()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir to tmpDir: %v", err)
	}

	tmpl, _ := Get("minimal")
	if err := Apply(tmpl, false); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	// Check that src/ was created
	if _, err := os.Stat(filepath.Join(tmpDir, "src")); os.IsNotExist(err) {
		t.Error("src/ directory should have been created")
	}
}
