package cmd

import (
	"path/filepath"
	"os"
	"errors"
)

func ensureAnimaDir() error {
 	homePath, err := os.UserHomeDir()
	if err != nil {
		return errors.New("Error getting home path")
	}

	animaDirPath := filepath.Join(homePath, ".anima", "entries")
	
	err = os.MkdirAll(animaDirPath, 0700)
	if err != nil {
		return errors.New("Error creating the anima directory")
	}

	return nil
}
