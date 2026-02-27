package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/TerraFaster/scaf/internal/structure"
	"github.com/TerraFaster/scaf/internal/ui"
)

var (
	structureDryRun bool
)

var structureCmd = &cobra.Command{
	Use:   "structure [template]",
	Short: "Create a project directory structure from a template",
	Long: `Create a directory structure from a named template.

Available templates:
  layered            Classic layered (api / service / repository)
  clean-architecture Clean Architecture rings
  hexagonal          Ports & Adapters
  ddd                Domain-Driven Design
  microservice       Microservice layout
  monolith           Monolithic application
  cli                CLI application structure
  minimal            Minimal structure (src / tests / docs)

Custom templates can be added to ~/.scaf/templates/.`,
	Example: `  scaf structure
  scaf structure clean-architecture
  scaf structure microservice --dry-run`,
	Args: cobra.MaximumNArgs(1),
	RunE: runStructure,
}

func init() {
	structureCmd.Flags().BoolVar(&structureDryRun, "dry-run", false, "Show what would be created without writing")
}

func runStructure(cmd *cobra.Command, args []string) error {
	available := structure.AvailableTemplates()
	sort.Strings(available)

	var templateName string

	if len(args) > 0 {
		templateName = args[0]
	} else {
		items := make([]ui.SelectItem, len(available))
		for i, t := range available {
			tmpl, _ := structure.Get(t)
			label := t
			if tmpl.Description != "" {
				label = fmt.Sprintf("%-22s — %s", t, tmpl.Description)
			}
			items[i] = ui.SelectItem{ID: t, Label: label}
		}
		chosen, err := ui.SelectOne("Select structure template:", items)
		if err != nil {
			return fmt.Errorf("selection cancelled: %w", err)
		}
		templateName = chosen
	}

	tmpl, ok := structure.Get(templateName)
	if !ok {
		return fmt.Errorf("unknown template %q\n\nAvailable: %v", templateName, available)
	}

	fmt.Printf("✔ Template: %s\n", tmpl.Name)
	if tmpl.Description != "" {
		fmt.Printf("  %s\n", tmpl.Description)
	}
	fmt.Println()

	if structureDryRun {
		fmt.Println("Dry run — would create:")
		return structure.Apply(tmpl, true)
	}

	if err := structure.Apply(tmpl, false); err != nil {
		return err
	}

	fmt.Printf("\n✔ Structure created: %s\n", tmpl.Name)
	return nil
}
