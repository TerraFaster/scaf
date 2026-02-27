package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/TerraFaster/scaf/internal/dockerignore"
	"github.com/TerraFaster/scaf/internal/fs"
	"github.com/TerraFaster/scaf/internal/ui"
)

var (
	dockerignoreStack   string
	dockerignoreForce   bool
	dockerignoreDryRun  bool
	dockerignoreYes     bool
)

var dockerignoreCmd = &cobra.Command{
	Use:   "dockerignore [stack]",
	Short: "Generate a .dockerignore file",
	Long: `Generate a .dockerignore file for the current project.

Available stacks: go, node, python, dotnet, unity, generic

If no stack is provided, the project type is auto-detected from files
in the current directory.`,
	Example: `  scaf dockerignore
  scaf dockerignore go
  scaf dockerignore node --force
  scaf dockerignore --dry-run`,
	Args: cobra.MaximumNArgs(1),
	RunE: runDockerignore,
}

func init() {
	dockerignoreCmd.Flags().StringVar(&dockerignoreStack, "stack", "", "Stack (go|node|python|dotnet|unity|generic)")
	dockerignoreCmd.Flags().BoolVar(&dockerignoreForce, "force", false, "Overwrite existing .dockerignore")
	dockerignoreCmd.Flags().BoolVar(&dockerignoreDryRun, "dry-run", false, "Print content without writing")
	dockerignoreCmd.Flags().BoolVarP(&dockerignoreYes, "yes", "y", false, "Skip confirmation prompts")
}

func runDockerignore(cmd *cobra.Command, args []string) error {
	stack := dockerignoreStack

	if stack == "" && len(args) > 0 {
		stack = args[0]
	}

	if stack == "" {
		stack = dockerignore.AutoDetectStack()
		fmt.Printf("✔ Auto-detected stack: %s\n", stack)
	}

	// Interactive fallback if auto-detect returned generic
	if stack == "generic" && !dockerignoreYes {
		available := dockerignore.AvailableStacks()
		items := make([]ui.SelectItem, len(available))
		for i, s := range available {
			items[i] = ui.SelectItem{ID: s, Label: s}
		}
		chosen, err := ui.SelectOne("Select project stack for .dockerignore:", items)
		if err == nil {
			stack = chosen
		}
	}

	content := dockerignore.Render(stack)

	fmt.Printf("✔ Stack: %s\n", stack)

	if dockerignoreDryRun {
		fmt.Println("\n--- .dockerignore content ---")
		fmt.Println(content)
		return nil
	}

	writer := fs.NewWriter()
	opts := fs.WriteOptions{
		Force: dockerignoreForce,
		Yes:   dockerignoreYes,
	}

	if err := writer.Write(".dockerignore", content, opts); err != nil {
		return err
	}

	fmt.Println("✔ File .dockerignore created")
	return nil
}
