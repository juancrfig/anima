package cli

import (
	"errors"
	"fmt"
	"time"

	"anima/internal/config"
	"github.com/spf13/cobra"
)

func YesterdayCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "yesterday",
		Short: "Create or open yesterday's journal entry.",
		Long:  `Creates a new journal entry for yesterday if one doesn't exist, then opens it.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// 1. Retrieve services from the context
			services, err := GetServices(cmd.Context())
			if err != nil {
				return err
			}

			// 2. Get location from the retrieved config
			location, err := services.Config.Get("location")
			if err != nil {
				if errors.Is(err, config.ErrKeyNotFound) {
					location = ""
				} else {
					return fmt.Errorf("failed to get location from config: %w", err)
				}
			}

			// 3. Call the shared logic with yesterday's date
			yesterday := time.Now().Add(-24 * time.Hour)
			return runJournalLogic(services.Store, location, yesterday)
		},
	}
	return cmd
}