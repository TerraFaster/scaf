package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/TerraFaster/scaf/internal/fs"
	"github.com/TerraFaster/scaf/internal/templates"
	"github.com/TerraFaster/scaf/internal/ui"
)

var (
	ignoreForce   bool
	ignoreBackup  bool
	ignoreDryRun  bool
	ignoreYes     bool
	ignoreAuto    bool
)

var popularIgnoreTemplates = []string{
	"python", "node", "unity", "dotnet", "rust", "go",
}

var ignoreCmd = &cobra.Command{
	Use:   "ignore [templates...]",
	Short: "Generate a .gitignore file",
	Long: `Generate a .gitignore file in the current directory.
Multiple templates can be specified (comma-separated or as multiple args).
If no template is provided, an interactive multi-select list will be shown.`,
	RunE: runIgnore,
}

func init() {
	ignoreCmd.Flags().BoolVar(&ignoreForce, "force", false, "Overwrite existing .gitignore file")
	ignoreCmd.Flags().BoolVar(&ignoreBackup, "backup", false, "Create .gitignore.bak before overwriting")
	ignoreCmd.Flags().BoolVar(&ignoreDryRun, "dry-run", false, "Print content to stdout without writing file")
	ignoreCmd.Flags().BoolVarP(&ignoreYes, "yes", "y", false, "Skip confirmation prompts")
	ignoreCmd.Flags().BoolVar(&ignoreAuto, "auto", false, "Auto-detect project type and generate .gitignore")
}

func runIgnore(cmd *cobra.Command, args []string) error {
	provider, err := templates.NewIgnoreProviderWithTTL(Cfg.CacheTTLHours)
	if err != nil {
		return fmt.Errorf("failed to initialize gitignore provider: %w", err)
	}

	var selectedTemplates []string

	if ignoreAuto {
		detected, err := autoDetectProject()
		if err != nil {
			return err
		}
		if len(detected) == 0 {
			return fmt.Errorf("could not auto-detect project type")
		}
		fmt.Printf("✔ Detected project types: %s\n", strings.Join(detected, ", "))
		selectedTemplates = detected
	} else if len(args) > 0 {
		// Parse comma-separated or space-separated
		for _, arg := range args {
			parts := strings.Split(arg, ",")
			for _, p := range parts {
				p = strings.TrimSpace(p)
				if p != "" {
					selectedTemplates = append(selectedTemplates, strings.ToLower(p))
				}
			}
		}
	} else if len(Cfg.DefaultIgnore) > 0 {
		// Use defaults from scaf.yaml
		fmt.Printf("\u2139 Using default ignore templates from scaf.yaml: %s\n",
			strings.Join(Cfg.DefaultIgnore, ", "))
		selectedTemplates = Cfg.DefaultIgnore
	} else {
		// Interactive multi-select — show popular list
		if _, err := provider.List(); err != nil {
			return fmt.Errorf("failed to load templates: %w", err)
		}
		popularItems := make([]ui.SelectItem, 0, len(popularIgnoreTemplates))
		for _, name := range popularIgnoreTemplates {
			popularItems = append(popularItems, ui.SelectItem{ID: name, Label: name})
		}
		chosen, err := ui.MultiSelect("Select gitignore templates (space to select, enter to confirm):", popularItems)
		if err != nil {
			return fmt.Errorf("selection cancelled: %w", err)
		}
		selectedTemplates = chosen
	}

	if len(selectedTemplates) == 0 {
		return fmt.Errorf("no templates selected")
	}

	// Fetch combined gitignore
	content, err := provider.Get(selectedTemplates)
	if err != nil {
		return fmt.Errorf("failed to fetch gitignore templates: %w", err)
	}

	fmt.Printf("✔ Templates selected: %s\n", strings.Join(selectedTemplates, ", "))

	if ignoreDryRun {
		fmt.Println("\n--- .gitignore content ---")
		fmt.Println(content)
		return nil
	}

	writer := fs.NewWriter()
	opts := fs.WriteOptions{
		Force:  ignoreForce,
		Backup: ignoreBackup,
		Yes:    ignoreYes,
	}

	if err := writer.Write(".gitignore", content, opts); err != nil {
		return err
	}

	fmt.Println("✔ File .gitignore created")
	return nil
}
