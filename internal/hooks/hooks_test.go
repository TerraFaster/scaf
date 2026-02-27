package hooks

import (
	"strings"
	"testing"
)

func TestRenderLight(t *testing.T) {
	content := Render(ModeLight, "go")
	if !strings.Contains(content, "#!/bin/sh") {
		t.Error("hook should start with shebang")
	}
	if !strings.Contains(content, "gofmt") {
		t.Error("go light hook should include gofmt check")
	}
}

func TestRenderStrict(t *testing.T) {
	content := Render(ModeStrict, "go")
	if !strings.Contains(content, "golangci-lint") {
		t.Error("go strict hook should include golangci-lint")
	}
	if !strings.Contains(content, "secret") {
		t.Error("strict hook should include secret scan")
	}
}

func TestRenderNodeLight(t *testing.T) {
	content := Render(ModeLight, "node")
	if !strings.Contains(content, "prettier") {
		t.Error("node light hook should include prettier")
	}
}

func TestRenderGeneric(t *testing.T) {
	content := Render(ModeLight, "generic")
	if !strings.Contains(content, "#!/bin/sh") {
		t.Error("hook should have shebang")
	}
}
