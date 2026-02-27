package readme

import (
	"strings"
	"testing"
)

func TestRenderMinimal(t *testing.T) {
	params := Params{
		Name:        "myapp",
		Description: "A test project",
		Author:      "Alice",
		Template:    "minimal",
	}
	content, err := Render(params)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	if !strings.Contains(content, "myapp") {
		t.Error("expected project name in output")
	}
	if !strings.Contains(content, "A test project") {
		t.Error("expected description in output")
	}
}

func TestRenderAllTemplates(t *testing.T) {
	for name := range readmeTemplates {
		params := Params{
			Name:     "testproject",
			Template: name,
		}
		_, err := Render(params)
		if err != nil {
			t.Errorf("Render(%q) failed: %v", name, err)
		}
	}
}

func TestRenderFallbackToMinimal(t *testing.T) {
	params := Params{
		Name:     "myapp",
		Template: "nonexistent-template",
	}
	content, err := Render(params)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	if !strings.Contains(content, "myapp") {
		t.Error("fallback render should contain project name")
	}
}

func TestAutoFillName(t *testing.T) {
	p := Params{}
	p.AutoFill()
	if p.Name == "" {
		t.Error("AutoFill should set Name from cwd")
	}
}

func TestAvailableTemplates(t *testing.T) {
	available := AvailableTemplates()
	if len(available) == 0 {
		t.Error("should have at least one template")
	}
	for _, expected := range []string{"minimal", "cli", "go"} {
		found := false
		for _, a := range available {
			if a == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected template %q in available list", expected)
		}
	}
}
