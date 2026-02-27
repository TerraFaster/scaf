// Package gitprofile manages switching git user.name/email per repo.
package gitprofile

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/TerraFaster/scaf/internal/config"
)

// EnsureGitRepo checks that the current directory is inside a git repository.
func EnsureGitRepo() error {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("not a git repository (or git is not installed)")
	}
	return nil
}

// Apply sets local git user.name and user.email for the current repo.
func Apply(profile config.GitProfile) error {
	if err := EnsureGitRepo(); err != nil {
		return err
	}
	if err := runGit("config", "user.name", profile.Name); err != nil {
		return fmt.Errorf("failed to set user.name: %w", err)
	}
	if err := runGit("config", "user.email", profile.Email); err != nil {
		return fmt.Errorf("failed to set user.email: %w", err)
	}
	return nil
}

// Current returns the current local git user.name and user.email.
func Current() (name, email string, err error) {
	if err := EnsureGitRepo(); err != nil {
		return "", "", err
	}
	name = gitConfigGet("user.name")
	email = gitConfigGet("user.email")
	return name, email, nil
}

// FindProfile returns the profile name that matches the current git config, if any.
func FindProfile(profiles map[string]config.GitProfile) (string, bool) {
	name := gitConfigGet("user.name")
	email := gitConfigGet("user.email")
	for k, p := range profiles {
		if p.Name == name && p.Email == email {
			return k, true
		}
	}
	return "", false
}

func gitConfigGet(key string) string {
	var out bytes.Buffer
	cmd := exec.Command("git", "config", key)
	cmd.Stdout = &out
	cmd.Stderr = nil
	_ = cmd.Run()
	return strings.TrimSpace(out.String())
}

func runGit(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
