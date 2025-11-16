package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "anima",
	Short: "Anima is a personal journal. Store your thoughts and experiences safely!",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Root command has been executed")
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
