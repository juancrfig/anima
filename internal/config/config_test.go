package config

import (
	"os"
	"path/filepath"
	"testing"
    "time"

    "anima/internal/crypto"
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


// TestConfig_SecurityDefaults proves our app gets safe,
// default values if no security config is set.
func TestConfig_SecurityDefaults(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "defaults_config.json")

	// Create a new, empty config file
	cfg, err := New(configPath)
	require.NoError(t, err)

	// --- 1. Test CryptoParams ---
	t.Run("it should return default crypto params when empty", func(t *testing.T) {
		// These are our OWASP defaults
		wantParams := &crypto.Params{
			Time:    3,
			Memory:  65536,
			Threads: 1,
			SaltLen: 16,
			KeyLen:  32,
		}

		// This new method must parse from the map (or use defaults)
		gotParams, err := cfg.CryptoParams()

		require.NoError(t, err)
		assert.Equal(t, wantParams, gotParams)
	})

	// --- 2. Test SessionDuration ---
	t.Run("it should return default session duration when empty", func(t *testing.T) {
		// Default to 0 (non-expiring session)
		wantDuration := 0 * time.Minute

		gotDuration, err := cfg.SessionDuration()

		require.NoError(t, err)
		assert.Equal(t, wantDuration, gotDuration)
	})

	// --- 3. Test DBPath ---
	t.Run("it should return an error for DBPath when empty", func(t *testing.T) {
		// DBPath is different. It has no safe default.
		// The app MUST fail if it's not set.
		_, err := cfg.DBPath()
		assert.Error(t, err, "DBPath() should return an error if not set")
		assert.ErrorIs(t, err, ErrKeyNotFound)
	})
}

// TestConfig_SecurityCustom proves we can load custom,
// typed security values from the K-V store.
func TestConfig_SecurityCustom(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "custom_config.json")

	cfg, err := New(configPath)
	require.NoError(t, err)

	// --- Act ---
	// Set custom security values using the existing K-V store
	require.NoError(t, cfg.Set(keyCryptoTime, "1"))
	require.NoError(t, cfg.Set(keyCryptoMemory, "1024"))
	require.NoError(t, cfg.Set(keyCryptoThreads, "2"))
	require.NoError(t, cfg.Set(keyCryptoSaltLen, "32"))
	require.NoError(t, cfg.Set(keyCryptoKeyLen, "32"))
	require.NoError(t, cfg.Set(keySessionDuration, "15"))
	require.NoError(t, cfg.Set(keyDBPath, "/tmp/test.db"))

	// --- Assert CryptoParams ---
	wantParams := &crypto.Params{
		Time:    1,
		Memory:  1024,
		Threads: 2,
		SaltLen: 32,
		KeyLen:  32,
	}
	gotParams, err := cfg.CryptoParams()
	require.NoError(t, err)
	assert.Equal(t, wantParams, gotParams)

	// --- Assert SessionDuration ---
	wantDuration := 15 * time.Minute
	gotDuration, err := cfg.SessionDuration()
	require.NoError(t, err)
	assert.Equal(t, wantDuration, gotDuration)

	// --- Assert DBPath ---
	wantDBPath := "/tmp/test.db"
	gotDBPath, err := cfg.DBPath()
	require.NoError(t, err)
	assert.Equal(t, wantDBPath, gotDBPath)
}
