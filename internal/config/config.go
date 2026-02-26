package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// DefaultConfigPath returns ~/.config/tplm/config.yaml.
func DefaultConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ConfigDir, ConfigApp, ConfigFile)
}

// Load reads and parses the YAML config file at the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf(ErrReadingConfig, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf(ErrParsingConfig, err)
	}

	// Resolve ~ in project paths.
	home, _ := os.UserHomeDir()
	for i := range cfg.Projects {
		cfg.Projects[i].Path = expandHome(cfg.Projects[i].Path, home)
	}

	return &cfg, nil
}

// FindProject returns the project with the given name, or nil.
func (c *Config) FindProject(name string) *Project {
	for i := range c.Projects {
		if c.Projects[i].Name == name {
			return &c.Projects[i]
		}
	}
	return nil
}

// GetLayout returns the layout for a project, falling back to a single-window default.
func (c *Config) GetLayout(proj *Project) Layout {
	if proj.Layout != "" {
		if l, ok := c.Layouts[proj.Layout]; ok {
			return l
		}
	}
	// Default: single window named "main" with one pane.
	return Layout{
		Windows: []Window{{Name: "main"}},
	}
}

func expandHome(path, home string) string {
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(home, path[2:])
	}
	return path
}

// ExampleConfig returns a starter YAML config string.
func ExampleConfig() string {
	return `# tplm configuration
# Place this file at ~/.config/tplm/config.yaml

projects:
  - name: my-api
    path: ~/Projects/my-api
    layout: dev
    on_start:
      - window: editor
        command: nvim .
      - window: server
        command: "echo 'start your server here'"

  - name: frontend
    path: ~/Projects/frontend
    layout: fullstack

layouts:
  dev:
    windows:
      - name: editor
        panes:
          - size: "70%"
          - split: horizontal
            size: "30%"
      - name: server
        panes:
          - size: "100%"

  fullstack:
    windows:
      - name: frontend
        panes:
          - size: "50%"
          - split: horizontal
            size: "50%"
      - name: backend
        panes:
          - size: "60%"
          - split: horizontal
            size: "40%"
`
}
