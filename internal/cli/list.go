package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/rmvaldesd/tplm/internal/tmux"
)

var listCmd = &cobra.Command{
	Use:   ListUse,
	Short: ListShort,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(OutputProjects)
		for _, p := range cfg.Projects {
			fmt.Printf("  %-20s %s\n", p.Name, p.Path)
		}

		sessions, err := tmux.ListSessions()
		if err != nil {
			return err
		}

		fmt.Println()
		fmt.Println(OutputActiveSessions)
		if len(sessions) == 0 {
			fmt.Println(OutputNone)
		}
		for _, s := range sessions {
			attached := OutputNotAttached
			if s.Attached {
				attached = OutputAttached
			}
			fmt.Printf("  %s %-20s %d windows\n", attached, s.Name, s.Windows)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
