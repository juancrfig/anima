package cli

import (
    "os"
    "os/exec"
    "fmt"
)


func openInEditor(filePath string) error {
    editor := os.Getenv("EDITOR")
    if editor == "" {
        editor = "vim"
    }

    executable, err := exec.LookPath(editor)
    if err != nil {
        return fmt.Errorf("Could not find editor '%s' in PATH", editor)
    }

    cmd := exec.Command(executable, filePath)

    // Connect the command's stdin, stdout, and stderr to the user's terminal
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    return cmd.Run()
}
