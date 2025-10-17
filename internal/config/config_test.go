package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	// --- Setup ---
	// Use a temporary directory to avoid interfering with any real config file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	t.Run("it should set and get a value", func(t *testing.T) {
		cfg, err := New(configPath)
		require.NoError(t, err)

		err = cfg.Set("location", "Cucuta, Colombia")
		require.NoError(t, err)
		
		value, err := cfg.Get("location")

		assert.NoError(t, err)
		assert.Equal(t, "Cucuta, Colombia", value)
	})

	t.Run("it should persist a value to the file", func(t *testing.T) {
		// Use a separate config path for this test to ensure isolation
		persistConfigPath := filepath.Join(tempDir, "persist_config.json")
		
		// First, create a config, set a value, which implicitly saves it.
		cfg1, err := New(persistConfigPath)
		require.NoError(t, err)
		err = cfg1.Set("user.name", "Juanes")
		require.NoError(t, err)

		// --- Act ---
		// Now, create a *new* config instance from the same file path.
		// This simulates the app restarting and loading the config from disk.
		cfg2, err := New(persistConfigPath)
		require.NoError(t, err)
		value, err := cfg2.Get("user.name")

		// --- Assert ---
		assert.NoError(t, err)
		assert.Equal(t, "Juanes", value)

		// Also verify the file was actually created
		_, err = os.Stat(persistConfigPath)
		assert.NoError(t, err, "Config file should exist on disk")
	})

	t.Run("it should return an error when getting a non-existent key", func(t *testing.T) {
		cfg, err := New(configPath)
		require.NoError(t, err)

		_, err = cfg.Get("non.existent.key")

		assert.Error(t, err)
	})
}
