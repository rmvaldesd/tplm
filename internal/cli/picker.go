package cli

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/rmvaldesd/tplm/internal/ui"
)

var pickerCmd = &cobra.Command{
	Use:   PickerUse,
	Short: PickerShort,
	Long:  PickerLong,
	RunE: func(cmd *cobra.Command, args []string) error {
		m := ui.NewPicker(cfg)
		p := tea.NewProgram(m, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			return fmt.Errorf(ErrRunningPicker, err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pickerCmd)
}
