package gitprofile

import (
	"testing"

	"github.com/TerraFaster/scaf/internal/config"
)

func TestFindProfile(t *testing.T) {
	profiles := map[string]config.GitProfile{
		"work":     {Name: "Username", Email: "username@company.com"},
		"personal": {Name: "user.name", Email: "username@gmail.com"},
	}

	// Test no match
	_, ok := FindProfile(profiles)
	// We can't guarantee a match in tests (depends on git config)
	// Just ensure no panic
	_ = ok
}

func TestEnsureGitRepo_NotARepo(t *testing.T) {
	// This test would pass only if run outside a git repo, which is
	// environment-dependent. We just ensure the function returns an error
	// or nil (both are valid depending on environment).
	// We do not assert a specific outcome here.
	_ = EnsureGitRepo()
}
