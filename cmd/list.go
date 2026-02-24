package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/rmvaldesd/tplm/internal/tmux"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Print projects and active tmux sessions",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Projects:")
		for _, p := range cfg.Projects {
			fmt.Printf("  %-20s %s\n", p.Name, p.Path)
		}

		sessions, err := tmux.ListSessions()
		if err != nil {
			return err
		}

		fmt.Println()
		fmt.Println("Active Sessions:")
		if len(sessions) == 0 {
			fmt.Println("  (none)")
		}
		for _, s := range sessions {
			attached := " "
			if s.Attached {
				attached = "*"
			}
			fmt.Printf("  %s %-20s %d windows\n", attached, s.Name, s.Windows)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
