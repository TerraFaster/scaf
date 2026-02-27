// Package editorconfig generates .editorconfig files.
package editorconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Params holds editorconfig generation options.
type Params struct {
	Variant string
	// Overrides
	IndentStyle           string
	IndentSize            string
	EndOfLine             string
	InsertFinalNewline    bool
	TrimTrailingWhitespace bool
	Charset               string
}

// AutoDetectVariant detects the project type from files in cwd.
func AutoDetectVariant() string {
	cwd, _ := os.Getwd()
	if _, err := os.Stat(filepath.Join(cwd, "go.mod")); err == nil {
		return "go"
	}
	if _, err := os.Stat(filepath.Join(cwd, "package.json")); err == nil {
		return "node"
	}
	matches, _ := filepath.Glob(filepath.Join(cwd, "*.csproj"))
	if len(matches) > 0 {
		return "dotnet"
	}
	if info, err := os.Stat(filepath.Join(cwd, "Assets")); err == nil && info.IsDir() {
		return "unity"
	}
	return "default"
}

// Render generates the .editorconfig content.
func Render(p Params) string {
	base, ok := variants[p.Variant]
	if !ok {
		base = variants["default"]
	}

	// Apply overrides
	if p.IndentStyle != "" {
		base = replaceKey(base, "indent_style", p.IndentStyle)
	}
	if p.IndentSize != "" {
		base = replaceKey(base, "indent_size", p.IndentSize)
	}
	if p.EndOfLine != "" {
		base = replaceKey(base, "end_of_line", p.EndOfLine)
	}
	if p.Charset != "" {
		base = replaceKey(base, "charset", p.Charset)
	}
	return base
}

// AvailableVariants returns the list of supported variant names.
func AvailableVariants() []string {
	keys := make([]string, 0, len(variants))
	for k := range variants {
		keys = append(keys, k)
	}
	return keys
}

func replaceKey(content, key, value string) string {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, key+" =") || strings.HasPrefix(trimmed, key+"=") {
			lines[i] = fmt.Sprintf("%s = %s", key, value)
		}
	}
	return strings.Join(lines, "\n")
}

var variants = map[string]string{
	"default": `# EditorConfig — https://editorconfig.org
root = true

[*]
charset = utf-8
end_of_line = lf
insert_final_newline = true
trim_trailing_whitespace = true
indent_style = space
indent_size = 4

[*.md]
trim_trailing_whitespace = false

[Makefile]
indent_style = tab
`,

	"strict": `# EditorConfig — https://editorconfig.org
root = true

[*]
charset = utf-8
end_of_line = lf
insert_final_newline = true
trim_trailing_whitespace = true
indent_style = space
indent_size = 2
max_line_length = 120

[*.md]
trim_trailing_whitespace = false
max_line_length = off

[Makefile]
indent_style = tab
indent_size = 4
`,

	"go": `# EditorConfig — https://editorconfig.org
root = true

[*]
charset = utf-8
end_of_line = lf
insert_final_newline = true
trim_trailing_whitespace = true
indent_style = tab
indent_size = 4

[*.go]
indent_style = tab
indent_size = 4

[*.mod]
indent_style = tab
indent_size = 4

[*.md]
trim_trailing_whitespace = false

[Makefile]
indent_style = tab
indent_size = 4
`,

	"node": `# EditorConfig — https://editorconfig.org
root = true

[*]
charset = utf-8
end_of_line = lf
insert_final_newline = true
trim_trailing_whitespace = true
indent_style = space
indent_size = 2

[*.json]
indent_size = 2

[*.md]
trim_trailing_whitespace = false

[Makefile]
indent_style = tab
`,

	"unity": `# EditorConfig — https://editorconfig.org
root = true

[*]
charset = utf-8
end_of_line = crlf
insert_final_newline = false
trim_trailing_whitespace = true
indent_style = space
indent_size = 4

[*.cs]
indent_style = space
indent_size = 4
end_of_line = crlf

[*.md]
trim_trailing_whitespace = false
`,

	"dotnet": `# EditorConfig — https://editorconfig.org
root = true

[*]
charset = utf-8
end_of_line = crlf
insert_final_newline = true
trim_trailing_whitespace = true
indent_style = space
indent_size = 4

[*.cs]
indent_style = space
indent_size = 4

[*.{json,csproj,sln}]
indent_size = 2

[*.md]
trim_trailing_whitespace = false
`,

	"custom": `# EditorConfig — https://editorconfig.org
root = true

[*]
charset = utf-8
end_of_line = lf
insert_final_newline = true
trim_trailing_whitespace = true
indent_style = space
indent_size = 4

[*.md]
trim_trailing_whitespace = false
`,
}
