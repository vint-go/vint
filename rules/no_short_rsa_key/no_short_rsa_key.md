---
title: noShortRsaKey
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noShortRsaKey`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noShortRsaKey:
    # rule options here
```

## Details

Ensures minimum RSA key length of 2048 bits.

This rule detects calls to `crypto/rsa.GenerateKey` where the key size is less than 2048 bits. RSA keys shorter than 2048 bits are considered insecure and can be factored with modern computing resources. NIST and other standards bodies recommend a minimum of 2048 bits for RSA keys, with 3072 or 4096 bits recommended for long-term security.

Short RSA keys put encrypted data and digital signatures at risk of being compromised. All new key generation should use at least 2048 bits, and 4096 bits is recommended for new deployments.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
import "crypto/rsa"

func generateKey() (*rsa.PrivateKey, error) {
    // 1024-bit RSA key is too short
    return rsa.GenerateKey(rand.Reader, 1024)
}
```

```golang
import "crypto/rsa"

func generateKey() (*rsa.PrivateKey, error) {
    // 512-bit RSA key is trivially breakable
    return rsa.GenerateKey(rand.Reader, 512)
}
```

### Valid

```golang
import "crypto/rsa"

func generateKey() (*rsa.PrivateKey, error) {
    // 2048-bit RSA key meets minimum requirement
    return rsa.GenerateKey(rand.Reader, 2048)
}
```

```golang
import "crypto/rsa"

func generateKey() (*rsa.PrivateKey, error) {
    // 4096-bit RSA key for long-term security
    return rsa.GenerateKey(rand.Reader, 4096)
}
```

```golang
import "crypto/ecdsa"
import "crypto/elliptic"

func generateKey() (*ecdsa.PrivateKey, error) {
    // Consider using ECDSA for equivalent security with shorter keys
    return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}
```
