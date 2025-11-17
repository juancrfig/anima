package cmd

import (
	"time"
	"errors"

	"github.com/juancrfig/anima/internal/journal"

	"github.com/spf13/cobra"
)

type ctxKey string
const entriesPathKey ctxKey = "entriesPath"
const lastOpenedEntry ctxKey = "lastOpenedEntry"

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

		if args[0] == "today" {
			todayDate := time.Now().Format("2006-01-02")
			err := journal.OpenEntry(todayDate)
			if err != nil {
				return err
			}
		}

		if args[0] == "yesterday" {
			yestDate := time.Now().AddDate(0,0,-1).Format("2006-01-02")
			err := journal.OpenEntry(yestDate)
			if err != nil {
				return err
			}
		}

		if _, err := time.Parse("2006-01-02", args[0]); err != nil {
			return errors.New("Date must be like YYYY-MM-DD and be valid")
		}
		err := journal.OpenEntry(args[0])
		if err != nil {
			return err
		}
		return nil
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
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
