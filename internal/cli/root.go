package cli

import (
    "fmt"
    "os"
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command {
    Use: "anima [command] [flags]",
    Short: "Anima is a personal AI-powered journal, designed to  get to know you",
    Long: `
    This is a command-line tool that serves you as a simple personal journal. 
    You can write your diary entries, and they will be saved securely in a JSON file on your local device.
    The more you write, the better you and Anima will get to know yourself.`,
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
