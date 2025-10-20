package cli

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

    "anima/internal/auth"
	"anima/internal/config"
	"anima/internal/storage"
)

// Define a custom context key type to avoid collisions.
type contextKey string

const (
	servicesKey contextKey = "services"
	dbStoreKey  contextKey = "dbStore"
)

// Services holds all shared dependencies for our commands.
type Services struct {
	Store  *storage.Storage
	Config *config.Config
    Auth *auth.Manager
    KeyManager *auth.KeyManager
}

// initServices performs the setup logic that was duplicated in every command.
// It initializes the dependencies and stores them in the command's context.
// This version correctly returns (context.Context, error) [2 values]
func initServices(ctx context.Context) (context.Context, error) {
	animaPath, err := GetAnimaPath()
	if err != nil {
		return nil, err
	}
	dbPath := filepath.Join(animaPath, "anima.db")
	configPath := filepath.Join(animaPath, "config.json")

	cfg, err := config.New(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

    authMgr := auth.New()

    cryptoParams, err := cfg.CryptoParams()
    if err != nil {
        return nil, fmt.Errorf("Could not load crypto params:  %w", err)
    }

    keyMgr := auth.NewKeyManager(cryptoParams)

	store, err := storage.New(dbPath, cfg, authMgr)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Create the services struct
	services := &Services{
		Store:  store,
		Config: cfg,
        Auth: authMgr,
        KeyManager: keyMgr,
	}

	// Store the services struct in the context
	newCtx := context.WithValue(ctx, servicesKey, services)
	// Store the store separately for the cleanup hook
	newCtx = context.WithValue(newCtx, dbStoreKey, store)

	return newCtx, nil
}

// GetServices retrieves the shared services from the command's context.
func GetServices(ctx context.Context) (*Services, error) {
	services, ok := ctx.Value(servicesKey).(*Services)
	if !ok {
		return nil, errors.New("unable to retrieve services from context")
	}
	return services, nil
}
