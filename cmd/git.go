package cmd

import (
	"fmt"
	"sort"

	"github.com/TerraFaster/scaf/internal/config"
	"github.com/TerraFaster/scaf/internal/gitprofile"
	"github.com/spf13/cobra"
)

var gitCmd = &cobra.Command{
	Use:   "git",
	Short: "Git-related helpers (profile switching, etc.)",
	Long:  `Git-related helpers for scaf. Use "scaf git profile" to manage git identities.`,
}

var gitProfileCmd = &cobra.Command{
	Use:   "profile [name]",
	Short: "Switch local git user.name/email to a named profile",
	Long: `Switch the local git user.name and user.email to a named profile.

Profiles are defined in ~/.scaf/config.yaml:

  profiles:
    work:
      name: Username
      email: username@company.com
    personal:
      name: user.name
      email: username@gmail.com

Usage:
  scaf git profile work         # apply profile
  scaf git profile list         # list all profiles
  scaf git profile show         # show current profile
  scaf git profile current      # alias for show`,
	Example: `  scaf git profile work
  scaf git profile list
  scaf git profile show`,
	Args: cobra.MaximumNArgs(1),
	RunE: runGitProfile,
}

var gitProfileListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured git profiles",
	RunE:  runGitProfileList,
}

var gitProfileShowCmd = &cobra.Command{
	Use:     "show",
	Short:   "Show current local git user.name/email",
	Aliases: []string{"current"},
	RunE:    runGitProfileShow,
}

func init() {
	gitCmd.AddCommand(gitProfileCmd)
	gitCmd.AddCommand(gitProfileListCmd)
	gitCmd.AddCommand(gitProfileShowCmd)
}

func runGitProfile(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return runGitProfileList(cmd, args)
	}

	name := args[0]
	switch name {
	case "list":
		return runGitProfileList(cmd, []string{})
	case "show", "current":
		return runGitProfileShow(cmd, []string{})
	}

	if err := gitprofile.EnsureGitRepo(); err != nil {
		return err
	}

	profile, ok := Cfg.Profiles[name]
	if !ok {
		fmt.Printf("✖ Profile %q not found.\n\n", name)
		fmt.Println("Available profiles:")
		return runGitProfileList(cmd, []string{})
	}

	if err := gitprofile.Apply(profile); err != nil {
		return err
	}

	fmt.Printf("✔ Git profile applied: %s\n", name)
	fmt.Printf("  user.name  = %s\n", profile.Name)
	fmt.Printf("  user.email = %s\n", profile.Email)
	return nil
}

func runGitProfileList(cmd *cobra.Command, args []string) error {
	if len(Cfg.Profiles) == 0 {
		fmt.Println("No git profiles configured.")
		fmt.Println()
		fmt.Println("Add profiles to ~/.scaf/config.yaml:")
		fmt.Println()
		fmt.Println("  profiles:")
		fmt.Println("    work:")
		fmt.Println("      name: Username")
		fmt.Println("      email: username@company.com")
		fmt.Println("    personal:")
		fmt.Println("      name: user.name")
		fmt.Println("      email: username@gmail.com")
		return nil
	}

	keys := make([]string, 0, len(Cfg.Profiles))
	for k := range Cfg.Profiles {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Try to detect which profile is active
	activeName, _ := gitprofile.FindProfile(Cfg.Profiles)

	fmt.Println("Configured git profiles:")
	fmt.Println()
	for _, k := range keys {
		p := Cfg.Profiles[k]
		marker := "  "
		if k == activeName {
			marker = "* "
		}
		fmt.Printf("%s%-20s  %s <%s>\n", marker, k, p.Name, p.Email)
	}

	if activeName != "" {
		fmt.Printf("\n  * = active profile in current repo\n")
	}
	return nil
}

func runGitProfileShow(cmd *cobra.Command, args []string) error {
	name, email, err := gitprofile.Current()
	if err != nil {
		return err
	}

	fmt.Println("Current local git identity:")
	fmt.Printf("  user.name  = %s\n", orEmpty(name, "(not set)"))
	fmt.Printf("  user.email = %s\n", orEmpty(email, "(not set)"))

	// Show which named profile matches, if any
	if profileName, ok := gitprofile.FindProfile(Cfg.Profiles); ok {
		fmt.Printf("  profile    = %s\n", profileName)
	}
	return nil
}

func orEmpty(s, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}

// AddGitProfile adds a profile to config and saves it.
func AddGitProfile(name string, profile config.GitProfile) error {
	if Cfg.Profiles == nil {
		Cfg.Profiles = make(map[string]config.GitProfile)
	}
	Cfg.Profiles[name] = profile
	return config.Save(Cfg)
}
