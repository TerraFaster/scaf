package cmd

import (
	"fmt"

	"github.com/TerraFaster/scaf/internal/fs"
	"github.com/TerraFaster/scaf/internal/readme"
	"github.com/TerraFaster/scaf/internal/ui"
	"github.com/spf13/cobra"
)

var (
	readmeName        string
	readmeDescription string
	readmeAuthor      string
	readmeLicense     string
	readmeRepo        string
	readmeBadges      bool
	readmeCI          bool
	readmeDocker      bool
	readmeForce       bool
	readmeDryRun      bool
	readmeYes         bool
)

var readmeCmd = &cobra.Command{
	Use:   "readme [template]",
	Short: "Generate a README.md file",
	Long: `Generate a README.md from a named template.

Available templates: minimal, cli, library, webapp, api, go, node, unity, custom

If no template is provided, an interactive picker is shown.
Parameters are auto-filled from the environment (project name from dir,
author from git config, license from LICENSE file).`,
	Example: `  scaf readme
  scaf readme go --name myapp --author "Username"
  scaf readme cli --dry-run
  scaf readme api --docker`,
	Args: cobra.MaximumNArgs(1),
	RunE: runReadme,
}

func init() {
	readmeCmd.Flags().StringVar(&readmeName, "name", "", "Project name (default: current directory name)")
	readmeCmd.Flags().StringVar(&readmeDescription, "description", "", "Project description")
	readmeCmd.Flags().StringVar(&readmeAuthor, "author", "", "Author name (default: git config user.name)")
	readmeCmd.Flags().StringVar(&readmeLicense, "license", "", "License (default: detected from LICENSE file)")
	readmeCmd.Flags().StringVar(&readmeRepo, "repo", "", "Repository URL or module path")
	readmeCmd.Flags().BoolVar(&readmeBadges, "badges", false, "Include CI/status badges")
	readmeCmd.Flags().BoolVar(&readmeCI, "ci", false, "Include CI section")
	readmeCmd.Flags().BoolVar(&readmeDocker, "docker", false, "Include Docker section")
	readmeCmd.Flags().BoolVar(&readmeForce, "force", false, "Overwrite existing README.md")
	readmeCmd.Flags().BoolVar(&readmeDryRun, "dry-run", false, "Print content without writing")
	readmeCmd.Flags().BoolVarP(&readmeYes, "yes", "y", false, "Skip confirmation prompts")
}

func runReadme(cmd *cobra.Command, args []string) error {
	params := readme.Params{
		Name:        readmeName,
		Description: readmeDescription,
		Author:      readmeAuthor,
		License:     readmeLicense,
		Repo:        readmeRepo,
		Badges:      readmeBadges,
		CI:          readmeCI,
		Docker:      readmeDocker,
	}

	// Auto-fill missing params
	params.AutoFill()

	// Determine template
	var templateName string
	if len(args) > 0 {
		templateName = args[0]
	} else if !readmeYes {
		available := readme.AvailableTemplates()
		items := make([]ui.SelectItem, len(available))
		for i, t := range available {
			items[i] = ui.SelectItem{ID: t, Label: t}
		}
		chosen, err := ui.SelectOne("Select README template:", items)
		if err != nil {
			return fmt.Errorf("selection cancelled: %w", err)
		}
		templateName = chosen
	} else {
		templateName = "minimal"
	}
	params.Template = templateName

	// Interactive fill for missing fields
	if !readmeYes {
		if params.Description == "" {
			desc, err := ui.Prompt("Project description", "")
			if err == nil && desc != "" {
				params.Description = desc
			}
		}
	}

	content, err := readme.Render(params)
	if err != nil {
		return fmt.Errorf("failed to render README: %w", err)
	}

	fmt.Printf("✔ Template: %s\n", templateName)
	fmt.Printf("✔ Name: %s\n", params.Name)

	if readmeDryRun {
		fmt.Println("\n--- README.md content ---")
		fmt.Println(content)
		return nil
	}

	writer := fs.NewWriter()
	opts := fs.WriteOptions{
		Force: readmeForce,
		Yes:   readmeYes,
	}

	if err := writer.Write("README.md", content, opts); err != nil {
		return err
	}

	fmt.Println("✔ File README.md created")
	return nil
}
