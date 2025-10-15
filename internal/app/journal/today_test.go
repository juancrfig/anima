package journal

import (
    "os"
    "path/filepath"
    "testing"
    "time"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)


func TestGetOrCreateEntry(t *testing.T) {
    // Create a temporary directory for test files
    tempDir := t.TempDir()

    testDate := time.Date(2025, 10, 15, 0, 0, 0, 0, time.UTC)
    expectedFilename := "2025-10-15.md"
    expectedPath := filepath.Join(tempDir, expectedFilename)


    actualPath, err := GetOrCreateEntry(tempDir, testDate)

    require.NoError(t, err)
    assert.Equal(t, expectedPath, actualPath)

     _, err = os.Stat(expectedPath)
    assert.NoError(t, err, "The journal entry file should exist")
}
