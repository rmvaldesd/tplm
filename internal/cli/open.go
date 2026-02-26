package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/rmvaldesd/tplm/internal/config"
	"github.com/rmvaldesd/tplm/internal/tmux"
)

var openCmd = &cobra.Command{
	Use:   "open <project-name>",
	Short: "Create a session from project config and switch to it",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		proj := cfg.FindProject(name)
		if proj == nil {
			return fmt.Errorf("project %q not found in config", name)
		}

		return OpenProject(proj)
	},
}

func init() {
	rootCmd.AddCommand(openCmd)
}

// OpenProject creates a tmux session for the project (if needed) and switches to it.
func OpenProject(proj *config.Project) error {
	if tmux.SessionExists(proj.Name) {
		return tmux.SwitchClient(proj.Name)
	}

	if err := tmux.NewSession(proj.Name, proj.Path); err != nil {
		return fmt.Errorf("creating session: %w", err)
	}

	layout := cfg.GetLayout(proj)

	if err := tmux.ApplyLayout(proj.Name, layout, proj.Path); err != nil {
		return fmt.Errorf("applying layout: %w", err)
	}

	if len(proj.OnStart) > 0 {
		if err := tmux.RunOnStart(proj.Name, layout, proj.OnStart); err != nil {
			return fmt.Errorf("running on_start: %w", err)
		}
	}

	return tmux.SwitchClient(proj.Name)
}
