package tmux

// tmux binary.
const TmuxBin = "tmux"

// Subcommand names.
const (
	CmdNewSession     = "new-session"
	CmdKillSession    = "kill-session"
	CmdRenameSession  = "rename-session"
	CmdSwitchClient   = "switch-client"
	CmdNewWindow      = "new-window"
	CmdKillWindow     = "kill-window"
	CmdRenameWindow   = "rename-window"
	CmdSendKeys       = "send-keys"
	CmdSelectWindow   = "select-window"
	CmdSelectPane     = "select-pane"
	CmdSplitWindow    = "split-window"
	CmdListSessions   = "list-sessions"
	CmdListWindows    = "list-windows"
	CmdDisplayMessage = "display-message"
	CmdHasSession     = "has-session"
)

// Flags.
const (
	FlagDetached  = "-d"
	FlagSession   = "-s"
	FlagDir       = "-c"
	FlagTarget    = "-t"
	FlagName      = "-n"
	FlagFormat    = "-F"
	FlagPrint    = "-p"
	FlagVertical = "-v"
	FlagHoriz    = "-h"
)

// Format strings for tmux queries.
const (
	SessionListFormat = "#{session_name}\t#{session_windows}\t#{session_attached}\t#{session_path}"
	WindowListFormat  = "#{window_index}\t#{window_name}\t#{window_active}"
	SessionNameFormat = "#{session_name}"
)

// Error substrings used to detect expected failure modes.
const (
	ErrNoServer  = "no server"
	ErrNoCurrent = "no current"
)

// SendKeys constants.
const KeyEnter = "Enter"

// Split direction values.
const SplitVertical = "vertical"
