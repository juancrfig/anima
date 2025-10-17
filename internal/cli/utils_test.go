package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAnimaPath(t *testing.T) {
	// ARRANGE
	// We get the real home directory to construct our expected path.
	// This makes the test independent of any specific user's machine.
	homeDir, err := os.UserHomeDir()
	require.NoError(t, err, "Test setup failed: could not get user home directory")
	expectedPath := filepath.Join(homeDir, ".anima")

	// ACT
	actualPath, err := GetAnimaPath()

	// ASSERT
	assert.NoError(t, err)
	assert.Equal(t, expectedPath, actualPath, "GetAnimaPath() should return the correct default directory")
}

func TestCreateFileIfNotExists(t *testing.T) {
	// t.TempDir() creates a temporary directory for this specific test
	// and automatically cleans it up when the test finishes.
	// This is the standard, safe way to test filesystem operations.
	tempDir := t.TempDir()

	t.Run("it should create a file and its parent directories if they do not exist", func(t *testing.T) {
		// ARRANGE
		// Define a path inside the temp directory that includes a non-existent subdirectory.
		testPath := filepath.Join(tempDir, "new_subdir", "test_entry.md")

		// ACT
		err := CreateFileIfNotExists(testPath)

		// ASSERT
		require.NoError(t, err, "CreateFileIfNotExists should not return an error on first creation")
		// os.Stat returns info about the file. If it returns an error, the file doesn't exist.
		_, err = os.Stat(testPath)
		assert.NoError(t, err, "File should exist at the specified path after creation")
	})

	t.Run("it should do nothing and not return an error if the file already exists", func(t *testing.T) {
		// ARRANGE
		// First, create the file to establish the "already exists" state.
		testPath := filepath.Join(tempDir, "existing_entry.md")
		err := CreateFileIfNotExists(testPath)
		require.NoError(t, err, "Test setup failed: could not create initial file")

		// ACT
		// Call the function a second time on the same path.
		err = CreateFileIfNotExists(testPath)

		// ASSERT
		assert.NoError(t, err, "CreateFileIfNotExists should not return an error if the file already exists")
	})
}