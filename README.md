# tplm — Tmux Project Layout Manager

A Go CLI tool for managing tmux sessions with predefined project layouts. Define your workspaces in YAML, browse them in a floating picker, and switch between projects instantly.

## Features

- **YAML-based config** — define projects, layouts, and startup commands
- **Floating picker** — Bubbletea TUI inside a `tmux display-popup` for browsing projects and active sessions
- **Auto session creation** — select a project and tplm creates the session with the configured layout and runs startup commands
- **Session management** — kill and rename sessions, close individual windows directly from the picker
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
bind-key l display-popup -E -w 80% -h 60% "tplm picker"
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

Press `prefix + l` inside tmux to open the floating picker:

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
│  hjkl navigate  ⏎ select  q quit     │
└───────────────────────────────────────┘
```

### Keybindings

| Key | Context | Action |
|---|---|---|
| `j` / `k` / arrows | Anywhere | Navigate up/down |
| `l` / `Enter` | On a collapsed session | Expand to show windows |
| `l` / `Enter` | On an expanded session | Move to first window / toggle collapse |
| `h` | On an expanded session | Collapse windows |
| `h` | On a window | Jump to parent session |
| `Enter` | On a window | Switch to that window |
| `Enter` | On a project | Create session (if needed) and switch |
| `d` | On a session | Kill session (with `y/n` confirmation) |
| `d` | On a window | Kill window (with `y/n` confirmation) |
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
| `command` | no | Command to run in this pane on session creation |

## Layout Examples

Layouts define how windows are split into panes. Each window starts with one pane; additional panes are created by splitting from it. `horizontal` splits side-by-side (left/right), `vertical` splits top/bottom.

### Simple: single window, 2 panes side-by-side

```
 ┌──────────┬──────────┐
 │          │          │
 │  editor  │  terminal│
 │  70%     │  30%     │
 │          │          │
 └──────────┴──────────┘
```

```yaml
layouts:
  simple:
    windows:
      - name: code
        panes:
          - size: "70%"
          - split: horizontal
            size: "30%"
```

### Two panes stacked vertically

```
 ┌─────────────────────┐
 │       editor        │
 │       70%           │
 ├─────────────────────┤
 │       terminal      │
 │       30%           │
 └─────────────────────┘
```

```yaml
layouts:
  stacked:
    windows:
      - name: code
        panes:
          - size: "70%"
          - split: vertical
            size: "30%"
```

### Two vertical panes, right side split horizontally

A main editor on the left, with a terminal and a server log stacked on the right.

```
 ┌──────────┬──────────┐
 │          │ terminal │
 │          │ 50%      │
 │  editor  ├──────────┤
 │  60%     │ logs     │
 │          │ 50%      │
 └──────────┴──────────┘
```

```yaml
layouts:
  ide:
    windows:
      - name: workspace
        panes:
          - size: "60%"
          - split: horizontal
            size: "40%"
          - split: vertical
            size: "50%"
```

The first pane takes 60%. The second pane splits horizontally (side-by-side), taking 40% of the window width. The third pane splits vertically (top/bottom) from the second pane, each getting 50% of that column.

### Complex: two vertical columns, each split into two rows with different sizes

```
 ┌──────────┬──────────┐
 │  editor  │ server   │
 │  70%     │ 60%      │
 ├──────────┼──────────┤
 │ terminal │ logs     │
 │  30%     │ 40%      │
 └──────────┴──────────┘
```

```yaml
layouts:
  quad:
    windows:
      - name: dev
        panes:
          - size: "50%"
          - split: horizontal
            size: "50%"
          - split: vertical
            size: "40%"
```

> **Note:** tmux splits are relative to the pane being split, not the whole window. The third pane (split vertical at 40%) splits the right column into 60%/40% top/bottom. To also split the left column, you would use a second window or run `on_start` commands.

### Multi-window layout

Combine multiple windows for a full project workspace:

```yaml
layouts:
  fullstack:
    windows:
      - name: editor
        panes:
          - size: "70%"
            command: "nvim ."
          - split: horizontal
            size: "30%"
      - name: servers
        panes:
          - size: "50%"
            command: "go run ./cmd/api"
          - split: horizontal
            size: "50%"
            command: "npm run dev"
      - name: logs
        panes:
          - size: "100%"
            command: "tail -f /var/log/app.log"

projects:
  - name: my-app
    path: ~/Projects/my-app
    layout: fullstack
```

This creates three windows: `editor` (70/30 split), `servers` (50/50 split), and `logs` (single pane). Each pane can optionally specify a `command` to run on creation — this is an alternative to using `on_start`, and gives you per-pane control rather than per-window.

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
