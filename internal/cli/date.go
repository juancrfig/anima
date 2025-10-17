// internal/cli/date.go
package cli

import (
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"anima/internal/config"
	"anima/internal/storage"

	"github.com/spf13/cobra"
)

func DateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "date [YYYY-MM-DD]",
		Short: "Create or open a journal entry for a specific date.",
		Long:  `Creates or opens a journal entry for the date specified in YYYY-MM-DD format.`,
		// We enforce that exactly one argument must be provided.
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// 1. Parse the date argument
			dateStr := args[0]
			parsedDate, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				return fmt.Errorf("invalid date format: %q. Please use YYYY-MM-DD", dateStr)
			}

			// 2. Setup (This is identical to our other commands)
			animaPath, err := GetAnimaPath()
			if err != nil {
				return err
			}
			dbPath := filepath.Join(animaPath, "anima.db")
			configPath := filepath.Join(animaPath, "config.json")
			cfg, err := config.New(configPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			location, err := cfg.Get("location")
			if err != nil {
				if errors.Is(err, config.ErrKeyNotFound) {
					location = ""
				} else {
					return fmt.Errorf("failed to get location from config: %w", err)
				}
			}
			store, err := storage.New(dbPath)
			if err != nil {
				return fmt.Errorf("failed to initialize database: %w", err)
			}
			defer store.Close()

			// 3. Call the shared logic with the parsed date
			return runJournalLogic(store, location, parsedDate)
		},
	}
	return cmd
}