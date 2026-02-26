package ui

// Color palette (lipgloss ANSI color codes).
const (
	ColorAccent    = "170"
	ColorMuted     = "241"
	ColorDim       = "238"
	ColorText      = "252"
	ColorGreen     = "42"
	ColorRed       = "196"
)

// UI symbols.
const (
	SymbolActive       = "●"
	SymbolChevronRight = "▶"
	SymbolChevronDown  = "▼"
	SymbolSeparator    = "─"
	SymbolCursor       = "> "
	SymbolNoCursor     = "  "
)

// Section headers and title.
const (
	TitleText          = "tplm"
	HeaderProjects     = "Projects"
	HeaderSessions     = "Active Sessions"
	HintText           = "d:kill  r:rename"
)

// Messages.
const (
	MsgNoProjects  = "(no projects configured)"
	MsgNoSessions  = "(no active sessions)"
	MsgHelpBar     = "hjkl navigate  ⏎ select  d kill  r rename  q quit"
	MsgConfirmKill = "  Kill %s %q? (y/n)"
	MsgError       = "  Error: %v"
	MsgSession     = "session"
	MsgWindow      = "window"
)

// Rename input settings.
const (
	RenameCharLimit = 64
	RenameWidth     = 40
	RenamePrompt    = "Rename: "
)

// Key names for rename input handling.
const (
	KeyEnter = "enter"
	KeyEsc   = "esc"
)

// Default picker width when terminal size is unknown.
const DefaultPickerWidth = 60
