package cache

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCacheSetGet(t *testing.T) {
	dir := t.TempDir()
	c := &Cache{baseDir: dir, ttlHours: 168}

	if err := c.Set("ns", "key1", "hello world"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	val, ok, err := c.Get("ns", "key1")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if !ok {
		t.Fatal("expected cache hit")
	}
	if val != "hello world" {
		t.Fatalf("expected 'hello world', got %q", val)
	}
}

func TestCacheMiss(t *testing.T) {
	dir := t.TempDir()
	c := &Cache{baseDir: dir, ttlHours: 168}

	_, ok, err := c.Get("ns", "nonexistent")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if ok {
		t.Fatal("expected cache miss")
	}
}

func TestCacheTTLExpired(t *testing.T) {
	dir := t.TempDir()
	c := &Cache{baseDir: dir, ttlHours: 1}

	if err := c.Set("ns", "oldkey", "stale"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Backdate the file modification time
	path := c.keyPath("ns", "oldkey")
	past := time.Now().Add(-2 * time.Hour)
	if err := os.Chtimes(path, past, past); err != nil {
		t.Fatalf("Chtimes failed: %v", err)
	}

	_, ok, err := c.Get("ns", "oldkey")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if ok {
		t.Fatal("expected expired cache miss")
	}
}

func TestCachePathTraversal(t *testing.T) {
	dir := t.TempDir()
	c := &Cache{baseDir: dir, ttlHours: 168}

	err := c.Set("ns", "../../etc/passwd", "bad")
	if err != nil {
		// Expected: path traversal should be sanitized or blocked
		// The sanitizeKey function converts / to _, so it should be safe
		t.Logf("Set with traversal key returned error (acceptable): %v", err)
	}

	// Verify the key was sanitized - should not have written outside base dir
	dangerPath := filepath.Join(dir, "..", "..", "etc", "passwd")
	if _, err := os.Stat(dangerPath); err == nil {
		t.Fatal("path traversal was not prevented!")
	}
}

func TestCacheClear(t *testing.T) {
	dir := t.TempDir()
	c := &Cache{baseDir: dir, ttlHours: 168}

	_ = c.Set("ns", "k1", "v1")
	_ = c.Set("ns2", "k2", "v2")

	if err := c.Clear(); err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Fatal("expected cache dir to be removed after Clear")
	}
}

func TestSanitizeKey(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"mit", "mit"},
		{"apache-2.0", "apache-2.0"},
		{"../../etc", "______etc"},
		{"key with spaces", "key_with_spaces"},
	}
	for _, tt := range tests {
		got := sanitizeKey(tt.input)
		if got != tt.expected {
			t.Errorf("sanitizeKey(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}
