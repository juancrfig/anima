package journal

import (
    "fmt"
    "os"
    "path/filepath"
    "time"
)


func GetOrCreateEntry(journalDir string, date time.Time) (string, error) {
    filename := date.Format("2006-01-02") + ".md"
    entryPath := filepath.Join(journalDir, filename)

    _, err := os.Stat(entryPath)
    if err == nil {
        return entryPath, nil
    }

    if !os.IsNotExist(err) {
        return "", fmt.Errorf("Failed to check for entry file: %w", err)
    }

    file, err := os.Create(entryPath)
    if err != nil {
        return "", fmt.Errorf("Failed to create entry file: %w", err)
    }
    file.Close()

    return entryPath, nil
}
