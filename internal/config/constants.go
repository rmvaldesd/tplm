package config

// Config file path components.
const (
	ConfigDir  = ".config"
	ConfigApp  = "tplm"
	ConfigFile = "config.yaml"
)

// Error message templates.
const (
	ErrReadingConfig = "reading config: %w"
	ErrParsingConfig = "parsing config: %w"
)
