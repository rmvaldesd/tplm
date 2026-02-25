package tmux

import "fmt"

// NewSession creates a new detached session with a name and working directory.
func NewSession(name, path string) error {
	return RunSilent("new-session", "-d", "-s", name, "-c", path)
}

// KillSession kills the session with the given name.
func KillSession(name string) error {
	return RunSilent("kill-session", "-t", name)
}

// RenameSession renames a session.
func RenameSession(oldName, newName string) error {
	return RunSilent("rename-session", "-t", oldName, newName)
}

// SwitchClient switches the current client to the given session.
func SwitchClient(name string) error {
	return RunSilent("switch-client", "-t", name)
}

// NewWindow creates a new window in the given session.
func NewWindow(session, name string) error {
	return RunSilent("new-window", "-t", session, "-n", name)
}

// KillWindow kills a specific window. Target format: "session:windowIndex".
func KillWindow(target string) error {
	return RunSilent("kill-window", "-t", target)
}

// RenameWindow renames the current window in a session.
func RenameWindow(target, name string) error {
	return RunSilent("rename-window", "-t", target, name)
}

// SendKeys sends keystrokes to a target pane.
func SendKeys(target, keys string) error {
	return RunSilent("send-keys", "-t", target, keys, "Enter")
}

// SelectWindow selects the first window in a session.
func SelectWindow(target string) error {
	return RunSilent("select-window", "-t", target)
}

// SelectPane selects a specific pane.
func SelectPane(target string) error {
	return RunSilent("select-pane", "-t", fmt.Sprintf("%s.0", target))
}
