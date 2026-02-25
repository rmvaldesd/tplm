package config

// Config is the top-level YAML configuration.
type Config struct {
	Projects []Project         `yaml:"projects"`
	Layouts  map[string]Layout `yaml:"layouts"`
}

// Project defines a workspace entry.
type Project struct {
	Name    string    `yaml:"name"`
	Path    string    `yaml:"path"`
	Layout  string    `yaml:"layout"`
	OnStart []OnStart `yaml:"on_start,omitempty"`
}

// OnStart defines a command to run in a specific window on session creation.
type OnStart struct {
	Window  string `yaml:"window"`
	Command string `yaml:"command"`
}

// Layout defines a set of windows and their pane splits.
type Layout struct {
	Windows []Window `yaml:"windows"`
}

// Window defines a named window with pane splits.
type Window struct {
	Name  string `yaml:"name"`
	Panes []Pane `yaml:"panes"`
}

// Pane defines a single pane with optional split direction and size.
type Pane struct {
	Split   string `yaml:"split,omitempty"`   // "horizontal" or "vertical"
	Size    string `yaml:"size,omitempty"`     // e.g. "70%"
	Command string `yaml:"command,omitempty"` // optional command to run on pane startup
}
