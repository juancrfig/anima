package cli

import (
    "anima/internal/app/journal"
    "fmt"
    "os"
    "path/filepath"
    "time"
    "github.com/spf13/cobra"
)


var todayCmd = &cobra.Command {
    Use: "today",
    Short: "Open or create today's journal entry",
    Long: `
    This command opens the journal entry for the current date.
    If an entry doesn't exist, it creates a new one.`,
    Run: func(cmd *cobra.Command, args []string) {
        homeDir, err := os.UserHomeDir()
        if err != nil {
            fmt.Printf("Error: could not find home directory: %v\n", err)
            os.Exit(1)
        }
        journalDir := filepath.Join(homeDir, ".anima", "entries")

        // Ensure the directory exists
        if err := os.MkdirAll(journalDir, 0755); err != nil {
            fmt.Printf("Error: could not create journal directory: %v\n", err)
            os.Exit(1)
        }

        entryPath, err := journal.GetOrCreateEntry(journalDir, time.Now())
        if err != nil {
            fmt.Printf("Error: could not get or create journal entry: %v\n", err)
            os.Exit(1)
        }

        if err := openInEditor(entryPath); err != nil {
            fmt.Printf("Error: could not open file in editor: %v\n", err)
            fmt.Printf("Your entry is located at: %s\n", entryPath)
            os.Exit(1)
        }
    },
}

func init() {
    rootCmd.AddCommand(todayCmd)
}
