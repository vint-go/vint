---
title: noWeakEncryptionAlgorithm
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noWeakEncryptionAlgorithm`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noWeakEncryptionAlgorithm:
    # rule options here
```

## Details

Detects the usage of the DES or RC4 encryption algorithms.

This rule flags calls to `crypto/des.NewCipher`, `crypto/des.NewTripleDESCipher`, and `crypto/rc4.NewCipher`. Both DES and RC4 are considered cryptographically broken and should not be used.

- **DES** uses a 56-bit key which can be brute-forced in hours with modern hardware. Triple DES (3DES) is also deprecated due to the Sweet32 attack.
- **RC4** has multiple known biases in its output stream and is vulnerable to practical attacks. It has been prohibited in TLS since RFC 7465.

Applications should use AES (Advanced Encryption Standard) with appropriate modes of operation (GCM, CTR) for symmetric encryption. AES provides strong security with 128, 192, or 256-bit keys.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
import "crypto/des"

func encrypt(key, plaintext []byte) ([]byte, error) {
    // DES is cryptographically broken
    block, err := des.NewCipher(key)
    if err != nil {
        return nil, err
    }
    // ...
}
```

```golang
import "crypto/rc4"

func encrypt(key []byte) (*rc4.Cipher, error) {
    // RC4 is cryptographically broken
    return rc4.NewCipher(key)
}
```

```golang
import "crypto/des"

func encrypt3DES(key, plaintext []byte) ([]byte, error) {
    // Triple DES is deprecated
    block, err := des.NewTripleDESCipher(key)
    if err != nil {
        return nil, err
    }
    // ...
}
```

### Valid

```golang
import "crypto/aes"
import "crypto/cipher"

func encrypt(key, plaintext []byte) ([]byte, error) {
    // AES-GCM provides strong authenticated encryption
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
```
