package fixtures

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
)

// Invalid: 1024-bit RSA key is too short
func generateShortKey1024() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 1024) // MATCH /RSA key length 1024 is too short, minimum 2048 bits recommended/
}

// Invalid: 512-bit RSA key is trivially breakable
func generateShortKey512() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 512) // MATCH /RSA key length 512 is too short, minimum 2048 bits recommended/
}

// Valid: 2048-bit RSA key meets minimum requirement
func generateKey2048() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 2048)
}

// Valid: 4096-bit RSA key for long-term security
func generateKey4096() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 4096)
}

// Valid: ECDSA key generation is not affected
func generateECDSAKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}
