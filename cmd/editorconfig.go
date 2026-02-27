package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/TerraFaster/scaf/internal/editorconfig"
	"github.com/TerraFaster/scaf/internal/fs"
	"github.com/TerraFaster/scaf/internal/ui"
)

var (
	editorconfigVariant    string
	editorconfigForce      bool
	editorconfigDryRun     bool
	editorconfigYes        bool
	editorconfigIndent     string
	editorconfigIndentSize string
	editorconfigEOL        string
	editorconfigCharset    string
)

var editorconfigCmd = &cobra.Command{
	Use:   "editorconfig [variant]",
	Short: "Generate a .editorconfig file",
	Long: `Generate a .editorconfig file for the current project.

Available variants: default, strict, go, node, unity, dotnet, custom

If no variant is provided, the project type is auto-detected from files
in the current directory (go.mod → go, package.json → node, etc.).`,
	Example: `  scaf editorconfig
  scaf editorconfig go
  scaf editorconfig node --indent-size 2
  scaf editorconfig --dry-run`,
	Args: cobra.MaximumNArgs(1),
	RunE: runEditorconfig,
}

func init() {
	editorconfigCmd.Flags().StringVar(&editorconfigVariant, "variant", "", "Variant (default|strict|go|node|unity|dotnet|custom)")
	editorconfigCmd.Flags().BoolVar(&editorconfigForce, "force", false, "Overwrite existing .editorconfig")
	editorconfigCmd.Flags().BoolVar(&editorconfigDryRun, "dry-run", false, "Print content without writing")
	editorconfigCmd.Flags().BoolVarP(&editorconfigYes, "yes", "y", false, "Skip confirmation prompts")
	editorconfigCmd.Flags().StringVar(&editorconfigIndent, "indent-style", "", "Override indent_style (space|tab)")
	editorconfigCmd.Flags().StringVar(&editorconfigIndentSize, "indent-size", "", "Override indent_size (2|4|...)")
	editorconfigCmd.Flags().StringVar(&editorconfigEOL, "eol", "", "Override end_of_line (lf|crlf|cr)")
	editorconfigCmd.Flags().StringVar(&editorconfigCharset, "charset", "", "Override charset (utf-8|utf-16be|...)")
}

func runEditorconfig(cmd *cobra.Command, args []string) error {
	variant := editorconfigVariant

	// From positional arg
	if variant == "" && len(args) > 0 {
		variant = args[0]
	}

	// Auto-detect
	if variant == "" {
		variant = editorconfig.AutoDetectVariant()
		fmt.Printf("✔ Auto-detected variant: %s\n", variant)
	}

	// Interactive fallback
	if variant == "" && !editorconfigYes {
		available := editorconfig.AvailableVariants()
		items := make([]ui.SelectItem, len(available))
		for i, v := range available {
			items[i] = ui.SelectItem{ID: v, Label: v}
		}
		chosen, err := ui.SelectOne("Select .editorconfig variant:", items)
		if err != nil {
			return fmt.Errorf("selection cancelled: %w", err)
		}
		variant = chosen
	}

	params := editorconfig.Params{
		Variant:               variant,
		IndentStyle:           editorconfigIndent,
		IndentSize:            editorconfigIndentSize,
		EndOfLine:             editorconfigEOL,
		Charset:               editorconfigCharset,
		InsertFinalNewline:    true,
		TrimTrailingWhitespace: true,
	}

	content := editorconfig.Render(params)

	fmt.Printf("✔ Variant: %s\n", variant)

	if editorconfigDryRun {
		fmt.Println("\n--- .editorconfig content ---")
		fmt.Println(content)
		return nil
	}

	writer := fs.NewWriter()
	opts := fs.WriteOptions{
		Force: editorconfigForce,
		Yes:   editorconfigYes,
	}

	if err := writer.Write(".editorconfig", content, opts); err != nil {
		return err
	}

	fmt.Println("✔ File .editorconfig created")
	return nil
}
