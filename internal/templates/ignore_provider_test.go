package templates

import (
	"testing"
)

func TestParseCSVList(t *testing.T) {
	raw := "1c,1c-bitrix,a-frame,actionscript,ada\npython,node,go"
	items := parseCSVList(raw)
	if len(items) == 0 {
		t.Fatal("expected items from CSV list")
	}
	found := false
	for _, item := range items {
		if item == "python" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected 'python' in parsed list")
	}
}

func TestParseCSVListEmpty(t *testing.T) {
	items := parseCSVList("")
	if len(items) != 0 {
		t.Fatalf("expected empty list, got %d items", len(items))
	}
}

func TestIgnoreUserTemplateOverride(t *testing.T) {
	import_dir := t.TempDir()

	p := &IgnoreProvider{userDir: import_dir}

	// Non-existent user template should error
	_, err := p.getUserTemplate("nonexistent")
	if err == nil {
		t.Fatal("expected error for non-existent user template")
	}
}

func TestIgnoreUserTemplatePathTraversal(t *testing.T) {
	dir := t.TempDir()
	p := &IgnoreProvider{userDir: dir}

	_, err := p.getUserTemplate("../../etc/passwd")
	if err == nil {
		t.Fatal("expected path traversal to be blocked")
	}
}
