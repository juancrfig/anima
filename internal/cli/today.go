package cli

import (
    "fmt"
    "github.com/spf13/cobra"
)


var todayCmd = &cobra.Command {
    Use: "today",
    Short: "Open or create today's journal entry",
    Long: `
    This command opens the journal entry for the current date.
    If an entry doesn't exist, it creates a new one.`,
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("today's journal entry opened/created")
    },
}

func init() {
    rootCmd.AddCommand(todayCmd)
}
