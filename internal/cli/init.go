package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/rmvaldesd/tplm/internal/config"
)

var initCmd = &cobra.Command{
	Use:   InitUse,
	Short: InitShort,
	RunE: func(cmd *cobra.Command, args []string) error {
		path := cfgPath

		// Create parent directory.
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf(ErrCreatingDir, err)
		}

		// Don't overwrite existing config.
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf(ErrConfigExists, path)
		}

		if err := os.WriteFile(path, []byte(config.ExampleConfig()), 0o644); err != nil {
			return fmt.Errorf(ErrWritingConfig, err)
		}

		fmt.Printf(OutputCreatedConfig, path)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
