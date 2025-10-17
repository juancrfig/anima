// internal/config/config.go
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

// ErrKeyNotFound is returned when a key is not found in the configuration.
var ErrKeyNotFound = errors.New("key not found")

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