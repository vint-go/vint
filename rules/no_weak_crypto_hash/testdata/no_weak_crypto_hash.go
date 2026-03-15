package fixtures

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
)

// Invalid: MD5 Sum is cryptographically broken
func hashPasswordMd5(password string) []byte {
	h := md5.Sum([]byte(password)) // MATCH /use of weak cryptographic hash function md5/
	return h[:]
}

// Invalid: SHA1 New is cryptographically weak
func signDataSha1(data []byte) []byte {
	h := sha1.New() // MATCH /use of weak cryptographic hash function sha1/
	h.Write(data)
	return h.Sum(nil)
}

// Invalid: MD5 New
func hashWithMd5New(data []byte) []byte {
	h := md5.New() // MATCH /use of weak cryptographic hash function md5/
	h.Write(data)
	return h.Sum(nil)
}

// Invalid: SHA1 Sum
func hashWithSha1Sum(data []byte) [20]byte {
	return sha1.Sum(data) // MATCH /use of weak cryptographic hash function sha1/
}

// Valid: SHA-256 is a secure hash function
func hashDataSha256(data []byte) []byte {
	h := sha256.Sum256(data)
	return h[:]
}

// Valid: SHA-512 is a secure hash function
func hashDataSha512(data []byte) []byte {
	h := sha512.New()
	h.Write(data)
	return h.Sum(nil)
}
