package initproject

import (
	"fmt"
	"os"
	"os/exec"
)

// NodeMode defines the type of Node project.
type NodeMode string

const (
	NodeModeDefault NodeMode = "default"
	NodeModeTS      NodeMode = "ts"
	NodeModeAPI     NodeMode = "api"
	NodeModeCLI     NodeMode = "cli"
	NodeModeLib     NodeMode = "lib"
)

// NodeOptions holds options for Node project init.
type NodeOptions struct {
	ProjectName string
	Mode        NodeMode
	Docker      bool
	DryRun      bool
}

// NodeFiles returns all files to create for a Node project.
func NodeFiles(opts NodeOptions) []FileEntry {
	name := opts.ProjectName
	files := []FileEntry{
		{
			Path:    "src/index.js",
			Content: fmt.Sprintf("'use strict';\n\nconsole.log('Hello from %s!');\n", name),
		},
		{
			Path:    "tests/.gitkeep",
			Content: "",
		},
		{
			Path:    "configs/.gitkeep",
			Content: "",
		},
		{
			Path: ".gitignore",
			Content: `node_modules/
dist/
build/
.env
.env.local
coverage/
*.log
.DS_Store
Thumbs.db
`,
		},
		{
			Path: ".editorconfig",
			Content: `root = true

[*]
charset = utf-8
end_of_line = lf
insert_final_newline = true
trim_trailing_whitespace = true
indent_style = space
indent_size = 2

[*.md]
trim_trailing_whitespace = false
`,
		},
		{
			Path: "README.md",
			Content: fmt.Sprintf(`# %s

A Node.js project.

## Installation

`+"```"+`sh
npm install
`+"```"+`

## Usage

`+"```"+`sh
npm start
`+"```"+`
`, name),
		},
		{
			Path: ".eslintrc.json",
			Content: `{
  "env": {
    "node": true,
    "es2021": true
  },
  "extends": "eslint:recommended",
  "parserOptions": {
    "ecmaVersion": "latest"
  }
}
`,
		},
		{
			Path: ".prettierrc",
			Content: `{
  "semi": true,
  "singleQuote": true,
  "tabWidth": 2,
  "trailingComma": "es5"
}
`,
		},
	}

	if opts.Mode == NodeModeTS {
		files[0] = FileEntry{
			Path:    "src/index.ts",
			Content: fmt.Sprintf("console.log('Hello from %s!');\n", name),
		}
		files = append(files, FileEntry{
			Path: "tsconfig.json",
			Content: `{
  "compilerOptions": {
    "target": "ES2020",
    "module": "commonjs",
    "lib": ["ES2020"],
    "outDir": "./dist",
    "rootDir": "./src",
    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "forceConsistentCasingInFileNames": true,
    "resolveJsonModule": true,
    "declaration": true,
    "declarationMap": true,
    "sourceMap": true
  },
  "include": ["src/**/*"],
  "exclude": ["node_modules", "dist"]
}
`,
		})
	}

	if opts.Docker {
		files = append(files, FileEntry{
			Path: "Dockerfile",
			Content: fmt.Sprintf(`FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production

FROM node:18-alpine
WORKDIR /app
COPY --from=builder /app/node_modules ./node_modules
COPY . .
EXPOSE 3000
CMD ["node", "src/index.js"]
`),
		})
		files = append(files, FileEntry{
			Path: ".dockerignore",
			Content: `node_modules/
.git
*.md
coverage/
.env
`,
		})
	}

	return files
}

// RunNpmInit runs npm init -y in the current directory.
func RunNpmInit() error {
	cmd := exec.Command("npm", "init", "-y")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
