package dockerignore

import (
	"strings"
	"testing"
)

func TestRenderGo(t *testing.T) {
	content := Render("go")
	if !strings.Contains(content, ".git") {
		t.Error("go .dockerignore should contain .git")
	}
	if !strings.Contains(content, "/bin/") {
		t.Error("go .dockerignore should contain /bin/")
	}
}

func TestRenderNode(t *testing.T) {
	content := Render("node")
	if !strings.Contains(content, "node_modules/") {
		t.Error("node .dockerignore should contain node_modules/")
	}
}

func TestRenderUnknownFallsToGeneric(t *testing.T) {
	content := Render("unknown-stack")
	if content == "" {
		t.Error("unknown stack should fall back to generic")
	}
}

func TestAllStacks(t *testing.T) {
	for name := range stacks {
		content := Render(name)
		if content == "" {
			t.Errorf("stack %q produced empty content", name)
		}
	}
}

func TestAvailableStacks(t *testing.T) {
	stacks := AvailableStacks()
	if len(stacks) == 0 {
		t.Error("should have available stacks")
	}
}
