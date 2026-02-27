<div align="center">

# ⚡ scaf

**Bootstrap projects in seconds. Standardize them for good.**

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue?style=flat-square)](LICENSE)
[![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20macOS%20%7C%20Linux-lightgrey?style=flat-square)](#installation)

</div>

---

## 🚀 Quick usage overview

```bash
# Full Go service bootstrap
scaf init go github.com/user/myapp --service --docker
scaf hooks --strict

# Or generate just the basics
scaf license MIT
scaf auto
scaf readme minimal
````

No copy-pasting.
No searching for templates.
No re-creating the same structure again.

---

## What is scaf?

`scaf` is a zero-config CLI that bootstraps and standardizes projects locally.

Generate:

* `LICENSE`
* `.gitignore`
* `README.md`
* `.editorconfig`
* `.dockerignore`
* Pre-commit hooks
* Full project structures

Safe by default. Offline-capable. No global side effects.

---

## Why scaf?

Starting a project should not require:

* Manually copying license text
* Searching for the correct `.gitignore`
* Writing boilerplate README files
* Reconfiguring git identity
* Rebuilding the same folder structure again

`scaf` eliminates that friction — in seconds.

---

## Core Capabilities

### 📄 File Generation

* SPDX license generation with author/year injection
* Stack-aware `.gitignore`
* Smart README templates (auto-filled from git & project metadata)
* `.editorconfig` and `.dockerignore` with auto-detection

### 🚀 Project Bootstrap

* `scaf init go` — CLI, service, or library layouts
* `scaf init node` — JS/TS, API-ready, optional Docker
* Architecture templates (clean, hexagonal, DDD, microservice, layered, etc.)
* Stack-aware pre-commit hooks (light or strict mode)

### 🔒 Safe & Local

* Never modifies global git config
* Never overwrites without `--force`
* Dry-run support everywhere
* Caches remote templates locally
* Works fully offline after first use

---

## Quick Start

### Generate project basics

```bash
scaf license MIT
scaf auto
scaf readme minimal
scaf editorconfig
```

### Scaffold a Go service

```bash
scaf init go github.com/user/app --service --docker
scaf hooks --strict
```

### Scaffold a Node.js app

```bash
scaf init node my-app --ts --docker
scaf hooks
```

### Switch git identity per repository

```bash
scaf git profile work
scaf git profile personal
```

---

## Installation

### Prebuilt binaries

Download from:
[https://github.com/TerraFaster/scaf/releases](https://github.com/TerraFaster/scaf/releases)

### Install with Go

```bash
go install github.com/TerraFaster/scaf@latest
```

Make sure `$(go env GOPATH)/bin` is in your `PATH`.

### Winget

```bash
winget install TerraFaster.Scaf
```

### Build from source

<details>
<summary>Click to see build instructions</summary>

```bash
git clone https://github.com/TerraFaster/scaf.git
cd scaf

make build
make install
make test
make build-all
```

Requires Go 1.21+

</details>

---

## Command Overview

All commands support:

```
--dry-run   preview changes
--force     overwrite existing files
--yes       non-interactive mode
--verbose   detailed output
```

### Common Commands

```bash
scaf license [mit|apache-2.0|...]
scaf ignore [stack]
scaf auto
scaf readme [template]
scaf editorconfig
scaf dockerignore
scaf hooks [--strict]
scaf init go <module>
scaf init node <name>
scaf structure [template]
scaf git profile <name>
scaf config
scaf update
```

### Detailed command usage

<details>
<summary>Click to see detailed command usage</summary>

### License

```bash
scaf license MIT --author "Author" --year 2026
```

### Gitignore

```bash
scaf ignore unity,dotnet
scaf auto
```

### README

Templates:
`minimal`, `cli`, `library`, `webapp`, `api`, `go`, `node`, `unity`, `custom`

```bash
scaf readme go --docker
```

### Init (Go)

```bash
scaf init go github.com/user/app --service --docker
```

Creates standard `cmd/`, `internal/`, `pkg/`, configs, Makefile, etc.

### Init (Node)

```bash
scaf init node app --ts --api --docker
```

Creates `src/`, `tests/`, ESLint, Prettier, configs, etc.

### Structure Templates

`layered`, `clean-architecture`, `hexagonal`, `ddd`,
`microservice`, `monolith`, `cli`, `minimal`

```bash
scaf structure clean-architecture
```

</details>

---

## Configuration

Global config: `~/.scaf/config.yaml`

Supports:

* Default license
* Default gitignore stacks
* Cache TTL
* Named git profiles
* Author defaults

<details>
<summary>Example config</summary>

```yaml
default_license: mit

default_ignore:
  - unity
  - dotnet

cache_ttl_hours: 168

profiles:
  work:
    name: Username
    email: username@company.com
  personal:
    name: user.name
    email: username@gmail.com
```

</details>

---

## User Templates

Override any template by placing files in:

```
~/.scaf/templates/licenses/
~/.scaf/templates/gitignore/
```

User templates always take priority.

---

## Security

* HTTPS-only external requests
* No global git modifications
* Safe file writes with traversal protection
* Never deletes files

---

## Contributing

Contributions are welcome.

```bash
make test
make build
```

Please keep changes focused and include tests where applicable.

---

## License

MIT License — see `LICENSE`.
