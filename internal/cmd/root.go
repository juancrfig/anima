package cmd

import (
	"time"
	"errors"
	"context"
	"os"
	"log"

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
		if _, err := time.Parse("2006-01-02", args[0]); err != nil {
			return errors.New("Date must be like YYYY-MM-DD and be valid")
		}
		lastEntryPath, err := journal.OpenEntry(args[0])
		if err != nil {
			return err
		}

		ctx := context.WithValue(cmd.Context(), lastOpenedEntry, lastEntryPath)
		cmd.SetContext(ctx)

		return nil
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {

		v := cmd.Context().Value(lastOpenedEntry)
		path, ok := v.(string) 
		if !ok {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		hasFrontmatter, err := journal.DetectFrontmatter(f)
		if err != nil {
			return err
		}
		log.Printf("hasFrontmatter: %v", hasFrontmatter)

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
