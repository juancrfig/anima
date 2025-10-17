package cli

import (
	"errors"
	"fmt"
	"time"

	"anima/internal/config"
	"github.com/spf13/cobra"
)

func TodayCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "today",
		Short: "Create or open today's journal entry.",
		Long: `Creates a new journal entry for the current date if one doesn't exist,
then opens it in your default text editor for editing.`,
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
					fmt.Println("Location not set. Use 'anima config set location \"city, country\"' to set it.")
				} else {
					return fmt.Errorf("failed to get location from config: %w", err)
				}
			}

			// 3. Call the shared logic with today's date
			return runJournalLogic(services.Store, location, time.Now())
		},
	}
	return cmd
}