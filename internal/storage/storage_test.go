// internal/storage/storage_test.go
package storage

import (
    "database/sql"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)



func TestGetOrCreateEntryByDate(t *testing.T) {
	dbPath := t.TempDir() + "/test_anima.db"
	storage, err := New(dbPath)
	require.NoError(t, err, "Failed to initialize storage")
	defer storage.Close()

	today := time.Now().UTC().Truncate(24 * time.Hour)
	location := "Test City"

	// --- 1. First call: Entry should be created ---
	entry1, created, err := storage.GetOrCreateEntryByDate(today, location)
	require.NoError(t, err)
	assert.True(t, created, "The entry should have been created on the first call")
	assert.NotZero(t, entry1.ID)
	assert.Equal(t, "", entry1.Content, "A new entry should have empty content")
	assert.Equal(t, location, entry1.Location)
	assert.Equal(t, today, entry1.Date.Truncate(24*time.Hour))

	// --- 2. Second call: Existing entry should be retrieved ---
	entry2, created, err := storage.GetOrCreateEntryByDate(today, location)
	require.NoError(t, err)
	assert.False(t, created, "The entry should not be created on the second call")
	assert.Equal(t, entry1.ID, entry2.ID, "The same entry ID should be returned")

	// --- 3. Verification: Ensure only one row exists in the DB ---
	var count int
	err = storage.db.QueryRow("SELECT COUNT(id) FROM entries WHERE date = ?", today).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count, "There should be exactly one entry for the given date")
}


func TestCreateAndGetEntry(t *testing.T) {
	
    dbPath := "test_anima.db"
	defer os.Remove(dbPath)

	// Initialize the storage
	storage, err := New(dbPath)
	require.NoError(t, err, "Failed to initialize storage")
	defer storage.Close()

	// --- Test Case ---
	t.Run("creates a new entry and retrieves it successfully", func(t *testing.T) {
		now := time.Now().UTC().Truncate(time.Second) // Truncate for consistent comparison
		entryToCreate := &Entry{
			Content:   "This is a test journal entry.",
			CreatedAt: now,
			Location:  "Cucuta, Colombia",
		}

		createdEntry, err := storage.CreateEntry(entryToCreate.Content, entryToCreate.Location, now)
		require.NoError(t, err, "CreateEntry should not return an error")

		retrievedEntry, err := storage.GetEntry(createdEntry.ID)
		require.NoError(t, err, "GetEntry should not return an error")

		assert.NotZero(t, createdEntry.ID, "Created entry ID should not be zero")
		assert.Equal(t, createdEntry.ID, retrievedEntry.ID, "Retrieved ID should match created ID")
		assert.Equal(t, "This is a test journal entry.", retrievedEntry.Content, "Content should match")
		assert.Equal(t, "Cucuta, Colombia", retrievedEntry.Location, "Location should match")
		assert.Equal(t, now.UTC().Truncate(time.Minute), retrievedEntry.CreatedAt, "Timestamp should match")
	})
}


func TestGetAndUpdateEntryByDate(t *testing.T) {
	dbPath := t.TempDir() + "/test_anima.db"
	storage, err := New(dbPath)
	require.NoError(t, err, "Failed to initialize storage")
	defer storage.Close()

	// --- Test Cases ---
	t.Run("returns an error if no entry is found for the date", func(t *testing.T) {
		today := time.Now().UTC().Truncate(24 * time.Hour)
		_, err := storage.GetEntryByDate(today)

		assert.Error(t, err, "Expected an error when no entry is found")
		// We expect the specific error sql.ErrNoRows from the database driver
		assert.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("gets and updates an existing entry", func(t *testing.T) {
	
		originalContent := "This is the original content."
		today := time.Now().UTC().Truncate(24 * time.Hour)

		_, err := storage.db.Exec(
			"INSERT INTO entries(content, location, created_at, date) VALUES(?, ?, ?, ?)",
			originalContent, "Test Location", time.Now().UTC(), today,
		)
		require.NoError(t, err, "Test setup failed: could not insert test entry")

		entry, err := storage.GetEntryByDate(today)
		require.NoError(t, err, "GetEntryByDate should not return an error for an existing entry")
		assert.Equal(t, originalContent, entry.Content)

		updatedContent := "This is the updated content."
		err = storage.UpdateEntryContent(entry.ID, updatedContent)
		require.NoError(t, err, "UpdateEntryContent should not return an error")

		finalEntry, err := storage.GetEntry(entry.ID)
		require.NoError(t, err)
		assert.Equal(t, updatedContent, finalEntry.Content, "Content should have been updated")
	})
}
