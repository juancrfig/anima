package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/argon2"
)


// Params defines the tunable parameters for Argon2id and AES
type Params struct {
	Time    uint32 // Argon2id time cost
	Memory  uint32 // Argon2id memory cost (in KiB)
	Threads uint8  // Argon2id parallelism
	SaltLen uint8  // Length of the salt
	KeyLen  uint8  // Length of the derived key (e.g., 32 for AES-256)
}



var ErrDecryption = errors.New("crypto: decryption failed")

// The fixed size (in bytes) of our metadata header.
const headerSize = 4 + 4 + 1 + 1 + 1 + 1


// Uses Argon2id to derive a key from a password and salt
func deriveKey(password, salt []byte, p *Params) []byte {
	return argon2.IDKey(password, salt, p.Time, p.Memory, p.Threads, uint32(p.KeyLen))
}

// Encrypt derives a key and encrypts the plaintext.
// The output is a single blob: [header][salt][nonce][ciphertext]
func Encrypt(plaintext, password []byte, p *Params) ([]byte, error) {
	// We only support AES-256 (32 bytes).
	if p.KeyLen != 32 {
		return nil, errors.New("crypto: invalid KeyLen, must be 32 for AES-256")
	}

	// 1. Generate a new random salt
	salt := make([]byte, p.SaltLen)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}

	// 2. Derive the encryption key from the password and salt
	key := deriveKey(password, salt, p)

	// 3. Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 4. Create the GCM cipher mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 5. Generate a new random nonce (Number used Once)
	nonceSize := gcm.NonceSize()
	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// 6. Encrypt the data
	// gcm.Seal appends the ciphertext and auth tag to the 'nil' first arg
	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	// 7. Construct the header ("recipe card")
	header := make([]byte, headerSize)
	binary.BigEndian.PutUint32(header[0:4], p.Time)
	binary.BigEndian.PutUint32(header[4:8], p.Memory)
	header[8] = p.Threads
	header[9] = p.SaltLen
	header[10] = uint8(nonceSize)
	header[11] = p.KeyLen

	// 8. Combine all parts into the final blob
	finalData := make([]byte, 0, headerSize+len(salt)+len(nonce)+len(ciphertext))
	finalData = append(finalData, header...)
	finalData = append(finalData, salt...)
	finalData = append(finalData, nonce...)
	finalData = append(finalData, ciphertext...)

	return finalData, nil
}

// Decrypt reads the self-describing blob, re-derives the key,
// and decrypts the data.
func Decrypt(data, password []byte) ([]byte, error) {
	// 1. Check if data is at least as long as our header
	if len(data) < headerSize {
		return nil, fmt.Errorf("crypto: data is too short: %w", ErrDecryption)
	}

	// 2. Parse the header to get the "recipe"
	header := data[:headerSize]
	p := &Params{
		Time:    binary.BigEndian.Uint32(header[0:4]),
		Memory:  binary.BigEndian.Uint32(header[4:8]),
		Threads: header[8],
		SaltLen: header[9],
		KeyLen:  header[11],
	}
	nonceSize := int(header[10])

	// 3. Check boundaries to prevent panics from corrupt data
	saltStart := headerSize
	saltEnd := saltStart + int(p.SaltLen)
	if saltEnd > len(data) {
		return nil, fmt.Errorf("crypto: corrupt data, invalid salt boundary: %w", ErrDecryption)
	}

	nonceStart := saltEnd
	nonceEnd := nonceStart + nonceSize
	if nonceEnd > len(data) {
		return nil, fmt.Errorf("crypto: corrupt data, invalid nonce boundary: %w", ErrDecryption)
	}

	// 4. Extract components from the blob
	salt := data[saltStart:saltEnd]
	nonce := data[nonceStart:nonceEnd]
	ciphertext := data[nonceEnd:]

	// 5. Re-derive the *exact same key* using the stored salt
	key := deriveKey(password, salt, p)

	// 6. Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("crypto: cipher init failed: %w", ErrDecryption)
	}

	// 7. Create the GCM cipher mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("crypto: gcm init failed: %w", ErrDecryption)
	}

	// 8. Sanity check: ensure stored nonce size matches cipher's need
	if nonceSize != gcm.NonceSize() {
		return nil, fmt.Errorf("crypto: nonce size mismatch: %w", ErrDecryption)
	}

	// 9. Decrypt and authenticate
	// gcm.Open checks the authentication tag *first*. If the password
	// is wrong (so key is wrong) or data is corrupt, this will fail.
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, ErrDecryption
	}

	return plaintext, nil
}
