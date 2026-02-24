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

// SessionExists checks if a session with the given name exists.
func SessionExists(name string) bool {
	err := RunSilent("has-session", "-t", name)
	return err == nil
}
