// Package dockerignore generates .dockerignore files.
package dockerignore

import (
	"os"
	"path/filepath"
	"strings"
)

// AutoDetectStack detects the project stack from files in cwd.
func AutoDetectStack() string {
	cwd, _ := os.Getwd()
	if _, err := os.Stat(filepath.Join(cwd, "go.mod")); err == nil {
		return "go"
	}
	if _, err := os.Stat(filepath.Join(cwd, "package.json")); err == nil {
		return "node"
	}
	if _, err := os.Stat(filepath.Join(cwd, "requirements.txt")); err == nil {
		return "python"
	}
	matches, _ := filepath.Glob(filepath.Join(cwd, "*.csproj"))
	if len(matches) > 0 {
		return "dotnet"
	}
	if info, err := os.Stat(filepath.Join(cwd, "Assets")); err == nil && info.IsDir() {
		return "unity"
	}
	return "generic"
}

// Render generates the .dockerignore content for a given stack.
func Render(stack string) string {
	content, ok := stacks[strings.ToLower(stack)]
	if !ok {
		content = stacks["generic"]
	}
	return content
}

// AvailableStacks returns the list of supported stack names.
func AvailableStacks() []string {
	keys := make([]string, 0, len(stacks))
	for k := range stacks {
		keys = append(keys, k)
	}
	return keys
}

var stacks = map[string]string{
	"go": `# .dockerignore — Go
.git
.gitignore
.github
.editorconfig
*.md
*.test
*.out
coverage.txt
coverage.html

# Build artifacts
/bin/
/dist/
/tmp/

# IDE / OS
.idea/
.vscode/
*.DS_Store
Thumbs.db

# Test / dev files
testdata/
*_test.go
`,

	"node": `# .dockerignore — Node.js
.git
.gitignore
.github
.editorconfig
*.md
README*

# Dependencies (reinstalled in container)
node_modules/
npm-debug.log*
yarn-debug.log*
yarn-error.log*

# Build output
dist/
build/
.next/
.nuxt/
.output/

# Test / dev
coverage/
.nyc_output/
*.test.js
*.spec.js
__tests__/
jest.config.*

# IDE / OS
.idea/
.vscode/
*.DS_Store
Thumbs.db

# Env files
.env
.env.local
.env.*.local
`,

	"python": `# .dockerignore — Python
.git
.gitignore
.github
.editorconfig
*.md

# Virtual envs
venv/
.venv/
env/
__pycache__/
*.pyc
*.pyo
*.pyd
.Python

# Dist / build
build/
dist/
*.egg-info/
.eggs/

# Tests
.pytest_cache/
.coverage
htmlcov/
*.test

# IDE / OS
.idea/
.vscode/
*.DS_Store
Thumbs.db

# Env
.env
.env.*
`,

	"dotnet": `# .dockerignore — .NET
.git
.gitignore
.github
.editorconfig
*.md

# Build output
bin/
obj/

# IDE
.idea/
.vscode/
*.user
*.suo
*.vs/
.vs/

# Test results
TestResults/

# OS
*.DS_Store
Thumbs.db
`,

	"unity": `# .dockerignore — Unity
.git
.gitignore
.github
*.md

# Unity generated
Library/
Temp/
Obj/
Build/
Builds/
Logs/
UserSettings/

# IDE
.idea/
.vscode/
*.DS_Store
Thumbs.db
`,

	"generic": `# .dockerignore — Generic
.git
.gitignore
.github
.editorconfig
*.md
README*
CHANGELOG*
LICENSE*

# IDE / OS
.idea/
.vscode/
*.DS_Store
Thumbs.db

# Temp / logs
tmp/
temp/
logs/
*.log

# Test
tests/
test/
coverage/
`,
}
