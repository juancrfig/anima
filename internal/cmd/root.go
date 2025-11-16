package cmd

import (
	"log"
	"io"
	"fmt"
	"os"
	"context"
	"path/filepath"

	"github.com/spf13/cobra"
)

type ctxKey string
const pathKey ctxKey = "animaPath"

var rootCmd = &cobra.Command{
	Use: "anima [date]",
	Short: "Anima is a personal journal. Store your thoughts and experiences safely!",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := ensureAnimaDir(cmd); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			InitialGreeting(nil)
			return nil
		}

		dateArg := args[0]
		return openEntry(dateArg)
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


func InitialGreeting(w io.Writer) {
	if w == nil {
		w = os.Stdout
	}
	fmt.Fprintf(w, "Hello! I'm Anima.")
}
