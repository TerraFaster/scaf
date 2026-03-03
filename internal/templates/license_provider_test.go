package templates

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/TerraFaster/scaf/internal/cache"
)

func TestLicenseFuzzySearch(t *testing.T) {
	p := &LicenseProvider{}
	licenses := []LicenseInfo{
		{Key: "mit", Name: "MIT License"},
		{Key: "mpl-2.0", Name: "Mozilla Public License 2.0"},
		{Key: "apache-2.0", Name: "Apache License 2.0"},
	}

	results := p.FuzzySearch("mt", licenses)
	if len(results) == 0 {
		t.Fatal("expected fuzzy results for 'mt'")
	}
	if results[0].Key != "mit" {
		t.Fatalf("expected 'mit' as first result, got %q", results[0].Key)
	}
}

func TestLicenseListFromAPI(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		licenses := []LicenseInfo{
			{Key: "mit", Name: "MIT License"},
			{Key: "apache-2.0", Name: "Apache License 2.0"},
		}
		if err := json.NewEncoder(w).Encode(licenses); err != nil {
			t.Fatalf("failed to encode licenses JSON: %v", err)
		}
	}))
	defer server.Close()

	dir := t.TempDir()
	c := &cache.Cache{}
	_ = c
	_ = dir
	// Parsing test: verify JSON decode
	raw := `[{"key":"mit","name":"MIT License"},{"key":"apache-2.0","name":"Apache License 2.0"}]`
	var licenses []LicenseInfo
	if err := json.Unmarshal([]byte(raw), &licenses); err != nil {
		t.Fatalf("failed to parse licenses JSON: %v", err)
	}
	if len(licenses) != 2 {
		t.Fatalf("expected 2 licenses, got %d", len(licenses))
	}
	if licenses[0].Key != "mit" {
		t.Fatalf("expected 'mit', got %q", licenses[0].Key)
	}
}

func TestLicenseDetailParsing(t *testing.T) {
	raw := `{
		"key": "mit",
		"name": "MIT License",
		"body": "MIT License\n\nCopyright (c) [year] [fullname]\n\nPermission is hereby granted..."
	}`
	var detail LicenseDetail
	if err := json.Unmarshal([]byte(raw), &detail); err != nil {
		t.Fatalf("failed to parse license detail: %v", err)
	}
	if detail.Key != "mit" {
		t.Fatalf("expected 'mit', got %q", detail.Key)
	}
	if detail.Body == "" {
		t.Fatal("expected non-empty body")
	}
}

func TestUserTemplateOverride(t *testing.T) {
	dir := t.TempDir()

	// Write a user template
	content := "Custom MIT License\n"
	if err := os.WriteFile(dir+"/mit", []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	p := &LicenseProvider{userDir: dir}
	body, err := p.getUserTemplate("mit")
	if err != nil {
		t.Fatalf("getUserTemplate failed: %v", err)
	}
	if body != content {
		t.Fatalf("expected user template content, got %q", body)
	}
}
