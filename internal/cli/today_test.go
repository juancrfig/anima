package cli

import (
	"context"
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"

	"anima/internal/auth"
	"anima/internal/config"
	"anima/internal/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func mockOpenFileInEditor(filePath, content string) error {
	return ioutil.WriteFile(filePath, []byte(content), 0644)
}


// It creates a full, authenticated set of services.
func setupTestServices(t *testing.T) (*Services, *storage.Storage) {
	t.Helper()

	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_anima.db")
	configPath := filepath.Join(tempDir, "test_config.json")

	// 1. Create mock config
	cfg, err := config.New(configPath)
	require.NoError(t, err)

	// 2. Create and pre-authenticate the auth manager
	masterKey := []byte("test-master-key-123")
	authMgr := auth.New()
	authMgr.SetPassword(masterKey, 0) // 0 = never expires

	// 3. Create KeyManager (needed by storage.New via services)
	// We need to set crypto params in config for this to work
	require.NoError(t, cfg.Set("security.crypto.time", "1"))
	require.NoError(t, cfg.Set("securiy.crypto.memory_kib", "1024"))
	cryptoParams, err := cfg.CryptoParams()
	require.NoError(t, err)
	keyMgr := auth.NewKeyManager(cryptoParams)

	// 4. Create the storage service
	store, err := storage.New(dbPath, cfg, authMgr)
	require.NoError(t, err, "Failed to initialize storage")
	// Override params for fast tests
	store.SetCryptoParamsForTesting(cryptoParams)

	// 5. Create the services struct
	services := &Services{
		Store:      store,
		Config:     cfg,
		Auth:       authMgr,
		KeyManager: keyMgr,
	}
	return services, store
}


func TestTodayCmd_Workflow(t *testing.T) {
	services, store := setupTestServices(t)
	testTime := time.Now().UTC().Truncate(24 * time.Hour)
	ctx := context.WithValue(context.Background(), servicesKey, services)
	ctx = context.WithValue(ctx, dbStoreKey, store)

	// We create a *mock* TodayCmd
	cmd := TodayCmd()
	cmd.SetContext(ctx)

	// We also need to override the real editor
	// with our mock one for the *real* journal logic.
	// This is tricky. For now, let's just test the logic directly.
	// We will refactor 'runTodayLogic' to 'runJournalLogic'
	// and test that, since that's what we *really* care about.

	t.Run("first run creates a new entry", func(t *testing.T) {
		// Arrange
		// We replace the 'OpenFileInEditor' with our mock
		// This is a global override for testing, which is simple.
		OpenFileInEditor = func(filePath string) error {
			return mockOpenFileInEditor(filePath, "First thoughts for the day.")
		}

		// Act
		err := runJournalLogic(store, "Test Location", testTime)
		require.NoError(t, err)

		// Assert
		entry, err := store.GetEntryByDate(testTime)
		require.NoError(t, err)
		assert.Equal(t, "First thoughts for the day.", entry.Content)
	})

	t.Run("second run updates the existing entry", func(t *testing.T) {
		// Arrange
		OpenFileInEditor = func(filePath string) error {
			return mockOpenFileInEditor(filePath, "First thoughts for the day. And some more ideas.")
		}

		// Act
		err := runJournalLogic(store, "Test Location", testTime)
		require.NoError(t, err)

		// Assert
		updatedEntry, err := store.GetEntryByDate(testTime)
		require.NoError(t, err)
		assert.Equal(t, "First thoughts for the day. And some more ideas.", updatedEntry.Content)

		var count int
		err = store.DB().QueryRow("SELECT COUNT(id) FROM entries").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 1, count, "Should only be one entry")
	})
}
