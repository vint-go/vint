package fixtures

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	mathrand "math/rand"
	"math"
)

// Invalid: Using math/rand for token generation
func generateToken() string {
	token := make([]byte, 32)
	for i := range token {
		token[i] = byte(mathrand.Intn(256)) // MATCH /use of insecure random number generator (math/rand), use crypto/rand instead/
	}
	return hex.EncodeToString(token)
}

// Invalid: Using math/rand for session ID
func generateSessionID() string {
	return fmt.Sprintf("%d", mathrand.Int63()) // MATCH /use of insecure random number generator (math/rand), use crypto/rand instead/
}

// Invalid: Using math/rand Int
func badInt() int {
	return mathrand.Int() // MATCH /use of insecure random number generator (math/rand), use crypto/rand instead/
}

// Invalid: Using math/rand Float64
func badFloat() float64 {
	return mathrand.Float64() // MATCH /use of insecure random number generator (math/rand), use crypto/rand instead/
}

// Invalid: Using math/rand Read
func badRead() {
	buf := make([]byte, 16)
	mathrand.Read(buf) // MATCH /use of insecure random number generator (math/rand), use crypto/rand instead/
}

// Valid: Using crypto/rand for token generation
func generateTokenSecure() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}

// Valid: Using crypto/rand for session ID
func generateSessionIDSecure() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return "", err
	}
	return n.String(), nil
}

// Valid: Non-security use of math/rand (shuffling)
func shuffleSlice(s []int) {
	mathrand.Shuffle(len(s), func(i, j int) {
		s[i], s[j] = s[j], s[i]
	})
}

// Valid: Seeding math/rand
func seedRandom() {
	mathrand.Seed(42)
}
