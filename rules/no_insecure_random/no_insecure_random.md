---
title: noInsecureRandom
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noInsecureRandom`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noInsecureRandom:
    # rule options here
```

## Details

Detects the use of insecure random number sources from the `math/rand` package.

This rule flags calls to `math/rand` functions such as `rand.Int`, `rand.Intn`, `rand.Float64`, `rand.Read`, and others. The `math/rand` package uses a deterministic pseudo-random number generator (PRNG) that is not suitable for security-sensitive purposes. Its output can be predicted if the seed is known, and the default seed produces the same sequence across runs (prior to Go 1.20).

For security-sensitive operations such as generating tokens, passwords, cryptographic keys, session IDs, or nonces, use `crypto/rand` instead, which provides cryptographically secure random numbers from the operating system's entropy source.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
import "math/rand"

func generateToken() string {
    // Insecure random number generator
    token := make([]byte, 32)
    for i := range token {
        token[i] = byte(rand.Intn(256))
    }
    return hex.EncodeToString(token)
}
```

```golang
import "math/rand"

func generateSessionID() string {
    // Predictable random values
    return fmt.Sprintf("%d", rand.Int63())
}
```

### Valid

```golang
import "crypto/rand"

func generateToken() (string, error) {
    // Cryptographically secure random
    token := make([]byte, 32)
    _, err := rand.Read(token)
    if err != nil {
        return "", err
    }
    return hex.EncodeToString(token), nil
}
```

```golang
import "crypto/rand"
import "math/big"

func generateSessionID() (string, error) {
    n, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
    if err != nil {
        return "", err
    }
    return n.String(), nil
}
```

```golang
import "math/rand"

// Non-security use of math/rand is acceptable (e.g., shuffling, sampling)
func shuffleSlice(s []int) {
    rand.Shuffle(len(s), func(i, j int) {
        s[i], s[j] = s[j], s[i]
    })
}
```
