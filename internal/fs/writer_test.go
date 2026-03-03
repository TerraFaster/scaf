package fs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriterCreatesFile(t *testing.T) {
	dir := t.TempDir()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(orig); err != nil {
			t.Fatalf("failed to restore cwd: %v", err)
		}
	}()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	w := NewWriter()
	opts := WriteOptions{Force: true}
	if err := w.Write("LICENSE", "MIT License content", opts); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "LICENSE"))
	if err != nil {
		t.Fatalf("could not read file: %v", err)
	}
	if string(data) != "MIT License content" {
		t.Fatalf("unexpected content: %q", string(data))
	}
}

func TestWriterForceOverwrite(t *testing.T) {
	dir := t.TempDir()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(orig); err != nil {
			t.Fatalf("failed to restore cwd: %v", err)
		}
	}()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	w := NewWriter()
	opts := WriteOptions{Force: true}

	_ = w.Write("LICENSE", "original", opts)
	_ = w.Write("LICENSE", "overwritten", opts)

	data, _ := os.ReadFile(filepath.Join(dir, "LICENSE"))
	if string(data) != "overwritten" {
		t.Fatalf("expected overwritten content, got %q", string(data))
	}
}

func TestWriterBackupMode(t *testing.T) {
	dir := t.TempDir()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(orig); err != nil {
			t.Fatalf("failed to restore cwd: %v", err)
		}
	}()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	// Write original file
	originalPath := filepath.Join(dir, "LICENSE")
	if err := os.WriteFile(originalPath, []byte("original"), 0644); err != nil {
		t.Fatal(err)
	}

	w := NewWriter()
	opts := WriteOptions{Backup: true, Force: true}
	if err := w.Write("LICENSE", "new content", opts); err != nil {
		t.Fatalf("Write with backup failed: %v", err)
	}

	// Check backup exists
	bakData, err := os.ReadFile(filepath.Join(dir, "LICENSE.bak"))
	if err != nil {
		t.Fatalf("backup file not found: %v", err)
	}
	if string(bakData) != "original" {
		t.Fatalf("backup content wrong: %q", string(bakData))
	}

	// Check new content
	newData, _ := os.ReadFile(originalPath)
	if string(newData) != "new content" {
		t.Fatalf("new content wrong: %q", string(newData))
	}
}

func TestWriterPathTraversalBlocked(t *testing.T) {
	dir := t.TempDir()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(orig); err != nil {
			t.Fatalf("failed to restore cwd: %v", err)
		}
	}()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	w := NewWriter()
	opts := WriteOptions{Force: true}
	err = w.Write("../../etc/passwd", "malicious", opts)
	if err == nil {
		t.Fatal("expected path traversal to be blocked")
	}
}
