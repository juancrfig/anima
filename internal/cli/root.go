package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "anima",
	Short: "A personal journal command-line tool",
	Run: func(cmd *cobra.Command, args []string) {
        cmd.Println("Hello, Journal!")
	},
}

func Execute() error {
	return rootCmd.Execute()
}
