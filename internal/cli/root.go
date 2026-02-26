package cli

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

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
