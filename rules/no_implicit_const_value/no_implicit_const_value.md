---
title: noImplicitConstValue
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noImplicitConstValue`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noImplicitConstValue:
    # rule options here
```

## Details

Only the first constant in a `const` group has an explicit value, the rest are implicitly repeated.

In a `const` block, if only the first constant has an explicit type and value, subsequent constants inherit the same value. This is different from `iota`-based patterns and is often unintentional.

Source: https://staticcheck.dev/docs/checks/#SA9004

## Examples

### Invalid

```golang
package main

const (
    A = 1
    B    // B = 1, probably meant to be different
    C    // C = 1, probably meant to be different
)
```

### Valid

```golang
package main

const (
    A = 1
    B = 2
    C = 3
)
```

```golang
package main

const (
    A = iota
    B
    C
)
```
