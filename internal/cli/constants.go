package cli

// Cobra command metadata.
const (
	RootUse   = "tplm"
	RootShort = "tmux Project Layout Manager"
	RootLong  = "A tool for managing tmux sessions with predefined project layouts."

	PickerUse   = "picker"
	PickerShort = "Open the interactive project picker"
	PickerLong  = "Opens a Bubbletea TUI for browsing and switching between projects and sessions.\nIntended to run inside tmux display-popup."

	OpenUse   = "open <project-name>"
	OpenShort = "Create a session from project config and switch to it"

	ListUse   = "list"
	ListShort = "Print projects and active tmux sessions"

	InitUse   = "init"
	InitShort = "Generate a starter config file"
)

// Flag names.
const FlagConfig = "config"

// Flag descriptions.
const FlagConfigDesc = "path to config file"

// Command names used for skipping config load.
const CmdInit = "init"

// Error message templates.
const (
	ErrLoadingConfig    = "loading config: %w\nRun 'tplm init' to create a starter config"
	ErrRunningPicker    = "running picker: %w"
	ErrProjectNotFound  = "project %q not found in config"
	ErrCreatingSession  = "creating session: %w"
	ErrApplyingLayout   = "applying layout: %w"
	ErrRunningOnStart   = "running on_start: %w"
	ErrCreatingDir      = "creating config directory: %w"
	ErrConfigExists     = "config already exists at %s"
	ErrWritingConfig    = "writing config: %w"
)

// User-facing output strings.
const (
	OutputProjects       = "Projects:"
	OutputActiveSessions = "Active Sessions:"
	OutputNone           = "  (none)"
	OutputCreatedConfig  = "Created starter config at %s\n"
	OutputAttached       = "*"
	OutputNotAttached    = " "
)
