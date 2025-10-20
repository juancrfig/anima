package cli

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

    "anima/internal/auth"
	"anima/internal/storage"
)

// runJournalLogic handles the core workflow for opening, editing, and saving
// a journal entry for any given date.
func runJournalLogic(store *storage.Storage, location string, journalDate time.Time) error {
	isNewEntry := false

	// 1. Attempt to get the entry for the specified date
	entry, err := store.GetEntryByDate(journalDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			isNewEntry = true
			// Create a temporary, in-memory placeholder. DO NOT save to DB yet.
			entry = &storage.Entry{Content: ""}
		} else if errors.Is(err, auth.ErrSessionExpired) {
            return errors.New("Session expired. Please run 'anima login' first")
        } else {
			// For any other error, we should stop.
			return fmt.Errorf("failed to retrieve entry: %w", err)
		}
	}

	// 2. Open content in the editor
	tempFile, err := ioutil.TempFile("", "anima-*.md")
	if err != nil {
		return fmt.Errorf("could not create temporary file: %w", err)
	}

	if _, err := tempFile.WriteString(entry.Content); err != nil {
        os.Remove(tempFile.Name())
		return fmt.Errorf("could not write to temporary file: %w", err)
	}
	tempFile.Close() // Close file before handing off to editor

	if err := OpenFileInEditor(tempFile.Name()); err != nil {
        os.Remove(tempFile.Name())
		return fmt.Errorf("failed to open editor: %w", err)
	}

	// 3. After editing, read content and decide whether to save
	updatedContent, err := ioutil.ReadFile(tempFile.Name())
	if err != nil {
        os.Remove(tempFile.Name())
		return fmt.Errorf("could not read updated content from temp file: %w", err)
	}

	updatedContentStr := strings.TrimSpace(string(updatedContent))
	originalContentStr := strings.TrimSpace(entry.Content)

    // Decide whether to save
	if isNewEntry {
		if updatedContentStr == "" {
			fmt.Println("No content added. Journal entry not created.")
            os.Remove(tempFile.Name())
			return nil
		}
		// It's a new entry and it has content, so we CREATE it.
		_, err := store.CreateEntry(updatedContentStr, location, journalDate)
		if err != nil {
            if errors.Is(err, auth.ErrSessionExpired) {
                return fmt.Errorf("Session expired. Your changes were NOT saved, but are safe in: %s\nPlease run 'anima login' and manually copy the content", tempFile.Name())
            }
            os.Remove(tempFile.Name())
			return fmt.Errorf("failed to create new entry: %w", err)
		}

		fmt.Println("Journal entry saved.")
	} else {
		if updatedContentStr == originalContentStr {
			fmt.Println("No changes detected.")
            os.Remove(tempFile.Name())
			return nil
		}
		// It's an existing entry and the content has changed, so we UPDATE it.
		err := store.UpdateEntryContent(entry.ID, updatedContentStr)
		if err != nil {
			return fmt.Errorf("failed to update entry: %w", err)
		}
		fmt.Println("Journal entry updated.")
	}
    os.Remove(tempFile.Name())
	return nil
}
