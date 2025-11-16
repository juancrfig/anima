package cmd

import (
	"time"
	"errors"

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
		if _, err := time.Parse("2006-01-02", args[0]); err != nil {
			return errors.New("Date must be like YYYY-MM-DD and be valid")
		}
		journal.OpenEntry(args[0])
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		return
	}
}

func init() {
	rootCmd.Flags().StringP("config", "c", "", "config file path")
}
