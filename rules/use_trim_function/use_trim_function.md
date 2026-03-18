---
title: useTrimFunction
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useTrimFunction`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useTrimFunction:
    # rule options here
```

## Details

Replace manual trimming with `strings.TrimPrefix` or `strings.TrimSuffix`.

Checking `strings.HasPrefix` and then slicing is equivalent to using `strings.TrimPrefix`, which is clearer and less error-prone.

Source: https://staticcheck.dev/docs/checks/#S1017

## Examples

### Invalid

```golang
package main

import "strings"

func process(s string) string {
    if strings.HasPrefix(s, "http://") {
        s = s[len("http://"):]
    }
    return s
}
```

### Valid

```golang
package main

import "strings"

func process(s string) string {
    return strings.TrimPrefix(s, "http://")
}
```
