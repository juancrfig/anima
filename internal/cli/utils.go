package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// GetAnimaPath returns the absolute path to the ~/.anima directory.
func GetAnimaPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine user home directory: %w", err)
	}
	return filepath.Join(homeDir, ".anima"), nil
}

// CreateFileIfNotExists creates a file and its parent directories if it doesn't exist.
func CreateFileIfNotExists(path string) error {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("could not create directory %s: %w", dir, err)
		}
	}

	// 'OpenFile' with 'O_CREATE|O_EXCL' is an atomic "create if not exists" operation.
	// This prevents race conditions. We close it immediately as we only need to ensure it exists.
	file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		if os.IsExist(err) {
			return nil // File already exists, which is not an error for this function.
		}
		return fmt.Errorf("could not create file %s: %w", path, err)
	}
	return file.Close()
}

// OpenFileInEditor opens the specified file path in the user's default editor.
func OpenFileInEditor(filePath string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		// Provide sensible defaults for different OSes.
		switch runtime.GOOS {
		case "windows":
			editor = "notepad"
		case "darwin":
			editor = "open" // 'open' is more versatile on macOS than a specific editor.
		default:
			editor = "vim" // A safe bet on most Linux/Unix systems.
		}
	}

	cmd := exec.Command(editor, filePath)

	// We must connect the command's stdio to the terminal's stdio
	// so the user can interact with the editor (e.g., vim, nano).
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}