package journal

import (
	"path/filepath"
	"os"
	"os/exec"
	"errors"
)

func OpenEntry(date string) (string, error) {
	textEditor := os.Getenv("EDITOR")
	if textEditor == "" {
		return "", errors.New("No text editor detected") 
	}

	home, _ := os.UserHomeDir()
	s := filepath.Join(home,".anima", "entries")
	entryPath := filepath.Join(s, date + ".md")

	editorCmd := exec.Command(textEditor, entryPath)
	editorCmd.Stdin = os.Stdin
	editorCmd.Stdout = os.Stdout
	editorCmd.Stderr = os.Stderr

	err := editorCmd.Run()
	if err != nil {
		return "", err
	}
	return entryPath, nil
}


