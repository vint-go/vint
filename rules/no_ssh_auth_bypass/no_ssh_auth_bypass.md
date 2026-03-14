---
title: noSshAuthBypass
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noSshAuthBypass`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noSshAuthBypass:
    # rule options here
```

## Details

Detects stateful misuse of `ssh.PublicKeyCallback` that can lead to authentication bypass.

This rule uses SSA analysis to identify patterns where the `PublicKeyCallback` in `ssh.ServerConfig` is used in a way that allows authentication bypass. In certain SSH server implementations, the `PublicKeyCallback` is called to check if a public key is acceptable, but the callback does not perform the actual signature verification (which is handled by the SSH protocol). If server code makes authorization decisions based solely on which key was presented in the callback without verifying that the subsequent authentication step completed successfully, an attacker can present any public key and bypass authentication.

The SSH library's `PublicKeyCallback` should only be used to check if a key is in the authorized keys list. Authorization decisions should be deferred until after the full authentication handshake completes.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
import "golang.org/x/crypto/ssh"

func setupServer() *ssh.ServerConfig {
    config := &ssh.ServerConfig{
        PublicKeyCallback: func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
            // Storing state based on the key without verifying signature
            // This key might not actually be owned by the connecting user
            if isAuthorizedKey(key) {
                return &ssh.Permissions{
                    Extensions: map[string]string{
                        "role": "admin", // Granting role based on unverified key
                    },
                }, nil
            }
            return nil, fmt.Errorf("unauthorized")
        },
    }
    return config
}
```

### Valid

```golang
import "golang.org/x/crypto/ssh"

func setupServer() *ssh.ServerConfig {
    config := &ssh.ServerConfig{
        PublicKeyCallback: func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
            // Only check if the key is in the authorized list
            // Do not assign roles or permissions here
            if isAuthorizedKey(key) {
                return &ssh.Permissions{}, nil
            }
            return nil, fmt.Errorf("unauthorized")
        },
    }
    return config
}

// Assign roles after full authentication is verified
func handleConnection(conn *ssh.ServerConn) {
    user := conn.User()
    role := lookupRole(user)
    // ...
}
```
