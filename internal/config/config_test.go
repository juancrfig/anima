package config

import (
	"os"
	"path/filepath"
	"testing"
    "time"
    "encoding/base64"

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


// TestConfig_SecurityKeyStorage proves we can set and get the
// specific, high-security keys required for the "Recovery Key" architecture.
func TestConfig_SecurityKeyStorage(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "key_storage_config.json")

	cfg, err := New(configPath)
	require.NoError(t, err)

	// 1. On a new config, it should not be set up
	assert.False(t, cfg.IsSetup(), "IsSetup() should be false for a new config")

	// 2. Define our mock encrypted keys
	// (In a real setup, these are long, encrypted byte slices)
	mockEncryptedMasterKey := []byte("encrypted-master-key-data")
	mockEncryptedRecoveryKey := []byte("encrypted-recovery-key-data")

	// 3. Set the keys
	err = cfg.SetEncryptedMasterKey(mockEncryptedMasterKey)
	require.NoError(t, err)
	err = cfg.SetEncryptedRecoveryKey(mockEncryptedRecoveryKey)
	require.NoError(t, err)

	// 4. Now, the config should report as "set up"
	assert.True(t, cfg.IsSetup(), "IsSetup() should be true after setting keys")

	// 5. Get the keys back
	retrievedMaster, err := cfg.GetEncryptedMasterKey()
	require.NoError(t, err)
	assert.Equal(t, mockEncryptedMasterKey, retrievedMaster)

	retrievedRecovery, err := cfg.GetEncryptedRecoveryKey()
	require.NoError(t, err)
	assert.Equal(t, mockEncryptedRecoveryKey, retrievedRecovery)

	// 6. Verify they are stored as base64 strings in the raw map
	// This proves we are storing raw bytes safely in a JSON-friendly format.
	rawMaster, err := cfg.Get(keyEncryptedMaster)
	require.NoError(t, err)
	assert.Equal(t, base64.StdEncoding.EncodeToString(mockEncryptedMasterKey), rawMaster)
}

// TestConfig_IsSetup_False proves IsSetup is false if only one key is present.
func TestConfig_IsSetup_False(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "partial_config.json")

	cfg, err := New(configPath)
	require.NoError(t, err)

	// Set only the master key
	err = cfg.SetEncryptedMasterKey([]byte("only-master-key"))
	require.NoError(t, err)

	// It should still report as not set up
	assert.False(t, cfg.IsSetup(), "IsSetup() should be false if recovery key is missing")
}

