---
title: noInappropriateContextKey
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noInappropriateContextKey`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noInappropriateContextKey:
    # rule options here
```

## Details

Inappropriate key in call to `context.WithValue`.

The key used with `context.WithValue` should not be a built-in type like `string` or `int` to avoid collisions between packages. Instead, define an unexported type and use a value of that type as the key.

Source: https://staticcheck.dev/docs/checks/#SA1029

## Examples

### Invalid

```golang
package main

import "context"

func main() {
    // Wrong: using string as context key
    ctx := context.WithValue(context.Background(), "key", "value")
    _ = ctx
}
```

### Valid

```golang
package main

import "context"

type contextKey string

const myKey contextKey = "key"

func main() {
    // Correct: using custom type as context key
    ctx := context.WithValue(context.Background(), myKey, "value")
    _ = ctx
}
```
