package config

import (
	"encoding/json"
    "encoding/base64"
	"errors"
	"fmt"
	"os"
    "strconv"
	"sync"
    "time"
    "anima/internal/crypto"
)


// ErrKeyNotFound is returned when a key is not found in the configuration.
var ErrKeyNotFound = errors.New("key not found")

const (
    keyDBPath          = "db_path"
    keySessionDuration = "security.session_duration_minutes"
    keyCryptoTime      = "security.crypto.time"
    keyCryptoMemory    = "securiy.crypto.memory_kib"
    keyCryptoThreads   = "security.crypto.threads"
    keyCryptoSaltLen   = "security.crypto.salt_len"
    keyCryptoKeyLen    = "security.crypto.key_len"

    // Stores the master data key, encrypted by the user's password
    keyEncryptedMaster = "security.vault.master_key"
    // Stores the master data key, encrypted by the recovery phrase
    keyEncryptedRecovery = "security.vault.recovery_key"
)


// Config manages Anima's configuration settings.
type Config struct {
	filePath string
	data     map[string]string
	mu       sync.RWMutex
}

// New loads the configuration from a file or creates a new one if it doesn't exist.
func New(path string) (*Config, error) {
	c := &Config{
		filePath: path,
		data:     make(map[string]string),
	}

	if err := c.load(); err != nil {
		return nil, err
	}

	return c, nil
}

// load reads the configuration file from disk.
func (c *Config) load() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	file, err := os.ReadFile(c.filePath)
	if err != nil {
		// If the file doesn't exist, that's okay. We'll create it on the first Set.
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("could not read config file: %w", err)
	}

	return json.Unmarshal(file, &c.data)
}

// save writes the current configuration to disk.
func (c *Config) save() error {
	// Use MarshalIndent for a human-readable JSON file.
	data, err := json.MarshalIndent(c.data, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal config data: %w", err)
	}

	// WriteFile handles creating the file if it doesn't exist.
	return os.WriteFile(c.filePath, data, 0644)
}

// Get retrieves a value for a given key.
func (c *Config) Get(key string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, ok := c.data[key]
	if !ok {
		return "", fmt.Errorf("%w: %s", ErrKeyNotFound, key)
	}
	return value, nil
}

// Set saves a key-value pair and persists it to disk.
func (c *Config) Set(key, value string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value

	// Persist the change immediately.
	if err := c.save(); err != nil {
		return fmt.Errorf("could not save config file: %w", err)
	}
	return nil
}


// getWithDefault fetches a key, returning a default if not found.
// This helper does not need to lock, as its callers will.
func (c *Config) getWithDefault(key, defaultValue string) string {
	value, ok := c.data[key]
	if !ok {
		return defaultValue
	}
	return value
}

// parseInt parses a string to uint64.
func parseInt(s string, bitSize int) (uint64, error) {
	val, err := strconv.ParseUint(s, 10, bitSize)
	if err != nil {
		return 0, fmt.Errorf("could not parse value %q: %w", s, err)
	}
	return val, nil
}


// DBPath returns the database path.
// This is a required field and has no default.
func (c *Config) DBPath() (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	val, ok := c.data[keyDBPath]
	if !ok {
		return "", fmt.Errorf("%w: %s (this is a required field)", ErrKeyNotFound, keyDBPath)
	}
	return val, nil
}

// SessionDuration fetches and parses the session duration.
// Defaults to 0 (non-expiring) if not set.
func (c *Config) SessionDuration() (time.Duration, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Default to 0 minutes (non-expiring)
	valStr := c.getWithDefault(keySessionDuration, "0")

	minutes, err := parseInt(valStr, 64) // time.Duration is int64
	if err != nil {
		return 0, fmt.Errorf("invalid session duration: %w", err)
	}

	return time.Duration(minutes) * time.Minute, nil
}

// CryptoParams fetches and parses all crypto settings.
// It applies safe defaults for any missing values.
func (c *Config) CryptoParams() (*crypto.Params, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// 1. Get all values, applying defaults for each
	// These defaults are our OWASP recommendations
	timeStr := c.getWithDefault(keyCryptoTime, "3")
	memStr := c.getWithDefault(keyCryptoMemory, "65536")
	threadsStr := c.getWithDefault(keyCryptoThreads, "1")
	saltLenStr := c.getWithDefault(keyCryptoSaltLen, "16")
	keyLenStr := c.getWithDefault(keyCryptoKeyLen, "32")

	// 2. Parse all values
	timeVal, err := parseInt(timeStr, 32)
	if err != nil {
		return nil, err
	}

	memVal, err := parseInt(memStr, 32)
	if err != nil {
		return nil, err
	}

	threadsVal, err := parseInt(threadsStr, 8)
	if err != nil {
		return nil, err
	}

	saltLenVal, err := parseInt(saltLenStr, 8)
	if err != nil {
		return nil, err
	}

	keyLenVal, err := parseInt(keyLenStr, 8)
	if err != nil {
		return nil, err
	}

	// 3. Construct the struct
	params := &crypto.Params{
		Time:    uint32(timeVal),
		Memory:  uint32(memVal),
		Threads: uint8(threadsVal),
		SaltLen: uint8(saltLenVal),
		KeyLen:  uint8(keyLenVal),
	}

	return params, nil
}

// SetDBPath provides a typed helper for setting the DB path.
// This is safer than cfg.Set("db_path", ...).
func (c *Config) SetDBPath(path string) error {
	return c.Set(keyDBPath, path)
}


// setBytes stores raw bytes in the config as a Base64-encoded string.
func (c *Config) setBytes(key string, data []byte) error {
	encodedData := base64.StdEncoding.EncodeToString(data)
	return c.Set(key, encodedData)
}

// getBytes retrieves Base64-encoded data and decodes it back into raw bytes.
func (c *Config) getBytes(key string) ([]byte, error) {
	encodedData, err := c.Get(key)
	if err != nil {
		return nil, err // This will be ErrKeyNotFound if it doesn't exist
	}
	
	decodedData, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return nil, fmt.Errorf("could not decode config value for key %s: %w", key, err)
	}
	return decodedData, nil
}

// SetEncryptedMasterKey stores the password-encrypted master key.
func (c *Config) SetEncryptedMasterKey(keyData []byte) error {
	return c.setBytes(keyEncryptedMaster, keyData)
}

// GetEncryptedMasterKey retrieves the password-encrypted master key.
func (c *Config) GetEncryptedMasterKey() ([]byte, error) {
	return c.getBytes(keyEncryptedMaster)
}

// SetEncryptedRecoveryKey stores the recovery-phrase-encrypted master key.
func (c *Config) SetEncryptedRecoveryKey(keyData []byte) error {
	return c.setBytes(keyEncryptedRecovery, keyData)
}

// GetEncryptedRecoveryKey retrieves the recovery-phrase-encrypted master key.
func (c *Config) GetEncryptedRecoveryKey() ([]byte, error) {
	return c.getBytes(keyEncryptedRecovery)
}

// IsSetup checks if the Anima vault has been initialized.
// This is true only if *both* keys are present.
func (c *Config) IsSetup() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, masterOK := c.data[keyEncryptedMaster]
	_, recoveryOK := c.data[keyEncryptedRecovery]

	return masterOK && recoveryOK
}
