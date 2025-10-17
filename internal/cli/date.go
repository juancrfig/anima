package cli

import (
	"errors"
	"fmt"
	"time"

	"anima/internal/config"
	"github.com/spf13/cobra"
)

func DateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "date [YYYY-MM-DD]",
		Short: "Create or open a journal entry for a specific date.",
		Long:  `Creates or opens a journal entry for the date specified in YYYY-MM-DD format.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// 1. Parse the date argument
			dateStr := args[0]
			parsedDate, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				return fmt.Errorf("invalid date format: %q. Please use YYYY-MM-DD", dateStr)
			}

			// 2. Retrieve services from the context
			services, err := GetServices(cmd.Context())
			if err != nil {
				return err
			}

			// 3. Get location from the retrieved config
			location, err := services.Config.Get("location")
			if err != nil {
				if errors.Is(err, config.ErrKeyNotFound) {
					location = ""
				} else {
					return fmt.Errorf("failed to get location from config: %w", err)
				}
			}

			// 4. Call the shared logic with the parsed date
			return runJournalLogic(services.Store, location, parsedDate)
		},
	}
	return cmd
}