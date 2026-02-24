package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/rmvaldesd/tplm/internal/config"
)

var (
	cfgPath string
	cfg     *config.Config
)

var rootCmd = &cobra.Command{
	Use:   "tplm",
	Short: "tmux Project Layout Manager",
	Long:  "A tool for managing tmux sessions with predefined project layouts.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip config loading for init command.
		if cmd.Name() == "init" {
			return nil
		}

		var err error
		cfg, err = config.Load(cfgPath)
		if err != nil {
			return fmt.Errorf("loading config: %w\nRun 'tplm init' to create a starter config", err)
		}
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgPath, "config", config.DefaultConfigPath(), "path to config file")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
