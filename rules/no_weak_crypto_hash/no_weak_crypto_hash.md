---
title: noWeakCryptoHash
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noWeakCryptoHash`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noWeakCryptoHash:
    # rule options here
```

## Details

Detects the usage of weak cryptographic hash functions MD5 or SHA1.

This rule flags calls to `crypto/md5.New`, `crypto/md5.Sum`, `crypto/sha1.New`, and `crypto/sha1.Sum`. Both MD5 and SHA1 are considered cryptographically broken and should not be used for security purposes such as password hashing, digital signatures, certificate validation, or integrity verification of security-sensitive data.

MD5 is vulnerable to collision attacks and can be exploited in practical scenarios. SHA1 has known theoretical weaknesses and demonstrated collision attacks (SHAttered). Applications should use SHA-256 or SHA-3 for cryptographic hashing, and specialized password hashing algorithms (bcrypt, scrypt, argon2) for password storage.

Note: Non-security uses of MD5/SHA1 (such as checksums for data deduplication) may be acceptable but should be reviewed.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
import "crypto/md5"

func hashPassword(password string) []byte {
    // MD5 is cryptographically broken
    h := md5.Sum([]byte(password))
    return h[:]
}
```

```golang
import "crypto/sha1"

func signData(data []byte) []byte {
    // SHA1 is cryptographically weak
    h := sha1.New()
    h.Write(data)
    return h.Sum(nil)
}
```

### Valid

```golang
import "crypto/sha256"

func hashData(data []byte) []byte {
    // SHA-256 is a secure hash function
    h := sha256.Sum256(data)
    return h[:]
}
```

```golang
import "golang.org/x/crypto/bcrypt"

func hashPassword(password string) ([]byte, error) {
    // bcrypt is appropriate for password hashing
    return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
```

```golang
import "crypto/sha512"

func hashData(data []byte) []byte {
    h := sha512.New()
    h.Write(data)
    return h.Sum(nil)
}
```
