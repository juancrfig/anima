package crypto

import (
	"bytes"
	"testing"
)

// DefaultParams provides the OWASP-recommended settings
var DefaultParams = &Params{
	Time:    3,
	Memory:  65536, // 64 MB
	Threads: 1,
	SaltLen: 16, // 16 bytes
	KeyLen:  32, // 32 bytes for AES-256
}


// It proves that data encrypted with the default params can be successfully decrypted
func TestEncryptDecrypt_Roundtrip(t *testing.T) {
	password := []byte("my-super-secret-password")
	plaintext := []byte("anima protocol: standing by.")

	// 1. Encrypt the data using default params
	ciphertext, err := Encrypt(plaintext, password, DefaultParams)
	if err != nil {
		t.Fatalf("Encrypt() failed: %v", err)
	}

	if bytes.Equal(plaintext, ciphertext) {
		t.Fatal("Encrypt() returned plaintext, encryption failed.")
	}

	// 2. Decrypt the data
	// Notice it doesn't need params, it finds them in the ciphertext's header
	decrypted, err := Decrypt(ciphertext, password)
	if err != nil {
		t.Fatalf("Decrypt() failed: %v", err)
	}

	// 3. Verify
	if !bytes.Equal(plaintext, decrypted) {
		t.Fatalf("Decrypt() output does not match plaintext. got=%q, want=%q", decrypted, plaintext)
	}
}


// Proves the header logic is working
// We use weaker, faster params just for this test
// Decrypt() should correctly read the header and handle data not created with DefaultParams
func TestEncryptDecrypt_CustomParams(t *testing.T) {
	customParams := &Params{
		Time:    1,
		Memory:  1024, // 1 KB (fast for testing)
		Threads: 1,
		SaltLen: 32, // 32-byte salt (non-default)
		KeyLen:  32, // AES-256
	}

	password := []byte("custom-password")
	plaintext := []byte("data encrypted with different params")

	// 1. Encrypt with custom params
	ciphertext, err := Encrypt(plaintext, password, customParams)
	if err != nil {
		t.Fatalf("Encrypt() failed: %v", err)
	}

	// 2. Decrypt
	// Decrypt must discover the custom params from the header
	decrypted, err := Decrypt(ciphertext, password)
	if err != nil {
		t.Fatalf("Decrypt() failed: %v", err)
	}

	// 3. Verify
	if !bytes.Equal(plaintext, decrypted) {
		t.Fatal("Decrypt() output does not match plaintext for custom params.")
	}
}


// Proves that decryption fails if the password is incorrect.
// This MUST return our specific ErrDecryption
func TestDecrypt_WrongPassword(t *testing.T) {
	password := []byte("correct-password")
	wrongPassword := []byte("wrong-password")
	plaintext := []byte("secret data")

	ciphertext, err := Encrypt(plaintext, password, DefaultParams)
	if err != nil {
		t.Fatalf("Encrypt() failed: %v", err)
	}

	// Attempt to decrypt with the wrong password
	_, err = Decrypt(ciphertext, wrongPassword)
	if err == nil {
		t.Fatal("Decrypt() succeeded with wrong password, expected an error.")
	}

	// We must get our specific, opaque error
	if err != ErrDecryption {
		t.Fatalf("Decrypt() returned wrong error. got=%v, want=%v", err, ErrDecryption)
	}
}


// Proves that AES-GCM's authentication works. If any bit of the data (salt, nonce, or ciphertext)
// is changed, decryption MUST fail
func TestDecrypt_CorruptedData(t *testing.T) {
	password := []byte("correct-password")
	plaintext := []byte("secret data")

	ciphertext, err := Encrypt(plaintext, password, DefaultParams)
	if err != nil {
		t.Fatalf("Encrypt() failed: %v", err)
	}

	// Corrupt the data (flip a bit at the end)
	ciphertext[len(ciphertext)-1] = ^ciphertext[len(ciphertext)-1]

	_, err = Decrypt(ciphertext, password)
	if err == nil {
		t.Fatal("Decrypt() succeeded with corrupted data, expected an error.")
	}

	if err != ErrDecryption {
		t.Fatalf("Decrypt() returned wrong error. got=%v, want=%v", err, ErrDecryption)
	}
}

// TestDecrypt_InvalidHeader proves our header parsing is robust.
// It checks two failure modes:
// 1. Data that is too short to even contain a header.
// 2. Data with a "corrupt" header (e.g., claims salt is 255 bytes
//    but the data blob is only 20 bytes long).
func TestDecrypt_InvalidHeader(t *testing.T) {
	password := []byte("password")

	// 1. Test with data that is too short (less than headerSize)
	shortData := []byte{0x01, 0x02, 0x03}
	_, err := Decrypt(shortData, password)
	if err == nil {
		t.Fatal("Decrypt() succeeded with short data, expected error.")
	}
	t.Logf("Got expected short data error: %v", err)

	// 2. Test with valid length but corrupt header
	// This format is [time(4b)|mem(4b)|threads(1b)|saltLen(1b)|nonceSize(1b)|keyLen(1b)]
	// We'll set saltLen to 255, which is longer than the remaining data.
	header := []byte{
		0, 0, 0, 1, // time
		0, 0, 4, 0, // mem
		1,   // threads
		255, // saltLen (maliciously large)
		12,  // nonceSize
		32,  // keyLen
	}
	corruptData := append(header, []byte("some-data")...) // data < 255
	_, err = Decrypt(corruptData, password)
	if err == nil {
		t.Fatal("Decrypt() succeeded with corrupt header, expected error.")
	}
	t.Logf("Got expected corrupt header error: %v", err)
}
