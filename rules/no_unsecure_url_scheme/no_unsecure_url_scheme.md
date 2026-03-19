---
title: noUnsecureUrlScheme
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noUnsecureUrlScheme`
- This rule is not recommended. Enable it by explicitly configuring it.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noUnsecureUrlScheme:
    # no configuration options
```

## Details

Checks for usage of potentially unsecure URL schemes (`http`, `ws`) in string literals.
Using unencrypted URL schemes can expose sensitive data during transmission and
make applications vulnerable to man-in-the-middle attacks.
Secure alternatives like `https` and `wss` should be preferred when possible.

The rule will not warn on local URLs (`localhost`, `127.0.0.1`, `0.0.0.0`, `//::` IPv6 addresses).

Test files are automatically skipped.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
var url = "http://example.com"
```

```golang
const urlPattern = "http://%s:%d"
```

```golang
var wsURL = "ws://example.com"
```

### Valid

```golang
var url = "https://example.com"
```

```golang
var wsURL = "wss://example.com"
```

```golang
var localURL = "http://localhost:8080"
```

```golang
var loopbackURL = "http://127.0.0.1:80"
```
