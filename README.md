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

## Understanding the Layout Logic

When tplm creates a window, it follows these rules in order:

### Step 1: You always start with one full pane

The first entry in `panes:` is the **default pane** — it already exists when the window is created. It fills the entire window. You only set its `size` (used later as a reference) and optionally a `command`.

```
 ┌─────────────────────┐
 │                     │
 │     pane 0 (100%)   │
 │                     │
 └─────────────────────┘
```

### Step 2: Each new pane splits the *currently active* pane

Starting from the 2nd entry onward, each pane **splits the last pane that was created**. This is the key mental model:

- `split: horizontal` — splits **left/right** (new pane appears to the right)
- `split: vertical` — splits **top/bottom** (new pane appears below)

### Step 3: Sizes are relative to the pane being split, not the whole window

If pane 0 is 100% of the window and you split it horizontally at 40%, the new pane takes 40% of pane 0's space, leaving pane 0 with 60%.

### Putting it together: common recipes

#### Recipe A: Big pane on the LEFT, split right side

You want this:
```
 ┌──────────┬──────────┐
 │          │ terminal │
 │          │          │
 │  editor  ├──────────┤
 │  (big)   │ logs     │
 │          │          │
 └──────────┴──────────┘
```

Think step-by-step:

1. **Pane 0** — the editor, takes the full window
2. **Pane 1** — split `horizontal` (left/right) at 40% → creates the right column, editor stays on the left
3. **Pane 2** — split `vertical` (top/bottom) at 50% → splits *pane 1* (the right column) into two rows

```yaml
panes:
  - size: "60%"              # pane 0: editor (left, big)
  - split: horizontal        # pane 1: terminal (right-top)
    size: "40%"
  - split: vertical          # pane 2: logs (right-bottom, splits pane 1)
    size: "50%"
```

#### Recipe B: Big pane on the RIGHT, split left side

You want this:
```
 ┌──────────┬──────────┐
 │ terminal │          │
 │          │          │
 ├──────────┤  editor  │
 │ logs     │  (big)   │
 │          │          │
 └──────────┴──────────┘
```

The trick: you can't directly split "to the left". Instead, make the first split create the left column, then split *that* vertically, and the remaining big pane becomes the right side:

1. **Pane 0** — starts full, will become the left column
2. **Pane 1** — split `horizontal` at 60% → the new pane takes 60% on the right (this is your big editor)
3. **Pane 2** — split `vertical` at 50% → but this splits *pane 1* (the right/big one), not pane 0

The problem: pane 2 splits the last created pane (pane 1, the right side). To split pane 0 (the left side) instead, you need to reverse the order — make the big pane first and the small column second:

1. **Pane 0** — the big editor, takes the full window
2. **Pane 1** — split `horizontal` at 40% → right column appears
3. Now you have: `[editor 60% | pane1 40%]` — but editor is on the left, not the right

Since tmux always places the new split to the right (horizontal) or below (vertical), achieving "big pane on right" requires this approach:

```yaml
panes:
  - size: "40%"              # pane 0: left column (will be split)
  - split: horizontal        # pane 1: editor (right, big)
    size: "60%"
  - split: vertical          # pane 2: splits pane 1 (right)... not what we want
    size: "50%"
```

This doesn't work as expected because pane 2 splits the *last created pane* (pane 1, the right side). **The workaround**: accept that the left column pane stays unsplit, and use `on_start` or `command` fields to run what you need in each pane. Or, flip your thinking and put the big pane on the left (Recipe A) — which is the natural flow for tmux splits.

> **Rule of thumb**: tmux splits always go right or down. The simplest layouts have the big pane on the **left** or **top**, with the smaller split panes on the right or bottom side.

#### Recipe C: Three panes stacked vertically

```
 ┌─────────────────────┐
 │      editor         │
 ├─────────────────────┤
 │      terminal       │
 ├─────────────────────┤
 │      logs           │
 └─────────────────────┘
```

Each new pane splits the previous one top/bottom:

1. **Pane 0** — editor, full window
2. **Pane 1** — split `vertical` at 50% → bottom half
3. **Pane 2** — split `vertical` at 50% → splits pane 1's space in half

```yaml
panes:
  - size: "50%"              # pane 0: editor (top)
  - split: vertical          # pane 1: terminal (middle)
    size: "50%"
  - split: vertical          # pane 2: logs (bottom, splits pane 1)
    size: "50%"
```

#### Recipe D: Three panes side-by-side

```
 ┌───────┬───────┬───────┐
 │       │       │       │
 │ left  │ center│ right │
 │       │       │       │
 └───────┴───────┴───────┘
```

```yaml
panes:
  - size: "33%"              # pane 0: left
  - split: horizontal        # pane 1: center
    size: "50%"
  - split: horizontal        # pane 2: right (splits pane 1)
    size: "50%"
```

### Quick reference

| You want | First split | Second split |
|---|---|---|
| Big left + small right | `horizontal` (small %) | — |
| Big top + small bottom | `vertical` (small %) | — |
| Big left + 2 stacked right | `horizontal` | then `vertical` |
| Big top + 2 side-by-side bottom | `vertical` | then `horizontal` |
| 3 columns | `horizontal` | then `horizontal` |
| 3 rows | `vertical` | then `vertical` |

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
