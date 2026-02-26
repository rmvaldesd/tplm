# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Purpose

`tplm` (tmux Project Layout Manager) — a Go CLI tool that:
- Defines project workspaces with pane/window layouts via YAML config
- Provides a Bubbletea-powered floating picker (inside `tmux display-popup`) to browse and switch between projects and active sessions
- Auto-creates tmux sessions from config with layout application and startup commands

## Technology Stack

- **Language**: Go
- **TUI framework**: [Bubbletea](https://github.com/charmbracelet/bubbletea) — Elm-architecture TUI (Model/Update/View)
- **TUI components**: [Bubbles](https://github.com/charmbracelet/bubbles) — textinput for rename
- **Styling**: [Lipgloss](https://github.com/charmbracelet/lipgloss) — layout and color styling
- **tmux control**: `os/exec` calling `tmux` CLI commands (no external tmux library needed)
- **Config format**: YAML via `gopkg.in/yaml.v3`
- **CLI entrypoint**: [Cobra](https://github.com/spf13/cobra) for subcommands

## Architecture

```
tplm/
├── cmd/
│   └── tplm/
│       └── main.go              # package main — imports internal/cli
├── internal/
│   ├── cli/                     # Cobra commands
│   │   ├── root.go              # Root cobra command, --config flag, config loading
│   │   ├── picker.go            # `tplm picker` — launches Bubbletea TUI
│   │   ├── open.go              # `tplm open <project>` — create session + switch (no TUI)
│   │   ├── list.go              # `tplm list` — print projects/sessions to stdout
│   │   └── init.go              # `tplm init` — generate example config.yaml
│   ├── config/
│   │   ├── config.go            # Parse YAML config, resolve ~ paths, example config
│   │   └── types.go             # Config, Project, Layout, Window, Pane, OnStart structs
│   ├── tmux/
│   │   ├── exec.go              # Low-level: run tmux commands via exec.Command
│   │   ├── session.go           # Create/kill/rename/switch sessions, send-keys
│   │   ├── layout.go            # Apply layout: split-window, resize, on_start commands
│   │   └── query.go             # List sessions with metadata, check session existence
│   └── ui/
│       ├── picker.go            # Bubbletea model: two-section picker (projects + sessions)
│       ├── rename.go            # Bubbletea model: inline rename text input
│       ├── keys.go              # Key bindings: j/k, Enter, d, r, q/Esc
│       └── styles.go            # Lipgloss style definitions
├── go.mod
├── Makefile
└── tmux.conf.example            # Recommended keybindings
```

## Key Patterns

**tmux integration** — all tmux control goes through `internal/tmux/exec.go` using `exec.Command("tmux", ...)`. Never call tmux directly from UI or config code.

**Floating popup flow** — the picker runs inside a `tmux display-popup`:
```tmux
bind-key C-p display-popup -E -w 80% -h 60% "tplm picker"
```
The picker Bubbletea app runs inside the popup, and on selection it executes `tmux switch-client -t <session>` then exits.

**Bubbletea model structure** — the picker has three modes: normal, confirmKill, rename. The rename mode delegates to a sub-model (RenameModel) wrapping bubbles/textinput.

**Picker two-section list** — projects on top, active sessions on bottom. A single cursor index spans both sections. Session actions (kill, rename) only work on sessions.

**Config file** (`~/.config/tplm/config.yaml`):
```yaml
projects:
  - name: my-api
    path: ~/Projects/my-api
    layout: dev
    on_start:
      - window: editor
        command: nvim .
      - window: server
        command: go run ./cmd/server

layouts:
  dev:
    windows:
      - name: editor
        panes:
          - size: "70%"
          - split: horizontal
            size: "30%"
      - name: server
        panes:
          - size: "100%"
```

## CLI Commands

| Command | Description |
|---|---|
| `tplm picker` | Open the Bubbletea picker (inside `tmux display-popup`) |
| `tplm open <name>` | Create session from project config and switch (no TUI) |
| `tplm list` | Print projects and active sessions to stdout |
| `tplm init` | Generate a starter `~/.config/tplm/config.yaml` |

## Development Commands

```bash
# Run directly
go run ./cmd/tplm picker

# Build binary
go build -o tplm ./cmd/tplm

# Install to PATH
go install ./cmd/tplm

# Run tests
go test ./...

# Run a single package's tests
go test ./internal/tmux/...

# Lint (requires golangci-lint)
golangci-lint run
```

## Keybinding Convention

Keybindings are added to `~/.tmux.conf` by the user (not auto-applied). The `tmux.conf.example` shows recommended bindings:

```tmux
# Open tplm picker in a floating popup
bind-key C-p display-popup -E -w 80% -h 60% "tplm picker"
```

## Skills

**Always use the `golang-patterns` and `golang-pro` skills** when writing, reviewing, or modifying Go code in this project. Invoke them before making implementation decisions to ensure idiomatic patterns, proper concurrency handling, and best practices are followed.

## Code Conventions

### Constants — No Inline Literals

**Never use string literals, numeric literals, or format strings directly in code.** All literals must be extracted into named constants:

- **File-local constants**: If a constant is only used within a single file, declare it in a `const` block just below the imports at the top of that file.
- **Package-level constants**: If a constant is used across multiple files in the same package, place it in the package's `constants.go` file (e.g., `internal/tmux/constants.go`, `internal/cli/constants.go`, `internal/config/constants.go`, `internal/ui/constants.go`).

This applies to:
- Error message templates (e.g., `ErrFmtRenameWindow = "renaming window %q: %w"`)
- Format strings (e.g., `FmtSessionWindow = "%s:%d"`)
- User-facing messages and output strings
- Numeric values like file permissions, sizing estimates, field counts
- tmux command names, flags, and format strings (already in `tmux/constants.go`)

### Error Handling

- **Never ignore errors silently** — always handle or propagate errors from `os.UserHomeDir()`, `strconv.Atoi()`, and similar calls. If an error is intentionally ignored (e.g., cosmetic tmux focus operations), add a comment explaining why.
- **Use `fmt.Errorf` with `%w`** for error wrapping, always referencing a constant template.
- **Define sentinel errors** (`var ErrFoo = errors.New(...)`) for domain-specific error cases that callers may need to check with `errors.Is`.

### Go Idioms

- **Preallocate slices** when the size is known: use `make([]T, 0, len(source))` instead of `var s []T`.
- **Use `strings.Builder`** with `b.Grow()` for building strings with known approximate size.
- **Add doc comments** to all exported functions, types, and package-level variables.
- **Write table-driven tests** with subtests (`t.Run`) for all non-trivial logic.

## Dependencies

- `tmux` >= 3.2 (for `display-popup` support)
- Go >= 1.22
