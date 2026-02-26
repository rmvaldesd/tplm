package tmux

import "fmt"

// NewSession creates a new detached session with a name and working directory.
func NewSession(name, path string) error {
	return RunSilent(CmdNewSession, FlagDetached, FlagSession, name, FlagDir, path)
}

// KillSession kills the session with the given name.
func KillSession(name string) error {
	return RunSilent(CmdKillSession, FlagTarget, name)
}

// RenameSession renames a session.
func RenameSession(oldName, newName string) error {
	return RunSilent(CmdRenameSession, FlagTarget, oldName, newName)
}

// SwitchClient switches the current client to the given session.
func SwitchClient(name string) error {
	return RunSilent(CmdSwitchClient, FlagTarget, name)
}

// NewWindow creates a new window in the given session.
func NewWindow(session, name string) error {
	return RunSilent(CmdNewWindow, FlagTarget, session, FlagName, name)
}

// KillWindow kills a specific window. Target format: "session:windowIndex".
func KillWindow(target string) error {
	return RunSilent(CmdKillWindow, FlagTarget, target)
}

// RenameWindow renames the current window in a session.
func RenameWindow(target, name string) error {
	return RunSilent(CmdRenameWindow, FlagTarget, target, name)
}

// SendKeys sends keystrokes to a target pane.
func SendKeys(target, keys string) error {
	return RunSilent(CmdSendKeys, FlagTarget, target, keys, KeyEnter)
}

// SelectWindow selects the first window in a session.
func SelectWindow(target string) error {
	return RunSilent(CmdSelectWindow, FlagTarget, target)
}

// SelectPane selects a specific pane.
func SelectPane(target string) error {
	return RunSilent(CmdSelectPane, FlagTarget, fmt.Sprintf(FmtTargetPane0, target))
}
