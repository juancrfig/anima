package cli

import (
	"path/filepath"
	"testing"
	"time"
    "io/ioutil"
    "os"
    "database/sql"

	"anima/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


// This is a mock implementation of the editor-opening function.
// In our test, we don't want to actually open Vim or Notepad.
// We just want to simulate the user editing a file.
func mockOpenFileInEditor(filePath, content string) error {
	// Simulate the user writing new content to the file.
	return ioutil.WriteFile(filePath, []byte(content), 0644)
}

func TestTodayCmd_Workflow(t *testing.T) {
	// --- Setup ---
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_anima.db")
	store, err := storage.New(dbPath)
	require.NoError(t, err)
	defer store.Close()

    testTime := time.Now().UTC().Truncate(24 * time.Hour)

	// --- Act & Assert ---
	t.Run("first run creates a new entry", func(t *testing.T) {
		// --- Arrange ---
		// We pass our storage and mock editor function to a helper
		err := runTodayLogic(store, "First thoughts for the day.", testTime)
		require.NoError(t, err)

		// --- Assert ---
		// Check the database directly
		entry, err := store.GetEntryByDate(testTime)
		require.NoError(t, err)
		assert.Equal(t, "First thoughts for the day.", entry.Content)
	})

	t.Run("second run updates the existing entry", func(t *testing.T) {
		// --- Arrange ---
		// On the second run, we simulate adding more text.
		err := runTodayLogic(store, "First thoughts for the day. And some more ideas.", testTime)
		require.NoError(t, err)

		// --- Assert ---
		// Verify the content was updated
		updatedEntry, err := store.GetEntryByDate(testTime)
		require.NoError(t, err)
		assert.Equal(t, "First thoughts for the day. And some more ideas.", updatedEntry.Content)

		// The most important assertion: verify no new rows were created.
		var count int
		err = store.DB().QueryRow("SELECT COUNT(id) FROM entries").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 1, count, "Should only be one entry in the database")
	})
}

// runTodayLogic is a helper function that isolates the core logic of the TodayCmd's RunE.
// This makes the logic testable without involving Cobra directly.
func runTodayLogic(store *storage.Storage, newContent string, today time.Time) error {
	entry, err := store.GetEntryByDate(today)
	if err != nil {
		if err == sql.ErrNoRows {
			// Location is blank for the test
			entry, err = store.CreateEntry("", "", today)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	tempFile, err := ioutil.TempFile("", "anima-*.md")
	if err != nil {
		return err
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.WriteString(entry.Content); err != nil {
		return err
	}
	tempFile.Close()

	// Use our mock editor function here
	if err := mockOpenFileInEditor(tempFile.Name(), newContent); err != nil {
		return err
	}

	updatedContent, err := ioutil.ReadFile(tempFile.Name())
	if err != nil {
		return err
	}

	if string(updatedContent) != entry.Content {
		return store.UpdateEntryContent(entry.ID, string(updatedContent))
	}

	return nil
}
