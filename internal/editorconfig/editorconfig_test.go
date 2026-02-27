package editorconfig

import (
	"strings"
	"testing"
)

func TestRenderDefault(t *testing.T) {
	p := Params{Variant: "default"}
	content := Render(p)
	if !strings.Contains(content, "root = true") {
		t.Error("should contain 'root = true'")
	}
	if !strings.Contains(content, "indent_style") {
		t.Error("should contain indent_style")
	}
}

func TestRenderGoVariant(t *testing.T) {
	p := Params{Variant: "go"}
	content := Render(p)
	if !strings.Contains(content, "indent_style = tab") {
		t.Error("Go variant should use tab indent")
	}
}

func TestRenderNodeVariant(t *testing.T) {
	p := Params{Variant: "node"}
	content := Render(p)
	if !strings.Contains(content, "indent_size = 2") {
		t.Error("Node variant should use indent_size = 2")
	}
}

func TestRenderUnknownFallsToDefault(t *testing.T) {
	p := Params{Variant: "nonexistent"}
	content := Render(p)
	if content == "" {
		t.Error("should return default for unknown variant")
	}
}

func TestRenderWithOverride(t *testing.T) {
	p := Params{
		Variant:    "default",
		IndentSize: "2",
	}
	content := Render(p)
	if !strings.Contains(content, "indent_size = 2") {
		t.Errorf("override not applied; content:\n%s", content)
	}
}

func TestAllVariants(t *testing.T) {
	for name := range variants {
		p := Params{Variant: name}
		content := Render(p)
		if content == "" {
			t.Errorf("variant %q produced empty content", name)
		}
	}
}
