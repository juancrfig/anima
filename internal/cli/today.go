// internal/cli/today.go
package cli

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
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
			// --- 1. Setup Storage ---
			animaPath, err := GetAnimaPath()
			if err != nil {
				return err
			}
			// The database is now the single source of truth.
			dbPath := filepath.Join(animaPath, "anima.db")
			store, err := storage.New(dbPath)
			if err != nil {
				return fmt.Errorf("failed to initialize database: %w", err)
			}
			defer store.Close()

			// --- 2. Get or Create Entry ---
			today := time.Now()
			entry, err := store.GetEntryByDate(today)
			if err != nil {
				if err == sql.ErrNoRows {
					// Entry doesn't exist, so create a new one.
					// We'll pass empty content and a placeholder location for now.
					fmt.Println("Creating new journal entry for today...")
					entry, err = store.CreateEntry("", "Cucuta, Colombia")
					if err != nil {
						return fmt.Errorf("failed to create new entry: %w", err)
					}
				} else {
					// A real database error occurred.
					return fmt.Errorf("failed to retrieve entry: %w", err)
				}
			}

			// --- 3. The Temp File Editing Workflow ---
			// Create a temporary file to edit the content.
			tempFile, err := ioutil.TempFile("", "anima-*.md")
			if err != nil {
				return fmt.Errorf("could not create temporary file: %w", err)
			}
			defer os.Remove(tempFile.Name()) // IMPORTANT: Clean up the temp file.

			// Write the current entry content to the temp file.
			if _, err := tempFile.WriteString(entry.Content); err != nil {
				return fmt.Errorf("could not write to temporary file: %w", err)
			}
			tempFile.Close() // Close the file so the editor can open it.

			// --- 4. Open Editor and Wait ---
			fmt.Printf("Opening journal entry in editor...\n")
			if err := OpenFileInEditor(tempFile.Name()); err != nil {
				return fmt.Errorf("failed to open editor: %w", err)
			}

			// --- 5. Read Changes and Update Database ---
			updatedContent, err := ioutil.ReadFile(tempFile.Name())
			if err != nil {
				return fmt.Errorf("could not read updated content from temp file: %w", err)
			}

			// Only update if the content has actually changed.
			if string(updatedContent) != entry.Content {
				fmt.Println("Saving changes...")
				if err := store.UpdateEntryContent(entry.ID, string(updatedContent)); err != nil {
					return fmt.Errorf("failed to save updated entry: %w", err)
				}
			} else {
				fmt.Println("No changes detected.")
			}

			return nil
		},
	}
	return cmd
}
