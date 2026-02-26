package tmux

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ErrTmuxNoServer indicates the tmux server is not running.
var ErrTmuxNoServer = errors.New("tmux: no server running")

// SessionInfo holds metadata about a tmux session.
type SessionInfo struct {
	Name     string
	Windows  int
	Attached bool
	Path     string
}

// ListSessions returns all active tmux sessions.
func ListSessions() ([]SessionInfo, error) {
	out, err := Run(CmdListSessions, FlagFormat, SessionListFormat)
	if err != nil {
		// No server running means no sessions.
		if strings.Contains(err.Error(), ErrNoServer) || strings.Contains(err.Error(), ErrNoCurrent) {
			return nil, nil
		}
		return nil, err
	}

	if out == "" {
		return nil, nil
	}

	lines := strings.Split(out, "\n")
	sessions := make([]SessionInfo, 0, len(lines))
	for _, line := range lines {
		parts := strings.SplitN(line, "\t", 4)
		if len(parts) < 4 {
			continue
		}
		wins, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("parsing window count for session %q: %w", parts[0], err)
		}
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
	out, err := Run(CmdListWindows, FlagTarget, session, FlagFormat, WindowListFormat)
	if err != nil {
		return nil, err
	}

	if out == "" {
		return nil, nil
	}

	lines := strings.Split(out, "\n")
	windows := make([]WindowInfo, 0, len(lines))
	for _, line := range lines {
		parts := strings.SplitN(line, "\t", 3)
		if len(parts) < 3 {
			continue
		}
		idx, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("parsing window index %q: %w", parts[0], err)
		}
		active := parts[2] == "1"
		windows = append(windows, WindowInfo{
			Index:  idx,
			Name:   parts[1],
			Active: active,
		})
	}
	return windows, nil
}

// CurrentSession returns the name of the session the current client is attached to.
func CurrentSession() (string, error) {
	out, err := Run(CmdDisplayMessage, FlagPrint, SessionNameFormat)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// NeighborSession returns the name of an adjacent session to switch to
// before killing the current one. It returns the next session in the list,
// or the previous one if current is last. If current is the only session,
// it returns "", false.
func NeighborSession(current string) (string, bool) {
	sessions, err := ListSessions()
	if err != nil || len(sessions) <= 1 {
		return "", false
	}

	for i, s := range sessions {
		if s.Name == current {
			if i+1 < len(sessions) {
				return sessions[i+1].Name, true
			}
			return sessions[i-1].Name, true
		}
	}
	return "", false
}

// SessionExists checks if a session with the given name exists.
func SessionExists(name string) bool {
	err := RunSilent(CmdHasSession, FlagTarget, name)
	return err == nil
}
