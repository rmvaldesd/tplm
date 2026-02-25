package tmux

import (
	"fmt"
	"strings"

	"github.com/rmvaldesd/tplm/internal/config"
)

// ApplyLayout creates windows and splits panes according to the layout config.
func ApplyLayout(sessionName string, layout config.Layout, projectPath string) error {
	for i, win := range layout.Windows {
		target := fmt.Sprintf("%s:%d", sessionName, i)

		if i == 0 {
			// First window is created with the session; just rename it.
			if err := RenameWindow(target, win.Name); err != nil {
				return fmt.Errorf("renaming window %q: %w", win.Name, err)
			}
		} else {
			if err := NewWindow(sessionName, win.Name); err != nil {
				return fmt.Errorf("creating window %q: %w", win.Name, err)
			}
			// Set the working directory for the new window.
			if err := SendKeys(target, fmt.Sprintf("cd %s", shellEscape(projectPath))); err != nil {
				return fmt.Errorf("setting directory for window %q: %w", win.Name, err)
			}
		}

		// Run command in the first pane if specified.
		if len(win.Panes) > 0 && win.Panes[0].Command != "" {
			paneTarget := fmt.Sprintf("%s.0", target)
			if err := SendKeys(paneTarget, win.Panes[0].Command); err != nil {
				return fmt.Errorf("running command in pane 0 of window %q: %w", win.Name, err)
			}
		}

		// Split panes (skip the first pane — it exists by default).
		for j := 1; j < len(win.Panes); j++ {
			pane := win.Panes[j]
			args := []string{"split-window", "-t", target}

			// Default to horizontal split (side-by-side).
			if pane.Split == "vertical" {
				args = append(args, "-v")
			} else {
				args = append(args, "-h")
			}

			if pane.Size != "" {
				pct := strings.TrimSuffix(pane.Size, "%")
				args = append(args, "-p", pct)
			}

			args = append(args, "-c", projectPath)

			if err := RunSilent(args...); err != nil {
				return fmt.Errorf("splitting pane %d in window %q: %w", j, win.Name, err)
			}

			// Run command in this pane if specified.
			if pane.Command != "" {
				paneTarget := fmt.Sprintf("%s.%d", target, j)
				if err := SendKeys(paneTarget, pane.Command); err != nil {
					return fmt.Errorf("running command in pane %d of window %q: %w", j, win.Name, err)
				}
			}
		}

		// Select the first pane after all splits.
		_ = SelectPane(target)
	}

	// Select the first window.
	_ = SelectWindow(fmt.Sprintf("%s:0", sessionName))
	return nil
}

// RunOnStart sends the on_start commands to the appropriate windows.
func RunOnStart(sessionName string, layout config.Layout, commands []config.OnStart) error {
	// Build a map of window name -> index.
	winIndex := make(map[string]int)
	for i, w := range layout.Windows {
		winIndex[w.Name] = i
	}

	for _, cmd := range commands {
		idx, ok := winIndex[cmd.Window]
		if !ok {
			continue // Skip if window not found in layout.
		}
		target := fmt.Sprintf("%s:%d.0", sessionName, idx)
		if err := SendKeys(target, cmd.Command); err != nil {
			return fmt.Errorf("running on_start for window %q: %w", cmd.Window, err)
		}
	}
	return nil
}

func shellEscape(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}
