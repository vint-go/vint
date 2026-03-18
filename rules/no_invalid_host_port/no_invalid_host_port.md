---
title: noInvalidHostPort
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noInvalidHostPort`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noInvalidHostPort:
    # rule options here
```

## Details

Using an invalid host:port pair with a `net.Listen`-related function.

This check validates that the address argument to functions like `net.Listen`, `net.Dial`, and `http.ListenAndServe` is a valid host:port pair. Common mistakes include forgetting the colon, using an invalid port number, or using the wrong format.

Source: https://staticcheck.dev/docs/checks/#SA1020

## Examples

### Invalid

```golang
package main

import "net/http"

func main() {
    // Missing colon before port
    http.ListenAndServe("localhost8080", nil)
}
```

### Valid

```golang
package main

import "net/http"

func main() {
    // Correct host:port format
    http.ListenAndServe(":8080", nil)
}
```
