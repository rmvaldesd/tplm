package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/rmvaldesd/tplm/internal/config"
)

// Package-level state for Cobra command closures. This is acceptable for a CLI
// tool where commands run sequentially, but would be problematic in a library.
var (
	cfgPath string
	cfg     *config.Config
)

var rootCmd = &cobra.Command{
	Use:   RootUse,
	Short: RootShort,
	Long:  RootLong,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip config loading for init command.
		if cmd.Name() == CmdInit {
			return nil
		}

		var err error
		cfg, err = config.Load(cfgPath)
		if err != nil {
			return fmt.Errorf(ErrLoadingConfig, err)
		}
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgPath, FlagConfig, config.DefaultConfigPath(), FlagConfigDesc)
}

// Execute runs the root Cobra command. It exits with code 1 on error.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
