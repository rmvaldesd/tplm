package tmux

import (
	"strconv"
	"strings"
)

// SessionInfo holds metadata about a tmux session.
type SessionInfo struct {
	Name       string
	Windows    int
	Attached   bool
	Path       string
}

// ListSessions returns all active tmux sessions.
func ListSessions() ([]SessionInfo, error) {
	out, err := Run("list-sessions", "-F", "#{session_name}\t#{session_windows}\t#{session_attached}\t#{session_path}")
	if err != nil {
		// No server running means no sessions.
		if strings.Contains(err.Error(), "no server") || strings.Contains(err.Error(), "no current") {
			return nil, nil
		}
		return nil, err
	}

	if out == "" {
		return nil, nil
	}

	var sessions []SessionInfo
	for _, line := range strings.Split(out, "\n") {
		parts := strings.SplitN(line, "\t", 4)
		if len(parts) < 4 {
			continue
		}
		wins, _ := strconv.Atoi(parts[1])
		attached := parts[2] == "1"
		sessions = append(sessions, SessionInfo{
			Name:     parts[0],
			Windows:  wins,
			Attached: attached,
			Path:     parts[3],
		})
	}
	return sessions, nil
}

// WindowInfo holds metadata about a tmux window.
type WindowInfo struct {
	Index  int
	Name   string
	Active bool
}

// ListWindows returns all windows for the given session.
func ListWindows(session string) ([]WindowInfo, error) {
	out, err := Run("list-windows", "-t", session, "-F", "#{window_index}\t#{window_name}\t#{window_active}")
	if err != nil {
		return nil, err
	}

	if out == "" {
		return nil, nil
	}

	var windows []WindowInfo
	for _, line := range strings.Split(out, "\n") {
		parts := strings.SplitN(line, "\t", 3)
		if len(parts) < 3 {
			continue
		}
		idx, _ := strconv.Atoi(parts[0])
		active := parts[2] == "1"
		windows = append(windows, WindowInfo{
			Index:  idx,
			Name:   parts[1],
			Active: active,
		})
	}
	return windows, nil
}

// SessionExists checks if a session with the given name exists.
func SessionExists(name string) bool {
	err := RunSilent("has-session", "-t", name)
	return err == nil
}
