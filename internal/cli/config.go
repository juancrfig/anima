// internal/cli/config.go
package cli

import (
	"fmt"
	"path/filepath"

	"anima/internal/config"
	"github.com/spf13/cobra"
)

// ConfigCmd returns the root command for config-related operations.
func ConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage Anima configuration.",
		Long:  `Set or get configuration values for the Anima CLI.`,
	}

	// Add subcommands to the parent 'config' command
	cmd.AddCommand(configSetCmd())

	return cmd
}

// configSetCmd returns the command for 'config set'.
func configSetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set [key] [value]",
		Short: "Set a configuration key-value pair.",
		Long:  `Sets a configuration value. For example: anima config set location "Cucuta, Colombia"`,
		// Enforce that we get exactly two arguments.
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			value := args[1]

			animaPath, err := GetAnimaPath()
			if err != nil {
				return err
			}
			configPath := filepath.Join(animaPath, "config.json")

			cfg, err := config.New(configPath)
			if err != nil {
				return fmt.Errorf("could not initialize config: %w", err)
			}

			if err := cfg.Set(key, value); err != nil {
				return fmt.Errorf("could not set config value: %w", err)
			}

			return nil
		},
	}
	return cmd
}
