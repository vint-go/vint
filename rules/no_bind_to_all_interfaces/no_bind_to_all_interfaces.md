---
title: noBindToAllInterfaces
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noBindToAllInterfaces`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noBindToAllInterfaces:
    # rule options here
```

## Details

Detects when a network listener binds to all interfaces (0.0.0.0 or just `:port`).

This rule monitors calls to `net.Listen()` and `crypto/tls.Listen()` and flags addresses that match the pattern `^(0.0.0.0|:).*$`. Binding to all interfaces exposes the service to every network the host is connected to, including potentially untrusted networks.

In production environments, services should bind to specific interfaces to limit their attack surface. Binding to `0.0.0.0` means the service will accept connections from any network interface, which may include public-facing interfaces that should not have access to the service.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
// Binding to all interfaces
listener, err := net.Listen("tcp", "0.0.0.0:8080")
```

```golang
// Binding to all interfaces using shorthand
listener, err := net.Listen("tcp", ":8080")
```

```golang
// TLS listener on all interfaces
listener, err := tls.Listen("tcp", "0.0.0.0:443", tlsConfig)
```

### Valid

```golang
// Binding to localhost only
listener, err := net.Listen("tcp", "127.0.0.1:8080")
```

```golang
// Binding to a specific interface
listener, err := net.Listen("tcp", "192.168.1.10:8080")
```

```golang
// TLS listener on specific interface
listener, err := tls.Listen("tcp", "127.0.0.1:443", tlsConfig)
```
