package config

import (
	"os"
	"path/filepath"
	"testing"
)

// overrideHome temporarily redirects UserHomeDir to a temp dir for testing.
func overrideHome(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("HOME", dir)         // Unix
	t.Setenv("USERPROFILE", dir)  // Windows
	return dir
}

func TestDefaultsAreSet(t *testing.T) {
	cfg := Defaults()
	if cfg.CacheTTLHours != 168 {
		t.Fatalf("expected default TTL 168, got %d", cfg.CacheTTLHours)
	}
	if !cfg.Interactive {
		t.Fatal("expected interactive to default to true")
	}
}

func TestLoadAutoCreatesConfigOnFirstRun(t *testing.T) {
	home := overrideHome(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() should not error: %v", err)
	}
	if cfg.CacheTTLHours != 168 {
		t.Fatalf("expected default TTL, got %d", cfg.CacheTTLHours)
	}

	// Config file must have been created automatically.
	expected := filepath.Join(home, ".scaf", FileName)
	if _, err := os.Stat(expected); err != nil {
		t.Fatalf("config file was not auto-created at %s: %v", expected, err)
	}
}

func TestLoadReadsExistingConfig(t *testing.T) {
	home := overrideHome(t)

	// Pre-write a config file.
	dir := filepath.Join(home, ".scaf")
	_ = os.MkdirAll(dir, 0755)
	content := "default_license: apache-2.0\ncache_ttl_hours: 24\n"
	_ = os.WriteFile(filepath.Join(dir, FileName), []byte(content), 0644)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}
	if cfg.DefaultLicense != "apache-2.0" {
		t.Fatalf("expected 'apache-2.0', got %q", cfg.DefaultLicense)
	}
	if cfg.CacheTTLHours != 24 {
		t.Fatalf("expected TTL 24, got %d", cfg.CacheTTLHours)
	}
}

func TestSaveAndLoadRoundtrip(t *testing.T) {
	overrideHome(t)

	original := Config{
		DefaultLicense: "mit",
		DefaultIgnore:  []string{"python", "node"},
		CacheTTLHours:  48,
		Interactive:    false,
		DefaultAuthor:  "Author",
	}

	if err := Save(original); err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load() after Save() failed: %v", err)
	}
	if loaded.DefaultLicense != "mit" {
		t.Fatalf("expected 'mit', got %q", loaded.DefaultLicense)
	}
	if len(loaded.DefaultIgnore) != 2 {
		t.Fatalf("expected 2 ignore entries, got %d", len(loaded.DefaultIgnore))
	}
	if loaded.CacheTTLHours != 48 {
		t.Fatalf("expected TTL 48, got %d", loaded.CacheTTLHours)
	}
	if loaded.DefaultAuthor != "Author" {
		t.Fatalf("expected 'Author', got %q", loaded.DefaultAuthor)
	}
}

func TestPathIsUnderHomeScaf(t *testing.T) {
	home := overrideHome(t)

	path, err := Path()
	if err != nil {
		t.Fatalf("Path() failed: %v", err)
	}

	expected := filepath.Join(home, ".scaf", FileName)
	if path != expected {
		t.Fatalf("expected %q, got %q", expected, path)
	}
}

func TestLoadIsIdempotent(t *testing.T) {
	overrideHome(t)

	// Two consecutive loads must not error and must return the same values.
	cfg1, err := Load()
	if err != nil {
		t.Fatalf("first Load() failed: %v", err)
	}
	cfg2, err := Load()
	if err != nil {
		t.Fatalf("second Load() failed: %v", err)
	}
	if cfg1.CacheTTLHours != cfg2.CacheTTLHours {
		t.Fatal("successive Load() calls returned different TTL values")
	}
}
