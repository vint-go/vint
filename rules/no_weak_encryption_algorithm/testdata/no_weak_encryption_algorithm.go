package fixtures

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/rc4"
	"crypto/rand"
	"io"
)

// Invalid: DES NewCipher
func encryptDES(key, plaintext []byte) ([]byte, error) {
	block, err := des.NewCipher(key) // MATCH /use of weak encryption algorithm DES: use AES instead/
	if err != nil {
		return nil, err
	}
	_ = block
	return nil, nil
}

// Invalid: RC4 NewCipher
func encryptRC4(key []byte) (*rc4.Cipher, error) {
	return rc4.NewCipher(key) // MATCH /use of weak encryption algorithm RC4: use AES instead/
}

// Invalid: Triple DES NewTripleDESCipher
func encrypt3DES(key, plaintext []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key) // MATCH /use of weak encryption algorithm 3DES: use AES instead/
	if err != nil {
		return nil, err
	}
	_ = block
	return nil, nil
}

// Valid: AES-GCM provides strong authenticated encryption
func encryptAES(key, plaintext []byte) ([]byte, error) {
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
