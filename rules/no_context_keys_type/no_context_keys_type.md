---
title: noContextKeysType
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noContextKeysType`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noContextKeysType:
    # rule options here
```

## Details

Basic types should not be used as a key in `context.WithValue`.

Using basic types (like `string`, `int`, etc.) as context keys can cause collisions between packages. Instead, define an unexported custom type and use a value of that type as the key.

Source: https://github.com/mgechev/revive

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

```golang
package main

import "context"

func main() {
    // Wrong: using int as context key
    ctx := context.WithValue(context.Background(), 42, "value")
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
