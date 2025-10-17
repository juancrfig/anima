package cli

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"anima/internal/config"
	"anima/internal/storage"

	"github.com/spf13/cobra"
)


func TodayCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "today",
		Short: "Create or open today's journal entry.",
		Long: `Creates a new journal entry for the current date if one doesn't exist,
then opens it in your default text editor for editing.`,
		RunE: func(cmd *cobra.Command, args []string) error {
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
					fmt.Println("Location not set")
				} else {
					return fmt.Errorf("failed to get location from config: %w", err)
				}
			}

			store, err := storage.New(dbPath)
			if err != nil {
				return fmt.Errorf("failed to initialize database: %w", err)
			}
			defer store.Close()

			today := time.Now()
			isNewEntry := false

			entry, err := store.GetEntryByDate(today)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					isNewEntry = true
					entry = &storage.Entry{Content: ""}
				} else {
					return fmt.Errorf("failed to retrieve entry: %w", err)
				}
			}

			// Open content in editor
			tempFile, err := ioutil.TempFile("", "anima-*.md")
			if err != nil {
				return fmt.Errorf("could not create temporary file: %w", err)
			}
			defer os.Remove(tempFile.Name())

			if _, err := tempFile.WriteString(entry.Content); err != nil {
				return fmt.Errorf("could not write to temporary file: %w", err)
			}
			tempFile.Close()

			if err := OpenFileInEditor(tempFile.Name()); err != nil {
				return fmt.Errorf("failed to open editor: %w", err)
			}

			// Read and trim updated content once
			updatedContent, err := ioutil.ReadFile(tempFile.Name())
			if err != nil {
				return fmt.Errorf("could not read updated content: %w", err)
			}
			updatedStr := strings.TrimSpace(string(updatedContent))
			originalStr := strings.TrimSpace(entry.Content)

			switch {
			case isNewEntry && updatedStr == "":
				fmt.Println("No content added. Journal entry not created.")
			case isNewEntry:
				if _, err := store.CreateEntry(updatedStr, location, today); err != nil {
					return fmt.Errorf("failed to create new entry: %w", err)
				}
				fmt.Println("Journal entry saved")
			case !isNewEntry && updatedStr == originalStr:
				fmt.Println("No changes detected")
			case !isNewEntry:
				if err := store.UpdateEntryContent(entry.ID, updatedStr); err != nil {
					return fmt.Errorf("failed to update entry: %w", err)
				}
				fmt.Println("Journal entry updated")
			}

			return nil
		},
	}

	return cmd
}

