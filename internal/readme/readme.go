// Package readme generates README.md files from templates.
package readme

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

// Params holds all template variables.
type Params struct {
	Name        string
	Description string
	Author      string
	License     string
	Repo        string
	Badges      bool
	CI          bool
	Docker      bool
	Template    string
}

// AutoFill populates missing fields from the environment.
func (p *Params) AutoFill() {
	if p.Name == "" {
		cwd, _ := os.Getwd()
		p.Name = filepath.Base(cwd)
	}
	if p.Author == "" {
		p.Author = gitConfigGet("user.name")
	}
	if p.License == "" {
		if data, err := os.ReadFile("LICENSE"); err == nil {
			line := strings.SplitN(string(data), "\n", 2)[0]
			p.License = strings.TrimSpace(line)
		}
	}
	if p.Repo == "" {
		out, err := exec.Command("git", "remote", "get-url", "origin").Output()
		if err == nil {
			p.Repo = strings.TrimSpace(string(out))
		}
	}
}

// Render generates the README content for the given template name.
func Render(params Params) (string, error) {
	tmplStr, ok := readmeTemplates[params.Template]
	if !ok {
		tmplStr = readmeTemplates["minimal"]
	}

	tmpl, err := template.New("readme").Parse(tmplStr)
	if err != nil {
		return "", fmt.Errorf("template parse error: %w", err)
	}

	var sb strings.Builder
	if err := tmpl.Execute(&sb, params); err != nil {
		return "", fmt.Errorf("template execute error: %w", err)
	}
	return sb.String(), nil
}

// AvailableTemplates returns the list of supported template names.
func AvailableTemplates() []string {
	keys := make([]string, 0, len(readmeTemplates))
	for k := range readmeTemplates {
		keys = append(keys, k)
	}
	return keys
}

func gitConfigGet(key string) string {
	out, err := exec.Command("git", "config", key).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// bt is a backtick character — used to embed code fences inside Go template strings.
const bt = "`"

var readmeTemplates = map[string]string{
	"minimal": "# {{.Name}}\n\n" +
		"{{if .Description}}{{.Description}}{{else}}A short description of the project.{{end}}\n\n" +
		"## Installation\n\n" +
		bt + bt + bt + "sh\n# TODO\n" + bt + bt + bt + "\n\n" +
		"## Usage\n\n" +
		bt + bt + bt + "sh\n# TODO\n" + bt + bt + bt + "\n" +
		"{{if .License}}\n## License\n\n{{.License}}\n{{end}}" +
		"{{if .Author}}\n## Author\n\n{{.Author}}\n{{end}}",

	"cli": "# {{.Name}}\n\n" +
		"{{if .Description}}{{.Description}}{{else}}A command-line tool.{{end}}\n" +
		"{{if .Badges}}\n![Build](https://github.com/{{.Repo}}/actions/workflows/ci.yml/badge.svg)\n{{end}}\n" +
		"## Installation\n\n" +
		bt + bt + bt + "sh\ngo install {{.Repo}}@latest\n" + bt + bt + bt + "\n\n" +
		"## Usage\n\n" +
		bt + bt + bt + "\n{{.Name}} [command] [flags]\n\nCommands:\n  help    Show help\n  version Show version\n" + bt + bt + bt + "\n\n" +
		"## Flags\n\n" +
		"| Flag | Description | Default |\n|------|-------------|----------|\n| --help | Show help | - |\n\n" +
		"## Examples\n\n" +
		bt + bt + bt + "sh\n{{.Name}} --help\n" + bt + bt + bt + "\n" +
		"{{if .License}}\n## License\n\n{{.License}}\n{{end}}",

	"library": "# {{.Name}}\n\n" +
		"{{if .Description}}{{.Description}}{{else}}A reusable library.{{end}}\n\n" +
		"## Installation\n\n" +
		bt + bt + bt + "sh\ngo get {{.Repo}}\n" + bt + bt + bt + "\n\n" +
		"## Usage\n\n" +
		bt + bt + bt + "go\nimport \"{{.Repo}}\"\n\n// TODO: example\n" + bt + bt + bt + "\n\n" +
		"## API Reference\n\n### Functions\n\nTODO\n\n" +
		"## Contributing\n\nPull requests are welcome.\n" +
		"{{if .License}}\n## License\n\n{{.License}}\n{{end}}",

	"webapp": "# {{.Name}}\n\n" +
		"{{if .Description}}{{.Description}}{{else}}A web application.{{end}}\n" +
		"{{if .Badges}}\n![Build](https://github.com/{{.Repo}}/actions/workflows/ci.yml/badge.svg)\n{{end}}\n" +
		"## Getting Started\n\n### Prerequisites\n\n- Node.js >= 18\n- npm or yarn\n\n" +
		"### Installation\n\n" +
		bt + bt + bt + "sh\nnpm install\n" + bt + bt + bt + "\n\n" +
		"### Development\n\n" +
		bt + bt + bt + "sh\nnpm run dev\n" + bt + bt + bt + "\n\n" +
		"### Build\n\n" +
		bt + bt + bt + "sh\nnpm run build\n" + bt + bt + bt + "\n" +
		"{{if .Docker}}\n## Docker\n\n" +
		bt + bt + bt + "sh\ndocker build -t {{.Name}} .\ndocker run -p 3000:3000 {{.Name}}\n" + bt + bt + bt + "\n{{end}}" +
		"{{if .License}}\n## License\n\n{{.License}}\n{{end}}",

	"api": "# {{.Name}}\n\n" +
		"{{if .Description}}{{.Description}}{{else}}A REST API service.{{end}}\n" +
		"{{if .Badges}}\n![Build](https://github.com/{{.Repo}}/actions/workflows/ci.yml/badge.svg)\n{{end}}\n" +
		"## API Endpoints\n\n" +
		"| Method | Path | Description |\n|--------|------|-------------|\n| GET    | /health | Health check |\n\n" +
		"## Getting Started\n\n### Prerequisites\n\nTODO\n\n### Running\n\n" +
		bt + bt + bt + "sh\n# TODO\n" + bt + bt + bt + "\n" +
		"{{if .Docker}}\n## Docker\n\n" +
		bt + bt + bt + "sh\ndocker build -t {{.Name}} .\ndocker run -p 8080:8080 {{.Name}}\n" + bt + bt + bt + "\n{{end}}\n" +
		"## Configuration\n\n| Variable | Description | Default |\n|----------|-------------|----------|\n| PORT | Server port | 8080 |\n" +
		"{{if .License}}\n## License\n\n{{.License}}\n{{end}}",

	"go": "# {{.Name}}\n\n" +
		"{{if .Description}}{{.Description}}{{else}}A Go project.{{end}}\n" +
		"{{if .Badges}}\n[![Go Reference](https://pkg.go.dev/badge/{{.Repo}}.svg)](https://pkg.go.dev/{{.Repo}})\n{{end}}\n" +
		"## Requirements\n\n- Go 1.21+\n\n" +
		"## Installation\n\n" +
		bt + bt + bt + "sh\ngo get {{.Repo}}\n" + bt + bt + bt + "\n\n" +
		"## Usage\n\n" +
		bt + bt + bt + "go\n// TODO\n" + bt + bt + bt + "\n" +
		"{{if .Docker}}\n## Docker\n\n" +
		bt + bt + bt + "sh\ndocker build -t {{.Name}} .\n" + bt + bt + bt + "\n{{end}}\n" +
		"## Development\n\n" +
		bt + bt + bt + "sh\nmake build\nmake test\nmake lint\n" + bt + bt + bt + "\n" +
		"{{if .License}}\n## License\n\n{{.License}}\n{{end}}",

	"node": "# {{.Name}}\n\n" +
		"{{if .Description}}{{.Description}}{{else}}A Node.js project.{{end}}\n\n" +
		"## Requirements\n\n- Node.js >= 18\n\n" +
		"## Installation\n\n" +
		bt + bt + bt + "sh\nnpm install\n" + bt + bt + bt + "\n\n" +
		"## Scripts\n\n" +
		"| Script | Description |\n|--------|-------------|\n| npm start | Start the application |\n| npm test  | Run tests |\n| npm run build | Build for production |\n" +
		"{{if .License}}\n## License\n\n{{.License}}\n{{end}}",

	"unity": "# {{.Name}}\n\n" +
		"{{if .Description}}{{.Description}}{{else}}A Unity project.{{end}}\n\n" +
		"## Requirements\n\n- Unity 2022.x or newer\n\n" +
		"## Getting Started\n\n1. Clone the repository\n2. Open in Unity Hub\n3. Open the project\n\n" +
		"## Project Structure\n\n" +
		bt + bt + bt + "\nAssets/\n  Scripts/\n  Scenes/\n  Prefabs/\n  Materials/\n" + bt + bt + bt + "\n" +
		"{{if .License}}\n## License\n\n{{.License}}\n{{end}}",

	"custom": "# {{.Name}}\n\n" +
		"{{if .Description}}{{.Description}}{{end}}\n\n" +
		"## Overview\n\nTODO\n\n## Getting Started\n\nTODO\n\n## Contributing\n\nTODO\n" +
		"{{if .License}}\n## License\n\n{{.License}}\n{{end}}",
}
