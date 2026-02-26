package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/rmvaldesd/tplm/internal/config"
)

const (
	dirPermissions  = 0o755
	filePermissions = 0o644
)

var initCmd = &cobra.Command{
	Use:   InitUse,
	Short: InitShort,
	RunE: func(cmd *cobra.Command, args []string) error {
		path := cfgPath

		// Create parent directory.
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, dirPermissions); err != nil {
			return fmt.Errorf(ErrCreatingDir, err)
		}

		// Don't overwrite existing config.
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf(ErrConfigExists, path)
		}

		if err := os.WriteFile(path, []byte(config.ExampleConfig()), filePermissions); err != nil {
			return fmt.Errorf(ErrWritingConfig, err)
		}

		fmt.Printf(OutputCreatedConfig, path)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
