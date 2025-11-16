package cmd

import (
	"log"

	"github.com/juancrfig/anima/internal/journal"

	"github.com/spf13/cobra"
)

type ctxKey string
const entriesPathKey ctxKey = "entriesPath"

var rootCmd = &cobra.Command{
	Use: "anima [date]",
	Short: "Anima is a personal journal. Store your thoughts and experiences safely!",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return ensureAnimaDir(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			initialGreeting(nil)
			return nil
		}
		date := args[0]
		journal.OpenEntry(date)
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Panic(err)
	}
}

func init() {
	rootCmd.Flags().StringP("config", "c", "", "config file path")
}
