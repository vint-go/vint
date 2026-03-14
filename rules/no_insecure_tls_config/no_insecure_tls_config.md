---
title: noInsecureTlsConfig
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noInsecureTlsConfig`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noInsecureTlsConfig:
    # rule options here
```

## Details

Detects bad TLS connection settings in `crypto/tls.Config`.

This rule examines `tls.Config` struct instances for insecure configurations including:

- **InsecureSkipVerify set to true**: Disables certificate verification, making the connection vulnerable to man-in-the-middle attacks.
- **Weak cipher suites**: Usage of cipher suites that are known to be vulnerable.
- **MinVersion set too low**: Allowing TLS versions below 1.2 enables attacks exploiting known vulnerabilities in older TLS/SSL protocols (POODLE, BEAST, etc.).
- **MaxVersion capped too low**: Explicitly limiting the maximum TLS version may prevent the use of the latest security improvements.

TLS configurations should enforce at least TLS 1.2 as the minimum version, use only strong cipher suites, and always verify server certificates.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
import "crypto/tls"

// Skipping certificate verification
config := &tls.Config{
    InsecureSkipVerify: true,
}
```

```golang
import "crypto/tls"

// Allowing TLS 1.0
config := &tls.Config{
    MinVersion: tls.VersionTLS10,
}
```

```golang
import "crypto/tls"

// Using weak cipher suites
config := &tls.Config{
    CipherSuites: []uint16{
        tls.TLS_RSA_WITH_RC4_128_SHA,
        tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
    },
}
```

### Valid

```golang
import "crypto/tls"

// Secure TLS configuration
config := &tls.Config{
    MinVersion: tls.VersionTLS12,
    CipherSuites: []uint16{
        tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
        tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
        tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
        tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
    },
}
```

```golang
import "crypto/tls"

// Modern TLS 1.3 only configuration
config := &tls.Config{
    MinVersion: tls.VersionTLS13,
}
```
