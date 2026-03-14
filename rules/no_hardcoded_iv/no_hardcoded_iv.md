---
title: noHardcodedIv
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noHardcodedIv`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noHardcodedIv:
    # rule options here
```

## Details

Detects the use of hardcoded initialization vectors (IVs) or nonces for encryption.

This rule uses SSA analysis to identify cases where a constant or hardcoded byte slice is used as an IV or nonce for symmetric encryption operations. Using a static IV or nonce with the same key completely undermines the security of encryption schemes like AES-GCM, AES-CTR, and AES-CBC.

- In **GCM mode**, reusing a nonce with the same key allows an attacker to recover the authentication key and forge messages.
- In **CTR mode**, reusing a nonce with the same key reveals the XOR of the plaintexts.
- In **CBC mode**, reusing an IV with the same key leaks information about common prefixes.

IVs and nonces should always be generated using a cryptographically secure random number generator (`crypto/rand`), and must never be reused with the same key.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
import (
    "crypto/aes"
    "crypto/cipher"
)

func encrypt(key, plaintext []byte) ([]byte, error) {
    block, _ := aes.NewCipher(key)
    gcm, _ := cipher.NewGCM(block)

    // Hardcoded nonce - never do this
    nonce := []byte("hardcodednonce")
    return gcm.Seal(nil, nonce, plaintext, nil), nil
}
```

```golang
import (
    "crypto/aes"
    "crypto/cipher"
)

// Constant IV
var iv = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

func encryptCBC(key, plaintext []byte) ([]byte, error) {
    block, _ := aes.NewCipher(key)
    mode := cipher.NewCBCEncrypter(block, iv) // Static IV
    // ...
}
```

### Valid

```golang
import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "io"
)

func encrypt(key, plaintext []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }

    // Generate random nonce
    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, err
    }
    return gcm.Seal(nonce, nonce, plaintext, nil), nil
}
```
