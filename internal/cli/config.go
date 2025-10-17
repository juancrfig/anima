package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func ConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration settings.",
		Long:  `Get or set configuration values for Anima, such as your default location.`,
	}

	cmd.AddCommand(configSetCmd())
	// You can add a 'config get' command here later
	return cmd
}

func configSetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set [key] [value]",
		Short: "Set a configuration value.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			value := args[1]

			// 1. Retrieve services from the context
			services, err := GetServices(cmd.Context())
			if err != nil {
				return err
			}

			// 2. Use the config service to set the value
			if err := services.Config.Set(key, value); err != nil {
				return fmt.Errorf("failed to set config: %w", err)
			}

			fmt.Printf("Config updated: %s = %s\n", key, value)
			return nil
		},
	}
	return cmd
}