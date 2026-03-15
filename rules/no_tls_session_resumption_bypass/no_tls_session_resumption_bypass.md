---
title: noTlsSessionResumptionBypass
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noTlsSessionResumptionBypass`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noTlsSessionResumptionBypass:
    # rule options here
```

## Details

Detects TLS session resumption configurations that bypass `VerifyPeerCertificate` when `VerifyConnection` is not set.

This rule uses SSA analysis to identify `crypto/tls.Config` instances where `VerifyPeerCertificate` is set but `VerifyConnection` is not. When TLS session resumption occurs, the `VerifyPeerCertificate` callback is not invoked for resumed sessions because no new certificate exchange takes place. If `VerifyConnection` is not configured, custom certificate validation logic in `VerifyPeerCertificate` is silently skipped during session resumption, potentially allowing connections that would otherwise be rejected.

To ensure consistent certificate validation across both new and resumed TLS sessions, either set `VerifyConnection` alongside `VerifyPeerCertificate`, or disable session resumption by setting `SessionTicketsDisabled` to `true`.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
import "crypto/tls"

// VerifyPeerCertificate without VerifyConnection
tlsConfig := &tls.Config{
    VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
        // Custom validation logic - NOT called on session resumption
        return validateCert(rawCerts)
    },
}
```

### Valid

```golang
import "crypto/tls"

// Both VerifyPeerCertificate and VerifyConnection set
tlsConfig := &tls.Config{
    VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
        return validateCert(rawCerts)
    },
    VerifyConnection: func(state tls.ConnectionState) error {
        // Called for both new and resumed sessions
        return validateConnection(state)
    },
}
```

```golang
import "crypto/tls"

// Disable session resumption when only using VerifyPeerCertificate
tlsConfig := &tls.Config{
    SessionTicketsDisabled: true,
    VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
        return validateCert(rawCerts)
    },
}
```
