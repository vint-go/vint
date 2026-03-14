package fixtures

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

// Invalid: hardcoded nonce passed to Seal via variable
func encryptGCMHardcoded(key, plaintext []byte) ([]byte, error) {
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)
	nonce := []byte("hardcodednonce")
	return gcm.Seal(nil, nonce, plaintext, nil), nil // MATCH /hardcoded IV or nonce: use a cryptographically random value instead/
}

// Invalid: hardcoded IV in CBC encrypter
func encryptCBCHardcoded(key, plaintext []byte) {
	block, _ := aes.NewCipher(key)
	iv := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	cipher.NewCBCEncrypter(block, iv) // MATCH /hardcoded IV or nonce: use a cryptographically random value instead/
}

// Invalid: hardcoded IV directly in NewCBCEncrypter call
func encryptCBCDirectHardcoded(key, plaintext []byte) {
	block, _ := aes.NewCipher(key)
	cipher.NewCBCEncrypter(block, []byte("0123456789abcdef")) // MATCH /hardcoded IV or nonce: use a cryptographically random value instead/
}

// Invalid: hardcoded IV in CTR mode
func encryptCTRHardcoded(key, plaintext []byte) {
	block, _ := aes.NewCipher(key)
	iv := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	cipher.NewCTR(block, iv) // MATCH /hardcoded IV or nonce: use a cryptographically random value instead/
}

// Invalid: hardcoded IV in OFB mode
func encryptOFBHardcoded(key, plaintext []byte) {
	block, _ := aes.NewCipher(key)
	cipher.NewOFB(block, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}) // MATCH /hardcoded IV or nonce: use a cryptographically random value instead/
}

// Invalid: hardcoded IV in CFB encrypter
func encryptCFBHardcoded(key, plaintext []byte) {
	block, _ := aes.NewCipher(key)
	cipher.NewCFBEncrypter(block, []byte("abcdefghijklmnop")) // MATCH /hardcoded IV or nonce: use a cryptographically random value instead/
}

// Invalid: package-level hardcoded IV used in function
var staticIV = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

func encryptWithStaticIV(key, plaintext []byte) {
	block, _ := aes.NewCipher(key)
	cipher.NewCBCEncrypter(block, staticIV) // MATCH /hardcoded IV or nonce: use a cryptographically random value instead/
}

// Valid: random nonce for GCM
func encryptGCMRandom(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// Valid: random IV for CBC
func encryptCBCRandom(key, plaintext []byte) {
	block, _ := aes.NewCipher(key)
	iv := make([]byte, aes.BlockSize)
	_, _ = io.ReadFull(rand.Reader, iv)
	cipher.NewCBCEncrypter(block, iv)
}

// Valid: IV passed as function parameter (not hardcoded)
func encryptCBCParam(key, plaintext, iv []byte) {
	block, _ := aes.NewCipher(key)
	cipher.NewCBCEncrypter(block, iv)
}
