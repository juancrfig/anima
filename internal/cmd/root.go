package cmd

import (
	"log"
	"io"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "anima",
	Short: "Anima is a personal journal. Store your thoughts and experiences safely!",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := ensureAnimaDir(); err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		InitialGreeting(nil)
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
