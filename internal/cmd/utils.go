package cmd

import (
	"os"
	"errors"
	"os/exec"
	"path/filepath"
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

func openEntry(cmd *cobra.Command, date string) error {
	textEditor := os.Getenv("EDITOR")
	if textEditor == "" {
		return errors.New("No text editor detected") 
	}

	ctx := cmd.Context()
	entriesPath := ctx.Value(entriesPathKey)

	if s, ok := entriesPath.(string); ok {

		entryPath := filepath.Join(s, date + ".md")

		editorCmd := exec.Command(textEditor, entryPath)
		editorCmd.Stdin = os.Stdin
		editorCmd.Stdout = os.Stdout
		editorCmd.Stderr = os.Stderr

		err := editorCmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

func initialGreeting(w io.Writer) {
	if w == nil {
		w = os.Stdout
	}

	fmt.Fprintf(w, "Hello! I am Anima.")
}
