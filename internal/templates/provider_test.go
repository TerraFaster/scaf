package templates

import (
	"strings"
	"testing"
)

// Tests for embedded templates (offline, no network required)

func TestEmbeddedLicensesHasMIT(t *testing.T) {
	body, ok := embeddedLicenses["mit"]
	if !ok {
		t.Fatal("embedded MIT license not found")
	}
	if !strings.Contains(body, "MIT License") {
		t.Error("MIT body should contain 'MIT License'")
	}
	if !strings.Contains(body, "[year]") {
		t.Error("MIT body should contain [year] placeholder")
	}
}

func TestEmbeddedLicensesHasApache(t *testing.T) {
	_, ok := embeddedLicenses["apache-2.0"]
	if !ok {
		t.Error("embedded apache-2.0 not found")
	}
}

func TestEmbeddedGitignoreHasGo(t *testing.T) {
	body, ok := embeddedGitignore["go"]
	if !ok {
		t.Fatal("embedded go gitignore not found")
	}
	if !strings.Contains(body, "*.test") {
		t.Error("go gitignore should contain *.test")
	}
}

func TestEmbeddedGitignoreHasNode(t *testing.T) {
	body, ok := embeddedGitignore["node"]
	if !ok {
		t.Fatal("embedded node gitignore not found")
	}
	if !strings.Contains(body, "node_modules") {
		t.Error("node gitignore should contain node_modules")
	}
}

func TestEmbeddedLicensesCount(t *testing.T) {
	if len(embeddedLicenses) == 0 {
		t.Error("should have at least one embedded license")
	}
}

func TestEmbeddedGitignoreCount(t *testing.T) {
	if len(embeddedGitignore) == 0 {
		t.Error("should have at least one embedded gitignore template")
	}
}
