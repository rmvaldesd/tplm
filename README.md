# tplm — Tmux Project Layout Manager

A Go CLI tool for managing tmux sessions with predefined project layouts. Define your workspaces in YAML, browse them in a floating picker, and switch between projects instantly.

## Features

- **YAML-based config** — define projects, layouts, and startup commands
- **Floating picker** — Bubbletea TUI inside a `tmux display-popup` for browsing projects and active sessions
- **Auto session creation** — select a project and tplm creates the session with the configured layout and runs startup commands
- **Session management** — kill and rename sessions directly from the picker
- **Scriptable** — `tplm open <name>` for headless session creation

## Requirements

- tmux >= 3.2 (for `display-popup` support)
- Go >= 1.22

## Installation

```bash
git clone https://github.com/rmvaldesd/tplm.git
cd tplm
make install
```

This builds the binary and installs it to `/usr/local/bin/tplm`.

To uninstall:

```bash
make uninstall
```

## Setup

### 1. Generate a starter config

```bash
tplm init
```

This creates `~/.config/tplm/config.yaml` with example projects and layouts.

### 2. Add the tmux keybinding

Add this to your `~/.tmux.conf`:

```tmux
# Open tplm picker in a floating popup
bind-key C-p display-popup -E -w 80% -h 60% "tplm picker"
```

Reload tmux config:

```
tmux source-file ~/.tmux.conf
```

### 3. Edit the config

Edit `~/.config/tplm/config.yaml` to define your projects:

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

  - name: frontend
    path: ~/Projects/frontend
    layout: fullstack

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

  fullstack:
    windows:
      - name: frontend
        panes:
          - size: "50%"
          - split: horizontal
            size: "50%"
      - name: backend
        panes:
          - size: "60%"
          - split: horizontal
            size: "40%"
```

## Usage

### Picker (interactive)

Press `prefix + C-p` inside tmux to open the floating picker:

```
┌───────────────────────────────────────┐
│  tplm                  d:kill r:rename│
│ ───────────────────────────────────── │
│  Projects                             │
│ ───────────────────────────────────── │
│  > my-api            ~/Projects/api   │
│    frontend          ~/Projects/fe    │
│    infra             ~/Projects/infra │
│ ───────────────────────────────────── │
│  Active Sessions                      │
│ ───────────────────────────────────── │
│  ● my-api      3 windows             │
│  ● frontend    2 windows             │
│                                       │
│  ↑↓ navigate  ⏎ open/switch  q quit  │
└───────────────────────────────────────┘
```

### Keybindings

| Key | Context | Action |
|---|---|---|
| `j` / `k` / arrows | Anywhere | Navigate up/down |
| `Enter` | On a project | Create session from layout (if needed) and switch to it |
| `Enter` | On a session | Switch to it |
| `d` | On a session | Kill session (with `y/n` confirmation) |
| `r` | On a session | Rename session inline |
| `q` / `Esc` | Anywhere | Close picker |

### CLI Commands

```bash
# Open the interactive picker (usually called via tmux keybinding)
tplm picker

# Create a session from config and switch to it (no TUI)
tplm open my-api

# List projects and active sessions
tplm list

# Generate starter config
tplm init

# Use a custom config path
tplm --config /path/to/config.yaml list
```

## Config Reference

### Projects

| Field | Required | Description |
|---|---|---|
| `name` | yes | Project name (used as tmux session name) |
| `path` | yes | Working directory (`~` is expanded) |
| `layout` | no | Name of a layout defined in `layouts` |
| `on_start` | no | Commands to run in specific windows on session creation |

### Layouts

Each layout defines a list of windows. Each window has a name and a list of panes.

| Field | Required | Description |
|---|---|---|
| `name` | yes | Window name |
| `panes` | no | List of pane splits (first pane is the default, additional panes split from it) |

### Panes

| Field | Required | Description |
|---|---|---|
| `split` | no | `horizontal` (side-by-side) or `vertical` (top/bottom) |
| `size` | no | Percentage of the split, e.g. `"30%"` |

## Session Creation Flow

When you select a project that has no active session:

1. Creates a detached tmux session at the project path
2. Sets up windows and pane splits from the layout config
3. Runs `on_start` commands in the specified windows
4. Switches your client to the new session

If the session already exists, it simply switches to it.

## Development

```bash
make build    # compile binary
make test     # run tests
make lint     # run golangci-lint
make clean    # remove binary
```

## License

MIT
