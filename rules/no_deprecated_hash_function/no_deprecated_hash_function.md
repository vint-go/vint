---
title: noDeprecatedHashFunction
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noDeprecatedHashFunction`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noDeprecatedHashFunction:
    # rule options here
```

## Details

Detects the usage of deprecated hash functions MD4 and RIPEMD160.

This rule flags calls to `golang.org/x/crypto/md4.New` and `golang.org/x/crypto/ripemd160.New`. Both MD4 and RIPEMD160 are deprecated cryptographic hash functions that should not be used in new applications.

- **MD4** is severely broken and has practical collision attacks. It was deprecated by RFC 6150 in 2011.
- **RIPEMD160** produces a 160-bit hash that provides insufficient collision resistance by modern standards.

Applications should use SHA-256 or SHA-3 family hash functions for cryptographic hashing needs.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
import "golang.org/x/crypto/md4"

func hashData(data []byte) []byte {
    // MD4 is severely broken
    h := md4.New()
    h.Write(data)
    return h.Sum(nil)
}
```

```golang
import "golang.org/x/crypto/ripemd160"

func hashData(data []byte) []byte {
    // RIPEMD160 is deprecated
    h := ripemd160.New()
    h.Write(data)
    return h.Sum(nil)
}
```

### Valid

```golang
import "crypto/sha256"

func hashData(data []byte) []byte {
    h := sha256.Sum256(data)
    return h[:]
}
```

```golang
import "golang.org/x/crypto/sha3"

func hashData(data []byte) []byte {
    h := sha3.New256()
    h.Write(data)
    return h.Sum(nil)
}
```
