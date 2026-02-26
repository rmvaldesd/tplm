package tmux

import (
	"fmt"
	"strings"

	"github.com/rmvaldesd/tplm/internal/config"
)

// ApplyLayout creates windows and splits panes according to the layout config.
func ApplyLayout(sessionName string, layout config.Layout, projectPath string) error {
	for i, win := range layout.Windows {
		target := fmt.Sprintf(FmtSessionWindow, sessionName, i)

		if i == 0 {
			// First window is created with the session; just rename it.
			if err := RenameWindow(target, win.Name); err != nil {
				return fmt.Errorf(ErrFmtRenameWindow, win.Name, err)
			}
		} else {
			if err := NewWindow(sessionName, win.Name); err != nil {
				return fmt.Errorf(ErrFmtCreateWindow, win.Name, err)
			}
			// Set the working directory for the new window.
			if err := SendKeys(target, fmt.Sprintf(FmtCdCommand, shellEscape(projectPath))); err != nil {
				return fmt.Errorf(ErrFmtSetDir, win.Name, err)
			}
		}

		// Run command in the first pane if specified.
		if len(win.Panes) > 0 && win.Panes[0].Command != "" {
			paneTarget := fmt.Sprintf(FmtTargetPane0, target)
			if err := SendKeys(paneTarget, win.Panes[0].Command); err != nil {
				return fmt.Errorf(ErrFmtRunPaneCmd, 0, win.Name, err)
			}
		}

		// Split panes (skip the first pane — it exists by default).
		for j := 1; j < len(win.Panes); j++ {
			pane := win.Panes[j]
			args := []string{CmdSplitWindow, FlagTarget, target}

			// Default to horizontal split (side-by-side).
			if pane.Split == SplitVertical {
				args = append(args, FlagVertical)
			} else {
				args = append(args, FlagHoriz)
			}

			if pane.Size != "" {
				pct := strings.TrimSuffix(pane.Size, SizeSuffix)
				args = append(args, FlagPrint, pct)
			}

			args = append(args, FlagDir, projectPath)

			if err := RunSilent(args...); err != nil {
				return fmt.Errorf(ErrFmtSplitPane, j, win.Name, err)
			}

			// Run command in this pane if specified.
			if pane.Command != "" {
				paneTarget := fmt.Sprintf(FmtTargetPaneN, target, j)
				if err := SendKeys(paneTarget, pane.Command); err != nil {
					return fmt.Errorf(ErrFmtRunPaneCmd, j, win.Name, err)
				}
			}
		}

		// Select the first pane after all splits. Safe to ignore: cosmetic
		// focus operation — the layout is already applied at this point.
		_ = SelectPane(target)
	}

	// Select the first window. Safe to ignore: cosmetic focus operation
	// — all windows and panes are already created.
	_ = SelectWindow(fmt.Sprintf(FmtSessionFirst, sessionName))
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
		target := fmt.Sprintf(FmtSessionWindowPane, sessionName, idx)
		if err := SendKeys(target, cmd.Command); err != nil {
			return fmt.Errorf(ErrFmtRunOnStart, cmd.Window, err)
		}
	}
	return nil
}

func shellEscape(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}
