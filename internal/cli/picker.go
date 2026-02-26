package cli

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/rmvaldesd/tplm/internal/ui"
)

var pickerCmd = &cobra.Command{
	Use:   "picker",
	Short: "Open the interactive project picker",
	Long:  "Opens a Bubbletea TUI for browsing and switching between projects and sessions.\nIntended to run inside tmux display-popup.",
	RunE: func(cmd *cobra.Command, args []string) error {
		m := ui.NewPicker(cfg)
		p := tea.NewProgram(m, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			return fmt.Errorf("running picker: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pickerCmd)
}
