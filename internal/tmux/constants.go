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
	FlagDetached = "-d"
	FlagSession  = "-s"
	FlagDir      = "-c"
	FlagTarget   = "-t"
	FlagName     = "-n"
	FlagFormat   = "-F"
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

// Target format strings used to build tmux target specifiers.
const (
	FmtSessionWindow     = "%s:%d"   // session:window
	FmtTargetPane0       = "%s.0"    // target.pane0
	FmtTargetPaneN       = "%s.%d"   // target.paneN
	FmtSessionWindowPane = "%s:%d.0" // session:window.pane0
	FmtSessionFirst      = "%s:0"    // session:firstWindow
)

// Error substrings used to detect expected failure modes.
const (
	ErrNoServer  = "no server"
	ErrNoCurrent = "no current"
)

// Error message templates.
const (
	ErrFmtRun              = "tmux %s: %s (%w)"
	ErrFmtRenameWindow     = "renaming window %q: %w"
	ErrFmtCreateWindow     = "creating window %q: %w"
	ErrFmtSetDir           = "setting directory for window %q: %w"
	ErrFmtRunPaneCmd       = "running command in pane %d of window %q: %w"
	ErrFmtSplitPane        = "splitting pane %d in window %q: %w"
	ErrFmtRunOnStart       = "running on_start for window %q: %w"
	ErrFmtParseWinCount    = "parsing window count for session %q: %w"
	ErrFmtParseWinIndex    = "parsing window index %q: %w"
)

// Shell command templates.
const FmtCdCommand = "cd %s"

// Size suffix stripped when passing percentage to tmux.
const SizeSuffix = "%"

// Parsing constants.
const (
	SessionFieldCount = 4
	WindowFieldCount  = 3
	AttachedValue     = "1"
	ActiveValue       = "1"
)

// SendKeys constants.
const KeyEnter = "Enter"

// Split direction values.
const SplitVertical = "vertical"
