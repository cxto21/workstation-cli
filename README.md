<div align="center">

<img src="https://workstationai.pages.dev/logo.png" alt="Workstation Logo" width="200">

# Workstation CLI

### The Team Operating Environment
Control plane for AI-driven teams, shared state, and agent execution.

[![Rust](https://img.shields.io/badge/rust-stable-orange?style=for-the-badge&logo=rust)](https://www.rust-lang.org/)
[![License](https://img.shields.io/badge/license-MIT-blue?style=for-the-badge)](LICENSE)
[![Stars](https://img.shields.io/github/stars/cxto21/workstation-cli?style=for-the-badge&color=gold)](https://github.com/cxto21/workstation-cli/stargazers)

[Quick Start](#quick-start) • [Core Commands](#core-commands) • [Features](#features) • [Development](#development) • [Resources](#resources)

</div>

## What Is Workstation?

Workstation CLI unifies:
- Terminal runtime for multi-session execution
- SSOT synchronization through Git
- Team-oriented workspace organization
- Local docs and knowledge indexing

It is designed for humans supervising AI-assisted work in a reproducible terminal environment.

## Quick Start

One-command install:

```bash
curl -fsSL https://raw.githubusercontent.com/cxto21/workstation-cli/main/install.sh | bash
```

Manual build:

```bash
git clone https://github.com/cxto21/workstation-cli.git
cd workstation-cli
cargo build --release
install -m 755 target/release/workstation-cli ~/.local/bin/workstation-cli
workstation-cli --version
```

## Core Commands

```bash
# Launch TUI
workstation-cli

# Team state sync
workstation-cli sync status
workstation-cli sync pull
workstation-cli sync push

# SSOT bootstrap and validation
workstation-cli ssot init --org "Your Org" --member "Your Name" --project "core-platform"
workstation-cli ssot validate

# Integration checks
workstation-cli tools status
workstation-cli tools init
```

## Features

- Offices, desks, and tabs for structured work
- Jump mode navigation with minimal shortcut conflicts
- Runtime persistence across reconnects
- Multi-client terminal attachment
- Local docs and KB indexing inside the workspace
- Theme system and onboarding templates

## Development

```bash
cargo build
cargo test
```

Go compatibility bootstrap:

```bash
cd go
go test ./...
```

Contributions are welcome through pull requests on:
https://github.com/cxto21/workstation-cli

## Resources

- [Keyboard Shortcuts](docs/KEYBOARD_SHORTCUTS.md)
- [AI Agent Friendly Notes](docs/AI_AGENT_FRIENDLY.md)
- [Testing Guide](docs/TESTING.md)
- [Terminal Persistence](docs/TERMINAL_PERSISTENCE.md)
- [Roadmap](docs/todos/roadmap.md)

---

<div align="center">

Built for clear, team-first execution in the terminal.

[Star the project](https://github.com/cxto21/workstation-cli) • [Report issues](https://github.com/cxto21/workstation-cli/issues)

</div>
