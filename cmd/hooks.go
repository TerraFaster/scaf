package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/TerraFaster/scaf/internal/hooks"
)

var (
	hooksMode   string
	hooksForce  bool
	hooksDryRun bool
	hooksStack  string
)

var hooksCmd = &cobra.Command{
	Use:   "hooks",
	Short: "Install git pre-commit hooks",
	Long: `Install git pre-commit hooks for the current repository.

Modes:
  --light    Basic checks (whitespace, gofmt/prettier)
  --strict   Full checks (lint, secret scan, commit message)
  --custom   Placeholder for custom hooks

The hook is installed at .git/hooks/pre-commit and made executable.
The stack is auto-detected but can be overridden with --stack.`,
	Example: `  scaf hooks
  scaf hooks --light
  scaf hooks --strict
  scaf hooks --stack go --force`,
	RunE: runHooks,
}

func init() {
	hooksCmd.Flags().StringVar(&hooksMode, "mode", "light", "Hook mode: light|strict|custom")
	hooksCmd.Flags().BoolVar(&hooksForce, "force", false, "Overwrite existing pre-commit hook")
	hooksCmd.Flags().BoolVar(&hooksDryRun, "dry-run", false, "Print hook content without writing")
	hooksCmd.Flags().StringVar(&hooksStack, "stack", "", "Stack: go|node|generic (auto-detected if empty)")

	// Convenience flags
	hooksCmd.Flags().Bool("light", false, "Use light mode (shortcut for --mode light)")
	hooksCmd.Flags().Bool("strict", false, "Use strict mode (shortcut for --mode strict)")
}

func runHooks(cmd *cobra.Command, args []string) error {
	// Resolve mode from shorthand flags
	mode := hooks.Mode(hooksMode)
	if v, _ := cmd.Flags().GetBool("strict"); v {
		mode = hooks.ModeStrict
	}
	if v, _ := cmd.Flags().GetBool("light"); v {
		mode = hooks.ModeLight
	}

	// Auto-detect stack
	stack := hooksStack
	if stack == "" {
		stack = autoDetectStack()
	}

	content := hooks.Render(mode, stack)
	hookPath := hooks.HookPath()

	fmt.Printf("✔ Mode: %s\n", mode)
	fmt.Printf("✔ Stack: %s\n", stack)

	if hooksDryRun {
		fmt.Printf("\n--- %s ---\n", hookPath)
		fmt.Println(content)
		return nil
	}

	if err := hooks.EnsureGitRepo(); err != nil {
		return err
	}

	// Check existing
	if _, err := os.Stat(hookPath); err == nil && !hooksForce {
		return fmt.Errorf("hook already exists at %s (use --force to overwrite)", hookPath)
	}

	if err := hooks.Install(content, hookPath); err != nil {
		return err
	}

	fmt.Printf("✔ Hook installed: %s\n", hookPath)
	return nil
}

// autoDetectStack detects the project stack for hooks.
func autoDetectStack() string {
	detected, _ := autoDetectProject()
	for _, d := range detected {
		switch d {
		case "go", "node":
			return d
		}
	}
	return "generic"
}
