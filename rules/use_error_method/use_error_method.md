---
title: useErrorMethod
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useErrorMethod`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useErrorMethod:
    # rule options here
```

## Details

Simplify `fmt.Sprintf("%s", err)` to `err.Error()`.

When formatting an error value with `%s`, it is cleaner to call `err.Error()` directly.

Source: https://staticcheck.dev/docs/checks/#S1028

## Examples

### Invalid

```golang
package main

import "fmt"

func getMessage(err error) string {
    return fmt.Sprintf("%s", err)
}
```

### Valid

```golang
package main

func getMessage(err error) string {
    return err.Error()
}
```
