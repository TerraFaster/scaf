package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/TerraFaster/scaf/internal/fs"
	"github.com/TerraFaster/scaf/internal/templates"
)

var (
	autoForce  bool
	autoBackup bool
	autoDryRun bool
	autoYes    bool
)

var autoCmd = &cobra.Command{
	Use:   "auto",
	Short: "Auto-detect project type and generate .gitignore",
	Long: `Scan the current directory for well-known project markers and
automatically generate a combined .gitignore tailored to your stack.

Detected markers:
  package.json   → node
  Cargo.toml     → rust
  go.mod         → go
  requirements.txt → python
  *.csproj       → dotnet
  Assets/        → unity (directory)`,
	Example: `  scaf auto
  scaf auto --dry-run
  scaf auto --force`,
	RunE: runAuto,
}

func init() {
	autoCmd.Flags().BoolVar(&autoForce, "force", false, "Overwrite existing .gitignore file")
	autoCmd.Flags().BoolVar(&autoBackup, "backup", false, "Create .gitignore.bak before overwriting")
	autoCmd.Flags().BoolVar(&autoDryRun, "dry-run", false, "Print content to stdout without writing file")
	autoCmd.Flags().BoolVarP(&autoYes, "yes", "y", false, "Skip confirmation prompts")
}

func runAuto(cmd *cobra.Command, args []string) error {
	detected, err := autoDetectProject()
	if err != nil {
		return fmt.Errorf("detection failed: %w", err)
	}

	if len(detected) == 0 {
		fmt.Println("✖ Could not detect any known project type in the current directory.")
		fmt.Println()
		fmt.Println("  scaf looks for these markers:")
		fmt.Println("    package.json     → node")
		fmt.Println("    Cargo.toml       → rust")
		fmt.Println("    go.mod           → go")
		fmt.Println("    requirements.txt → python")
		fmt.Println("    *.csproj         → dotnet")
		fmt.Println("    Assets/          → unity")
		fmt.Println()
		fmt.Println("  You can still run: scaf ignore <template>")
		return fmt.Errorf("no project type detected")
	}

	fmt.Printf("✔ Detected: %s\n", strings.Join(detected, ", "))

	provider, err := templates.NewIgnoreProviderWithTTL(Cfg.CacheTTLHours)
	if err != nil {
		return fmt.Errorf("failed to initialize gitignore provider: %w", err)
	}

	content, err := provider.Get(detected)
	if err != nil {
		return fmt.Errorf("failed to fetch gitignore templates: %w", err)
	}

	if autoDryRun {
		fmt.Println("\n--- .gitignore content ---")
		fmt.Println(content)
		return nil
	}

	writer := fs.NewWriter()
	opts := fs.WriteOptions{
		Force:  autoForce,
		Backup: autoBackup,
		Yes:    autoYes,
	}

	if err := writer.Write(".gitignore", content, opts); err != nil {
		return err
	}

	fmt.Println("✔ File .gitignore created")
	return nil
}
