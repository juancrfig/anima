package cmd

import (
	"os"
	"errors"
	"os/exec"
)

func openEntry(entry string) error {
	textEditor := os.Getenv("EDITOR")
	if textEditor == "" {
		return errors.New("No text editor detected") 
	}

	cmd := exec.Command(textEditor, entry)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
