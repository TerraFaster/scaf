package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/TerraFaster/scaf/internal/initproject"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new project (go, node)",
	Long:  `Initialize a new project with standard structure, files, and configuration.`,
}

// ── scaf init go ──────────────────────────────────────────────────────────────

var (
	initGoModeCLI     bool
	initGoModeService bool
	initGoModeLibrary bool
	initGoDocker      bool
	initGoDryRun      bool
	initGoForce       bool
)

var initGoCmd = &cobra.Command{
	Use:   "go <module-name>",
	Short: "Initialize a new Go project",
	Long: `Initialize a new Go project with standard directory structure.

Creates:
  cmd/<name>/main.go  internal/  pkg/  configs/  scripts/
  Makefile  README.md  LICENSE  .gitignore  .editorconfig

Modes:
  --cli       CLI application
  --service   HTTP/gRPC service
  --library   Reusable library`,
	Example: `  scaf init go github.com/user/myapp
  scaf init go github.com/user/myapp --cli
  scaf init go github.com/user/myapp --docker`,
	Args: cobra.ExactArgs(1),
	RunE: runInitGo,
}

func init() {
	initGoCmd.Flags().BoolVar(&initGoModeCLI, "cli", false, "CLI application mode")
	initGoCmd.Flags().BoolVar(&initGoModeService, "service", false, "Service/API mode")
	initGoCmd.Flags().BoolVar(&initGoModeLibrary, "library", false, "Library mode")
	initGoCmd.Flags().BoolVar(&initGoDocker, "docker", false, "Include Dockerfile")
	initGoCmd.Flags().BoolVar(&initGoDryRun, "dry-run", false, "Show what would be created")
	initGoCmd.Flags().BoolVar(&initGoForce, "force", false, "Overwrite existing files")

	initNodeCmd.Flags().BoolVar(&initNodeTS, "ts", false, "TypeScript project")
	initNodeCmd.Flags().BoolVar(&initNodeAPI, "api", false, "API project mode")
	initNodeCmd.Flags().BoolVar(&initNodeCLI, "cli", false, "CLI application mode")
	initNodeCmd.Flags().BoolVar(&initNodeLib, "lib", false, "Library mode")
	initNodeCmd.Flags().BoolVar(&initNodeDocker, "docker", false, "Include Dockerfile")
	initNodeCmd.Flags().BoolVar(&initNodeDryRun, "dry-run", false, "Show what would be created")
	initNodeCmd.Flags().BoolVar(&initNodeForce, "force", false, "Overwrite existing files")

	initCmd.AddCommand(initGoCmd)
	initCmd.AddCommand(initNodeCmd)
}

func runInitGo(cmd *cobra.Command, args []string) error {
	moduleName := args[0]

	mode := initproject.GoMode("default")
	if initGoModeCLI {
		mode = initproject.GoModeCLI
	} else if initGoModeService {
		mode = initproject.GoModeService
	} else if initGoModeLibrary {
		mode = initproject.GoModeLibrary
	}

	opts := initproject.GoOptions{
		ModuleName: moduleName,
		Mode:       mode,
		Docker:     initGoDocker,
		DryRun:     initGoDryRun,
	}

	files := initproject.GoFiles(opts)

	if initGoDryRun {
		fmt.Printf("Would initialize Go project: %s\n\n", moduleName)
		for _, f := range files {
			if f.Content == "" {
				fmt.Printf("  CREATE dir  %s\n", filepath.Dir(f.Path)+"/")
			} else {
				fmt.Printf("  CREATE file %s\n", f.Path)
			}
		}
		fmt.Printf("  RUN         go mod init %s\n", moduleName)
		return nil
	}

	fmt.Printf("Initializing Go project: %s\n\n", moduleName)

	cwd, _ := os.Getwd()
	for _, f := range files {
		full := filepath.Join(cwd, filepath.FromSlash(f.Path))
		dir := filepath.Dir(full)

		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("cannot create dir %s: %w", dir, err)
		}

		if _, err := os.Stat(full); err == nil && !initGoForce {
			fmt.Printf("  SKIP   %s (exists)\n", f.Path)
			continue
		}

		if err := os.WriteFile(full, []byte(f.Content), 0644); err != nil {
			return fmt.Errorf("cannot write %s: %w", f.Path, err)
		}
		fmt.Printf("  CREATE %s\n", f.Path)
	}

	// Run go mod init
	fmt.Printf("\n  RUN    go mod init %s\n", moduleName)
	if err := initproject.RunGoModInit(moduleName); err != nil {
		fmt.Printf("  ⚠  go mod init failed: %v\n", err)
		fmt.Printf("  Run manually: go mod init %s\n", moduleName)
	}

	fmt.Printf("\n✔ Go project initialized: %s\n", moduleName)
	fmt.Println("\nNext steps:")
	fmt.Println("  make build    # build the project")
	fmt.Println("  make test     # run tests")
	return nil
}

// ── scaf init node ────────────────────────────────────────────────────────────

var (
	initNodeTS      bool
	initNodeAPI     bool
	initNodeCLI     bool
	initNodeLib     bool
	initNodeDocker  bool
	initNodeDryRun  bool
	initNodeForce   bool
)

var initNodeCmd = &cobra.Command{
	Use:   "node <project-name>",
	Short: "Initialize a new Node.js project",
	Long: `Initialize a new Node.js project with standard structure.

Creates:
  src/  tests/  configs/
  .gitignore  .editorconfig  README.md  .eslintrc.json  .prettierrc

Modes:
  --ts    TypeScript project (adds tsconfig.json)
  --api   REST API project
  --cli   CLI application
  --lib   Library`,
	Example: `  scaf init node myapp
  scaf init node myapp --ts
  scaf init node myapp --api --docker`,
	Args: cobra.ExactArgs(1),
	RunE: runInitNode,
}


func runInitNode(cmd *cobra.Command, args []string) error {
	projectName := args[0]

	mode := initproject.NodeMode("default")
	if initNodeTS {
		mode = initproject.NodeModeTS
	} else if initNodeAPI {
		mode = initproject.NodeModeAPI
	} else if initNodeCLI {
		mode = initproject.NodeModeCLI
	} else if initNodeLib {
		mode = initproject.NodeModeLib
	}

	opts := initproject.NodeOptions{
		ProjectName: projectName,
		Mode:        mode,
		Docker:      initNodeDocker,
		DryRun:      initNodeDryRun,
	}

	files := initproject.NodeFiles(opts)

	if initNodeDryRun {
		fmt.Printf("Would initialize Node.js project: %s\n\n", projectName)
		for _, f := range files {
			fmt.Printf("  CREATE file %s\n", f.Path)
		}
		fmt.Println("  RUN    npm init -y")
		return nil
	}

	fmt.Printf("Initializing Node.js project: %s\n\n", projectName)

	cwd, _ := os.Getwd()
	for _, f := range files {
		full := filepath.Join(cwd, filepath.FromSlash(f.Path))
		dir := filepath.Dir(full)

		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("cannot create dir %s: %w", dir, err)
		}

		if _, err := os.Stat(full); err == nil && !initNodeForce {
			fmt.Printf("  SKIP   %s (exists)\n", f.Path)
			continue
		}

		if err := os.WriteFile(full, []byte(f.Content), 0644); err != nil {
			return fmt.Errorf("cannot write %s: %w", f.Path, err)
		}
		fmt.Printf("  CREATE %s\n", f.Path)
	}

	// Run npm init
	fmt.Println("\n  RUN    npm init -y")
	if err := initproject.RunNpmInit(); err != nil {
		fmt.Printf("  ⚠  npm init failed: %v\n", err)
		fmt.Println("  Run manually: npm init -y")
	}

	fmt.Printf("\n✔ Node.js project initialized: %s\n", projectName)
	fmt.Println("\nNext steps:")
	fmt.Println("  npm install   # install dependencies")
	fmt.Println("  npm start     # start the project")
	return nil
}
