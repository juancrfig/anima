package auth

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"

	"anima/internal/crypto"
	"github.com/tyler-smith/go-bip39"
)

// ErrInvalidCredentials is a specific error returned when decryption fails,
// indicating either a wrong password or a wrong recovery phrase.
var ErrInvalidCredentials = errors.New("auth: invalid credentials")

// KeyManager handles the creation and recovery of master encryption keys.
type KeyManager struct {
	cryptoParams *crypto.Params
}

// SetupResult holds all the critical data generated during the one-time setup.
type SetupResult struct {
	MasterKey          []byte // The raw 32-byte key for AES-256
	RecoveryPhrase     string // The 24-word mnemonic
	EncryptedMasterKey []byte // MasterKey encrypted with the user's password
	EncryptedRecoveryKey []byte // MasterKey encrypted with the RecoveryPhrase
}

// NewKeyManager creates a new key manager with the given crypto parameters.
func NewKeyManager(params *crypto.Params) *KeyManager {
	return &KeyManager{
		cryptoParams: params,
	}
}

// GenerateRecoveryPhrase creates a new 24-word (256-bit) BIP39 mnemonic.
func (km *KeyManager) GenerateRecoveryPhrase() (string, error) {
	// 256 bits of entropy = 24 words
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return "", fmt.Errorf("could not generate entropy: %w", err)
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", fmt.Errorf("could not generate mnemonic: %w", err)
	}
	return mnemonic, nil
}

// Setup performs the initial one-time key generation.
// It creates a master key, a recovery phrase, and encrypts the
// master key with both the password and the phrase.
func (km *KeyManager) Setup(password []byte) (*SetupResult, error) {
	// 1. Generate a new 32-byte master key
	masterKey := make([]byte, 32) // 32 bytes for AES-256
	if _, err := io.ReadFull(rand.Reader, masterKey); err != nil {
		return nil, fmt.Errorf("could not generate master key: %w", err)
	}

	// 2. Generate a new recovery phrase
	recoveryPhrase, err := km.GenerateRecoveryPhrase()
	if err != nil {
		return nil, fmt.Errorf("could not generate recovery phrase: %w", err)
	}
	recoveryPhraseBytes := []byte(recoveryPhrase)

	// 3. Encrypt the master key with the password
	encryptedMasterKey, err := crypto.Encrypt(masterKey, password, km.cryptoParams)
	if err != nil {
		return nil, fmt.Errorf("could not encrypt master key with password: %w", err)
	}

	// 4. Encrypt the master key with the recovery phrase
	encryptedRecoveryKey, err := crypto.Encrypt(masterKey, recoveryPhraseBytes, km.cryptoParams)
	if err != nil {
		return nil, fmt.Errorf("could not encrypt master key with recovery phrase: %w", err)
	}

	return &SetupResult{
		MasterKey:          masterKey,
		RecoveryPhrase:     recoveryPhrase,
		EncryptedMasterKey: encryptedMasterKey,
		EncryptedRecoveryKey: encryptedRecoveryKey,
	}, nil
}

// RecoverMasterKey attempts to decrypt the given data blob using the credential.
// This is used for both 'login' (data=EncryptedMasterKey, credential=password)
// and 'recover' (data=EncryptedRecoveryKey, credential=recoveryPhrase).
func (km *KeyManager) RecoverMasterKey(encryptedKeyData, credential []byte) ([]byte, error) {
	plaintextKey, err := crypto.Decrypt(encryptedKeyData, credential)
	if err != nil {
		// We return our specific, opaque error so we don't leak
		// *why* it failed (e.g., "corrupt data" vs "wrong password").
		return nil, ErrInvalidCredentials
	}
	return plaintextKey, nil
}