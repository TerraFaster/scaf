<div align="center">

# ⚡ scaf

**Scaffold standard project files in seconds.**

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue?style=flat-square)](LICENSE)
[![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20macOS%20%7C%20Linux-lightgrey?style=flat-square)](#installation)

</div>

---

## What is scaf?

`scaf` is a zero-config CLI tool that generates the boring-but-critical files every project needs — `LICENSE` and `.gitignore` — with the right content, instantly. It pulls templates from [GitHub's License API](https://api.github.com/licenses) and [gitignore.io](https://www.toptal.com/developers/gitignore), caches them locally, and works fully offline after the first sync.

No more copy-pasting from the internet. No more forgetting the copyright year. Just `scaf`.

Usage:
```bash
scaf license                     # interactive license selector
scaf ignore                      # interactive .gitignore selector
scaf auto                        # auto-detects project type → generates .gitignore
```

Or specify license and templates directly:
```bash
scaf license MIT                 # generates LICENSE with MIT text
scaf ignore python,dotnet        # generates .gitignore for Python + .NET
```


---

## Features

- 📄 **LICENSE generation** — all SPDX licenses from GitHub, with year & author substitution
- 🚫 **.gitignore generation** — single or combined templates (e.g. `unity,dotnet`)
- 🤖 **Auto-detection** — scans your directory and picks the right templates automatically
- 🔍 **Fuzzy search** — typo in the license name? scaf suggests what you meant
- 💡 **Interactive TUI** — arrow-key navigation with live filtering when no name is given
- 💾 **Smart caching** — templates cached for 7 days; works offline after first run
- 📝 **Per-project config** — `scaf.yaml` in the project root stores your defaults
- 🔧 **User templates** — drop files in `~/.scaf/templates/` to override any template
- 🔒 **Secure** — TLS-only requests, path-traversal protection, no code execution

---

## Installation

### Pre-built Binaries

Download the latest release for your platform from the [Releases page](https://github.com/TerraFaster/scaf/releases):

| Platform | Architecture | File |
|---|---|---|
| Linux | x86_64 | `scaf-linux-amd64` |
| Linux | ARM64 | `scaf-linux-arm64` |
| macOS | Intel | `scaf-darwin-amd64` |
| macOS | Apple Silicon | `scaf-darwin-arm64` |
| Windows | x86_64 | `scaf-windows-amd64.exe` |
| Windows | ARM64 | `scaf-windows-arm64.exe` |

---

### macOS & Linux

```bash
# Download (adjust PLATFORM as needed: linux-amd64, darwin-arm64, etc.)
curl -L https://github.com/TerraFaster/scaf/releases/latest/download/scaf-linux-amd64 -o scaf

# Make executable
chmod +x scaf

# Move to a directory already in your PATH
sudo mv scaf /usr/local/bin/scaf

# Verify
scaf --help
```

---

### Windows

**Step 1 — Download** `scaf-windows-amd64.exe` from the [Releases page](https://github.com/TerraFaster/scaf/releases) and rename it to `scaf.exe`.

**Step 2 — Place it somewhere convenient**, e.g. `C:\Tools\scaf\scaf.exe`.

**Step 3 — Add that folder to your PATH** so you can run `scaf` from any terminal:

> **GUI:** Start → search *"Edit the system environment variables"* → click **Environment Variables…** → under *User variables* select **Path** → **Edit** → **New** → paste `C:\Tools\scaf` → OK × 3.

> **PowerShell (one-liner, current user):**
> ```powershell
> $dir = "$HOME\bin"
> New-Item -ItemType Directory -Force -Path $dir | Out-Null
> Copy-Item .\scaf.exe "$dir\scaf.exe"
> [Environment]::SetEnvironmentVariable("Path", $env:Path + ";$dir", "User")
> ```
> Restart your terminal, then run `scaf --help`.

> **Windows Package Manager (coming soon):**
> ```powershell
> winget install scaf
> ```

---

### Install with Go

If you have Go 1.21+ installed:

```bash
go install github.com/TerraFaster/scaf@latest
```

This places `scaf` in `$(go env GOPATH)/bin`. Make sure that directory is in your `PATH`:

```bash
# Bash / Zsh — add to ~/.bashrc or ~/.zshrc:
export PATH="$PATH:$(go env GOPATH)/bin"

# Fish:
fish_add_path (go env GOPATH)/bin

# PowerShell (Windows):
$env:Path += ";$(go env GOPATH)\bin"
```

---

### Build from Source

```bash
git clone https://github.com/TerraFaster/scaf.git
cd scaf

make build          # produces ./scaf binary
make install        # copies to /usr/local/bin (Unix)
make test           # run all tests
make cross          # build for all platforms → ./dist/
```

**Requirements:** Go 1.21 or later

---

## Usage

### `scaf license`

Generate a `LICENSE` file.

```bash
# Interactive — searchable list of all SPDX licenses
scaf license

# Direct — specify by key (case-insensitive)
scaf license MIT
scaf license apache-2.0
scaf license gpl-3.0

# With metadata
scaf license MIT --author "AUTHOR" --year 2026

# Preview without writing
scaf license MIT --dry-run

# Overwrite existing file silently
scaf license MIT --force

# Backup existing file before overwrite
scaf license MIT --backup

# Skip all confirmation prompts
scaf license MIT --yes
```

---

### `scaf ignore` / `scaf gitignore`

Generate a `.gitignore` file. Both commands are identical aliases.

```bash
# Interactive multi-select from popular templates
scaf ignore

# Single template
scaf ignore python

# Multiple templates — comma-separated or space-separated
scaf ignore unity,dotnet
scaf ignore python node rust

# Preview
scaf ignore go --dry-run
```

---

### `scaf auto`

Auto-detect your project stack and generate `.gitignore` automatically — no arguments needed.

```bash
scaf auto
```

Detection rules:

| Marker | Template |
|---|---|
| `package.json` | `node` |
| `Cargo.toml` | `rust` |
| `go.mod` | `go` |
| `requirements.txt` | `python` |
| `*.csproj` | `dotnet` |
| `Assets/` (directory) | `unity` |

Multiple matches produce a single combined `.gitignore`. Available flags:

```bash
scaf auto --dry-run     # preview only, no file written
scaf auto --force       # overwrite existing .gitignore
scaf auto --backup      # save .gitignore.bak before overwriting
scaf auto --yes         # no confirmation prompts
```

---

### `scaf config`

View or edit the global config file (`~/.scaf/config.yaml`).
The file is created automatically on the first run — no setup step needed.

```bash
scaf config           # print path + current contents
scaf config --edit    # open in $EDITOR (falls back to system default)
```
---

### `scaf update`

Force-refresh all cached templates (ignores TTL).

```bash
scaf update
```

---

## Configuration (`~/.scaf/config.yaml`)

scaf keeps a **single, user-wide** config file at `~/.scaf/config.yaml`.
It is created automatically with defaults the very first time you run any `scaf` command — no setup needed.
Edit it with `scaf config --edit` or open the file directly in any text editor.

```yaml
# scaf configuration file
# https://github.com/TerraFaster/scaf

# Default license key — used by `scaf license` when no argument is given
default_license: mit

# Default .gitignore templates — used by `scaf ignore` when no argument is given
default_ignore:
  - unity
  - dotnet

# Template cache TTL in hours (default: 168 = 7 days)
cache_ttl_hours: 168

# Launch interactive mode even when a default is configured
interactive: true

# Author name for LICENSE placeholders (fallback: git config user.name)
default_author: ""

# Year for LICENSE placeholders (fallback: current year)
default_year: ""
```

**Resolution priority for each value:**

```
CLI flag  >  scaf.yaml  >  git config / system default
```

---

## User Templates

Override any template with your own version by placing files in:

```
~/.scaf/templates/licenses/<key>    # e.g. "mit", "apache-2.0"
~/.scaf/templates/gitignore/<name>  # e.g. "python", "mycompany"
```

User templates always take priority over downloaded ones. The file contents are used verbatim.

---

## Cache

Templates are stored in `~/.scaf/cache/` and reused until the TTL expires.

```
~/.scaf/
├── cache/
│   ├── licenses/          # license list + individual bodies
│   └── gitignore/         # combined gitignore responses
└── templates/
    ├── licenses/           # user overrides
    └── gitignore/          # user overrides
```

Run `scaf update` to force a refresh. Adjust `cache_ttl_hours` in `scaf.yaml` to control expiry.

---

## Typical Workflows

**Starting a new project:**
```bash
mkdir my-app && cd my-app
git init
scaf license MIT                 # generate LICENSE
scaf auto                        # auto-detect stack → .gitignore
```

**Unity + .NET game project:**
```bash
scaf license apache-2.0 --author "ACME Corp"
scaf ignore unity,dotnet,visualstudio
```

**Python microservice:**
```bash
scaf license MIT
scaf ignore python,docker,linux
```

**Preview everything first:**
```bash
scaf license MIT --dry-run
scaf auto --dry-run
```

---

## Developer Guide

### Project Structure

```
scaf/
├── main.go
├── Makefile
├── go.mod
├── cmd/
│   ├── root.go          # cobra root, config loading, init & update commands
│   ├── license.go       # scaf license
│   ├── ignore.go        # scaf ignore / scaf gitignore
│   ├── auto.go          # scaf auto (standalone auto-detect command)
│   └── helpers.go       # autoDetectProject(), runCommand()
└── internal/
    ├── config/
    │   └── config.go    # ~/.scaf/config.yaml — auto-created on first run
    ├── templates/
    │   ├── license_provider.go   # GitHub License API + fuzzy search
    │   ├── ignore_provider.go    # gitignore.io API
    │   └── http.go               # TLS-enforced HTTP client
    ├── cache/
    │   └── cache.go     # TTL cache backed by ~/.scaf/cache/
    ├── ui/
    │   └── prompt.go    # Bubbletea TUI: SelectOne, MultiSelect with live fuzzy filter
    └── fs/
        └── writer.go    # safe file writer (backup, force, path-traversal guard)
```

### Running Tests

```bash
make test            # run all tests quietly
make test-verbose    # run with full output
```

Test coverage includes:

- `internal/cache` — TTL expiry, cache miss, path-traversal protection, Clear()
- `internal/config` — load defaults, round-trip save/load, Init with force flag
- `internal/templates` — license JSON parsing, fuzzy search, user overrides, CSV list parsing
- `internal/fs` — create file, force overwrite, backup mode, path traversal blocked

### Cross-Compilation

```bash
make cross
```

Or run any of build scripts depending on your OS:
```bash
build.bat      # Windows
build.ps1      # Windows PowerShell
build.sh       # Linux/macOS
```

Produces self-contained binaries in `./dist/` for all six targets (Linux/macOS/Windows × amd64/arm64).

### Adding a New Template Source

1. Create `internal/templates/yourtype_provider.go` with `List()` and `Get()` methods
2. Use `cache.Cache` for TTL-based persistence
3. Add `cmd/yourtype.go` with a `cobra.Command`
4. Register it in `cmd/root.go → Execute()`
5. Write tests in `internal/templates/yourtype_provider_test.go`

### Key Dependencies

| Package | Purpose |
|---|---|
| [`spf13/cobra`](https://github.com/spf13/cobra) | CLI framework |
| [`charmbracelet/bubbletea`](https://github.com/charmbracelet/bubbletea) | TUI engine |
| [`charmbracelet/bubbles`](https://github.com/charmbracelet/bubbles) | TUI components (text input) |
| [`charmbracelet/lipgloss`](https://github.com/charmbracelet/lipgloss) | TUI styling |
| [`lithammer/fuzzysearch`](https://github.com/lithammer/fuzzysearch) | Fuzzy filtering |
| [`gopkg.in/yaml.v3`](https://pkg.go.dev/gopkg.in/yaml.v3) | Config file parsing |

---

## Security

- All network requests use **HTTPS with TLS verification** (Go's default `http.Client`, no `InsecureSkipVerify`)
- Only **GET requests** are ever made — no data is sent to external services
- All file paths are **validated against path traversal** before any read or write
- **No external code is executed** at any point
- Cache keys are sanitized to `[a-zA-Z0-9\-_.]` before being used as filenames

---

## Contributing

Contributions are welcome and appreciated!

1. Fork the repository
2. Create a feature branch: `git checkout -b feat/amazing-feature`
3. Make your changes and add tests
4. Ensure everything passes: `make test`
5. Commit with a clear message and open a Pull Request

Please keep commits focused, follow the existing code style, and update tests for any behaviour change.

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
