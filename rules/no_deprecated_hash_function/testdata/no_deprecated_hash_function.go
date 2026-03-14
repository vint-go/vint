package fixtures

import (
	"crypto/sha256"

	"golang.org/x/crypto/md4"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
)

// Invalid: using MD4 deprecated hash function
func hashWithMD4(data []byte) []byte {
	h := md4.New() // MATCH /use of deprecated hash function md4.New: MD4 is cryptographically broken/
	h.Write(data)
	return h.Sum(nil)
}

// Invalid: using RIPEMD160 deprecated hash function
func hashWithRIPEMD160(data []byte) []byte {
	h := ripemd160.New() // MATCH /use of deprecated hash function ripemd160.New: RIPEMD160 has insufficient collision resistance/
	h.Write(data)
	return h.Sum(nil)
}

// Valid: using SHA-256
func hashWithSHA256(data []byte) []byte {
	h := sha256.Sum256(data)
	return h[:]
}

// Valid: using SHA-3
func hashWithSHA3(data []byte) []byte {
	h := sha3.New256()
	h.Write(data)
	return h.Sum(nil)
}
