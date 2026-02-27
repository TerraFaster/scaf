package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/TerraFaster/scaf/internal/cache"
	"github.com/TerraFaster/scaf/internal/config"
	"github.com/spf13/cobra"
)

// Cfg holds the configuration loaded at startup for the current invocation.
var Cfg config.Config

// Global flags
var (
	globalVerbose bool
	globalQuiet   bool
	globalConfig  string
)

var rootCmd = &cobra.Command{
	Use:   "scaf",
	Short: "A local bootstrap tool for initializing and standardizing projects",
	Long: `scaf — a fast CLI tool for bootstrapping and standardizing projects.

Config: ~/.scaf/config.yaml`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var cfgPath string
		if globalConfig != "" {
			cfgPath = globalConfig
		}
		_ = cfgPath

		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "⚠  Warning: could not load config: %v\n", err)
			cfg = config.Defaults()
		}
		Cfg = cfg
		return nil
	},
}

// ── scaf config ───────────────────────────────────────────────────────────────

var configEditFlag bool

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show or open the global config file (~/.scaf/config.yaml)",
	Long: `Show the path to the global config file and print its contents,
or open it directly in your default editor.

The config file is created automatically on the first run of scaf.`,
	Example: `  scaf config           # print path + contents
  scaf config --edit    # open in $EDITOR / system default`,
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := config.Path()
		if err != nil {
			return err
		}

		if configEditFlag {
			return openEditor(path)
		}

		fmt.Printf("Config file: %s\n\n", path)

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("cannot read config: %w", err)
		}
		fmt.Println(string(data))
		return nil
	},
}

// openEditor opens path in $EDITOR, falling back to the OS default opener.
func openEditor(path string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}

	if editor != "" {
		c := exec.Command(editor, path)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		return c.Run()
	}

	switch runtime.GOOS {
	case "windows":
		c := exec.Command("cmd", "/c", "start", "", path)
		return c.Run()
	case "darwin":
		return exec.Command("open", path).Run()
	default:
		return exec.Command("xdg-open", path).Run()
	}
}

// ── scaf update ───────────────────────────────────────────────────────────────

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Force-refresh all cached templates",
	Long:  `Delete all locally cached templates so they are re-downloaded on next use.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := cache.NewWithOptions(Cfg.CacheTTLHours)
		if err != nil {
			return err
		}
		fmt.Println("Clearing template cache...")
		if err := c.Clear(); err != nil {
			return fmt.Errorf("failed to clear cache: %w", err)
		}
		fmt.Println("✔ Cache cleared.")
		return nil
	},
}

// ── scaf version ──────────────────────────────────────────────────────────────

var Version = "2.0.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print scaf version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("scaf version %s\n", Version)
	},
}

// ── Execute ───────────────────────────────────────────────────────────────────

func Execute() {
	// Global persistent flags
	rootCmd.PersistentFlags().BoolVarP(&globalVerbose, "verbose", "v", false, "Verbose output")
	rootCmd.PersistentFlags().BoolVarP(&globalQuiet, "quiet", "q", false, "Minimal output")
	rootCmd.PersistentFlags().StringVar(&globalConfig, "config", "", "Custom config file path")

	configCmd.Flags().BoolVarP(&configEditFlag, "edit", "e", false, "Open config in $EDITOR")

	// Register all commands
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(updateCmd)

	rootCmd.AddCommand(licenseCmd)
	rootCmd.AddCommand(ignoreCmd)
	rootCmd.AddCommand(autoCmd)

	rootCmd.AddCommand(gitCmd)
	rootCmd.AddCommand(readmeCmd)
	rootCmd.AddCommand(editorconfigCmd)
	rootCmd.AddCommand(dockerignoreCmd)

	rootCmd.AddCommand(hooksCmd)
	rootCmd.AddCommand(initCmd)

	rootCmd.AddCommand(structureCmd)

	// Alias: gitignore → ignore
	gitignoreAlias := *ignoreCmd
	gitignoreAlias.Use = "gitignore [templates...]"
	gitignoreAlias.Short = "Alias for the ignore command"
	rootCmd.AddCommand(&gitignoreAlias)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
