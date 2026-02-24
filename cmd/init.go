package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/rmvaldesd/tplm/internal/config"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Generate a starter config file",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := cfgPath

		// Create parent directory.
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("creating config directory: %w", err)
		}

		// Don't overwrite existing config.
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("config already exists at %s", path)
		}

		if err := os.WriteFile(path, []byte(config.ExampleConfig()), 0o644); err != nil {
			return fmt.Errorf("writing config: %w", err)
		}

		fmt.Printf("Created starter config at %s\n", path)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
