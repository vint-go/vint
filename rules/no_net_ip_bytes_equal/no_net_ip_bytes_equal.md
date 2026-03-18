---
title: noNetIpBytesEqual
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noNetIpBytesEqual`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noNetIpBytesEqual:
    # rule options here
```

## Details

Using `bytes.Equal` to compare two `net.IP` values.

`net.IP` values should be compared using the `net.IP.Equal` method, not `bytes.Equal`. IPv4 addresses can be represented as both 4-byte and 16-byte slices, and `bytes.Equal` would incorrectly report them as different.

Source: https://staticcheck.dev/docs/checks/#SA1021

## Examples

### Invalid

```golang
package main

import (
    "bytes"
    "net"
)

func main() {
    ip1 := net.ParseIP("127.0.0.1")
    ip2 := net.IPv4(127, 0, 0, 1)
    // Wrong: different byte representations may not match
    if bytes.Equal(ip1, ip2) {
        // ...
    }
}
```

### Valid

```golang
package main

import "net"

func main() {
    ip1 := net.ParseIP("127.0.0.1")
    ip2 := net.IPv4(127, 0, 0, 1)
    // Correct: handles different representations
    if ip1.Equal(ip2) {
        // ...
    }
}
```
