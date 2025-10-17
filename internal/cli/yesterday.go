// internal/cli/yesterday.go
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

func YesterdayCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "yesterday",
		Short: "Create or open yesterday's journal entry.",
		Long:  `Creates a new journal entry for yesterday if one doesn't exist, then opens it.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// --- Setup (identical to today.go) ---
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

			// --- Call Shared Logic ---
			yesterday := time.Now().Add(-24 * time.Hour)
			return runJournalLogic(store, location, yesterday)
		},
	}
	return cmd
}
