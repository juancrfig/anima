package cli

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)


func TodayCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "today",
		Short: "Create or open today's journal entry.",
		Long: `Creates a new journal entry for the current date if one doesn't exist,
then opens it in your default text editor.`,
		// Use RunE to handle errors returned by our functions.
		RunE: func(cmd *cobra.Command, args []string) error {
			animaPath, err := GetAnimaPath()
			if err != nil {
				return err
			}

			entriesDir := filepath.Join(animaPath, "entries")
			todayFilename := time.Now().Format("2006-01-02") + ".md"
			entryPath := filepath.Join(entriesDir, todayFilename)

			
			if err := CreateFileIfNotExists(entryPath); err != nil {
				return err
			}

			fmt.Printf("Opening journal entry: %s\n", entryPath)
			if err := OpenFileInEditor(entryPath); err != nil {
				// Provide a helpful fallback message if the editor fails.
				return fmt.Errorf("failed to open editor: %w. Your entry is at: %s", err, entryPath)
			}

			return nil
		},
	}
	return cmd
}
