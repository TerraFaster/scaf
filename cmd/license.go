package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/TerraFaster/scaf/internal/fs"
	"github.com/TerraFaster/scaf/internal/templates"
	"github.com/TerraFaster/scaf/internal/ui"
	"github.com/spf13/cobra"
)

var (
	licenseAuthor string
	licenseYear   string
	licenseForce  bool
	licenseBackup bool
	licenseDryRun bool
	licenseYes    bool
)

var licenseCmd = &cobra.Command{
	Use:   "license [name]",
	Short: "Generate a LICENSE file",
	Long: `Generate a LICENSE file in the current directory.
If no license name is provided, an interactive list will be shown.`,
	RunE: runLicense,
}

func init() {
	licenseCmd.Flags().StringVar(&licenseAuthor, "author", "", "Author name (default: git config user.name)")
	licenseCmd.Flags().StringVar(&licenseYear, "year", "", "Year for the license (default: current year)")
	licenseCmd.Flags().BoolVar(&licenseForce, "force", false, "Overwrite existing LICENSE file")
	licenseCmd.Flags().BoolVar(&licenseBackup, "backup", false, "Create LICENSE.bak before overwriting")
	licenseCmd.Flags().BoolVar(&licenseDryRun, "dry-run", false, "Print content to stdout without writing file")
	licenseCmd.Flags().BoolVarP(&licenseYes, "yes", "y", false, "Skip confirmation prompts")
}

func runLicense(cmd *cobra.Command, args []string) error {
	provider, err := templates.NewLicenseProviderWithTTL(Cfg.CacheTTLHours)
	if err != nil {
		return fmt.Errorf("failed to initialize license provider: %w", err)
	}

	var selectedKey string

	if len(args) == 0 {
		// Use default_license from config if set
		if Cfg.DefaultLicense != "" {
			selectedKey = strings.ToLower(Cfg.DefaultLicense)
			fmt.Printf("ℹ Using default license from scaf.yaml: %s\n", selectedKey)
		} else {
			// Interactive mode
			licenses, err := provider.List()
			if err != nil {
				return fmt.Errorf("failed to load licenses: %w", err)
			}
			items := make([]ui.SelectItem, len(licenses))
			for i, l := range licenses {
				items[i] = ui.SelectItem{ID: l.Key, Label: l.Name}
			}
			chosen, err := ui.SelectOne("Select a license:", items)
			if err != nil {
				return fmt.Errorf("selection cancelled: %w", err)
			}
			selectedKey = chosen
		}
	} else {
		selectedKey = strings.ToLower(args[0])
	}

	// Fetch license body
	body, err := provider.Get(selectedKey)
	if err != nil {
		// fuzzy search fallback
		licenses, listErr := provider.List()
		if listErr == nil {
			suggestions := provider.FuzzySearch(selectedKey, licenses)
			queryLabel := selectedKey
			if len(args) > 0 {
				queryLabel = args[0]
			}
			fmt.Fprintf(os.Stderr, "✖ License %q not found\n", queryLabel)
			if len(suggestions) > 0 {
				fmt.Fprintln(os.Stderr, "Did you mean:")
				for _, s := range suggestions {
					fmt.Fprintf(os.Stderr, "  - %s\n", s.Key)
				}
			}
		}
		return fmt.Errorf("license not found: %s", selectedKey)
	}

	// Resolve author: flag > config > git config
	author := licenseAuthor
	if author == "" && Cfg.DefaultAuthor != "" {
		author = Cfg.DefaultAuthor
	}
	if author == "" {
		author = resolveGitAuthor()
	}

	// Resolve year: flag > config > current year
	year := licenseYear
	if year == "" && Cfg.DefaultYear != "" {
		year = Cfg.DefaultYear
	}
	if year == "" {
		year = strconv.Itoa(time.Now().Year())
	}

	// Replace placeholders
	content := body
	content = strings.ReplaceAll(content, "[year]", year)
	content = strings.ReplaceAll(content, "[fullname]", author)
	content = strings.ReplaceAll(content, "[name of copyright owner]", author)
	content = strings.ReplaceAll(content, "<year>", year)
	content = strings.ReplaceAll(content, "<name of author>", author)
	content = strings.ReplaceAll(content, "YEAR", year)
	content = strings.ReplaceAll(content, "AUTHOR", author)

	fmt.Printf("✔ License %s selected\n", strings.ToUpper(selectedKey))
	fmt.Printf("✔ Author: %s\n", author)
	fmt.Printf("✔ Year: %s\n", year)

	if licenseDryRun {
		fmt.Println("\n--- LICENSE content ---")
		fmt.Println(content)
		return nil
	}

	writer := fs.NewWriter()
	opts := fs.WriteOptions{
		Force:  licenseForce,
		Backup: licenseBackup,
		Yes:    licenseYes,
	}

	if err := writer.Write("LICENSE", content, opts); err != nil {
		return err
	}

	fmt.Println("✔ File LICENSE created")
	return nil
}

func resolveGitAuthor() string {
	// Try git config
	out, err := runCommand("git", "config", "user.name")
	if err == nil && strings.TrimSpace(out) != "" {
		return strings.TrimSpace(out)
	}
	return "Author"
}
