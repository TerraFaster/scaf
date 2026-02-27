// Package structure generates project directory structures.
package structure

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Template defines a structure template.
type Template struct {
	Name        string
	Description string
	Dirs        []string
	Files       map[string]string
}

// Get returns the template by name.
func Get(name string) (Template, bool) {
	t, ok := templates[name]
	return t, ok
}

// AvailableTemplates returns all available template names.
func AvailableTemplates() []string {
	keys := make([]string, 0, len(templates))
	for k := range templates {
		keys = append(keys, k)
	}
	return keys
}

// Apply creates the directory structure in cwd.
func Apply(t Template, dryRun bool) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	created := []string{}

	for _, dir := range t.Dirs {
		full := filepath.Join(cwd, filepath.FromSlash(dir))
		if dryRun {
			fmt.Printf("  CREATE dir  %s/\n", dir)
			continue
		}
		if err := os.MkdirAll(full, 0755); err != nil {
			return fmt.Errorf("cannot create dir %s: %w", dir, err)
		}
		created = append(created, dir+"/")
	}

	for path, content := range t.Files {
		full := filepath.Join(cwd, filepath.FromSlash(path))
		if dryRun {
			fmt.Printf("  CREATE file %s\n", path)
			continue
		}
		if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
			return fmt.Errorf("cannot create dir for %s: %w", path, err)
		}
		if _, err := os.Stat(full); os.IsNotExist(err) {
			if err := os.WriteFile(full, []byte(content), 0644); err != nil {
				return fmt.Errorf("cannot write %s: %w", path, err)
			}
			created = append(created, path)
		}
	}

	if !dryRun {
		for _, p := range created {
			if strings.HasSuffix(p, "/") {
				fmt.Printf("  ✔ mkdir  %s\n", p)
			} else {
				fmt.Printf("  ✔ create %s\n", p)
			}
		}
	}
	return nil
}

var templates = map[string]Template{
	"layered": {
		Name:        "layered",
		Description: "Classic layered (presentation / business / data) structure",
		Dirs: []string{
			"api", "service", "repository", "model", "config", "middleware",
		},
		Files: map[string]string{
			"README.md": "# Layered Architecture\n\nLayers: api → service → repository → model\n",
		},
	},
	"clean-architecture": {
		Name:        "clean-architecture",
		Description: "Clean Architecture (entities, use cases, adapters, frameworks)",
		Dirs: []string{
			"entities", "usecases", "adapters/controllers",
			"adapters/presenters", "adapters/gateways",
			"frameworks/web", "frameworks/db",
		},
		Files: map[string]string{
			"README.md": "# Clean Architecture\n\nRings: entities → use cases → adapters → frameworks\n",
		},
	},
	"hexagonal": {
		Name:        "hexagonal",
		Description: "Hexagonal (Ports & Adapters) architecture",
		Dirs: []string{
			"domain", "application/ports", "application/services",
			"infrastructure/adapters", "infrastructure/persistence",
			"infrastructure/http",
		},
		Files: map[string]string{
			"README.md": "# Hexagonal Architecture\n\nDomain → Application Ports → Infrastructure Adapters\n",
		},
	},
	"ddd": {
		Name:        "ddd",
		Description: "Domain-Driven Design structure",
		Dirs: []string{
			"domain/aggregates", "domain/entities", "domain/valueobjects",
			"domain/repositories", "domain/events",
			"application/commands", "application/queries",
			"infrastructure/persistence", "infrastructure/messaging",
			"interfaces/http", "interfaces/grpc",
		},
		Files: map[string]string{
			"README.md": "# DDD Structure\n\nDomain / Application / Infrastructure / Interfaces\n",
		},
	},
	"microservice": {
		Name:        "microservice",
		Description: "Microservice project structure",
		Dirs: []string{
			"cmd/server", "internal/handler", "internal/service",
			"internal/repository", "internal/domain",
			"api/proto", "api/openapi",
			"deployments/k8s", "deployments/helm",
			"configs", "scripts",
		},
		Files: map[string]string{
			"README.md":             "# Microservice\n\nStandard microservice layout.\n",
			"deployments/k8s/.gitkeep": "",
		},
	},
	"monolith": {
		Name:        "monolith",
		Description: "Monolithic application structure",
		Dirs: []string{
			"cmd/app", "internal/api", "internal/core",
			"internal/infra/db", "internal/infra/cache",
			"web/templates", "web/static",
			"configs", "scripts", "migrations",
		},
		Files: map[string]string{
			"README.md": "# Monolith\n\nStandard monolithic layout.\n",
		},
	},
	"cli": {
		Name:        "cli",
		Description: "CLI application structure",
		Dirs: []string{
			"cmd", "internal/command", "internal/config",
			"internal/output", "scripts",
		},
		Files: map[string]string{
			"README.md": "# CLI Application\n\nCommand-line tool structure.\n",
		},
	},
	"minimal": {
		Name:        "minimal",
		Description: "Minimal project structure",
		Dirs: []string{
			"src", "tests", "docs",
		},
		Files: map[string]string{
			"README.md": "# Project\n\nMinimal project structure.\n",
		},
	},
}
