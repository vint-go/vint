---
title: noAddressOfDereference
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noAddressOfDereference`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noAddressOfDereference:
    # rule options here
```

## Details

`&*x` gets simplified to `x`, it does not copy `x`.

Taking the address of a dereferenced pointer (`&*x`) is the same as `x` itself. It does not create a copy of the pointed-to value. If you intended to copy the value, you should dereference and then take the address of a local variable.

Source: https://staticcheck.dev/docs/checks/#SA4001

## Examples

### Invalid

```golang
package main

type T struct{ Name string }

func copy(t *T) *T {
    // Does not copy: &*t is the same as t
    return &*t
}
```

### Valid

```golang
package main

type T struct{ Name string }

func copyT(t *T) *T {
    // Creates a real copy
    v := *t
    return &v
}
```
