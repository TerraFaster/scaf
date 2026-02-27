package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/TerraFaster/scaf/internal/cache"
	"github.com/TerraFaster/scaf/internal/config"
)

// Cfg holds the configuration loaded at startup for the current invocation.
var Cfg config.Config

var rootCmd = &cobra.Command{
	Use:   "scaf",
	Short: "Generate standard project files (LICENSE, .gitignore)",
	Long: `scaf — a fast CLI tool for generating standard project service files.

It creates LICENSE and .gitignore files using up-to-date templates
from GitHub and gitignore.io, with local caching for offline use.

Global config is stored at ~/.scaf/config.yaml and is created
automatically on first run. Run "scaf config" to view or edit it.`,
	// Load (or auto-create) the global config before any subcommand runs.
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
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

		// Print path and current contents.
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

	// No $EDITOR set — use platform default opener.
	var openCmd string
	switch runtime.GOOS {
	case "windows":
		// "start" is a shell built-in; use cmd.exe
		c := exec.Command("cmd", "/c", "start", "", path)
		return c.Run()
	case "darwin":
		openCmd = "open"
	default: // linux and others
		openCmd = "xdg-open"
	}

	return exec.Command(openCmd, path).Run()
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
		fmt.Println("🔄 Clearing template cache...")
		if err := c.Clear(); err != nil {
			return fmt.Errorf("failed to clear cache: %w", err)
		}
		fmt.Println("✔ Cache cleared. Templates will be re-downloaded on next use.")
		return nil
	},
}

// ── Execute ───────────────────────────────────────────────────────────────────

func Execute() {
	configCmd.Flags().BoolVarP(&configEditFlag, "edit", "e", false, "Open config in $EDITOR")

	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(licenseCmd)
	rootCmd.AddCommand(ignoreCmd)
	rootCmd.AddCommand(autoCmd)
	rootCmd.AddCommand(updateCmd)

	// alias: gitignore → ignore
	gitignoreAlias := *ignoreCmd
	gitignoreAlias.Use = "gitignore [templates...]"
	gitignoreAlias.Short = "Alias for the ignore command"
	rootCmd.AddCommand(&gitignoreAlias)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
