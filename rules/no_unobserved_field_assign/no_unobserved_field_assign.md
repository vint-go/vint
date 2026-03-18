---
title: noUnobservedFieldAssign
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noUnobservedFieldAssign`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noUnobservedFieldAssign:
    # rule options here
```

## Details

Field assignment that will never be observed. Did you mean to use a pointer receiver?

Assigning to a field of a struct with a value receiver means the change is made to a copy and will be lost when the method returns. You likely intended to use a pointer receiver.

Source: https://staticcheck.dev/docs/checks/#SA4005

## Examples

### Invalid

```golang
package main

type Counter struct {
    count int
}

// Value receiver: changes to count are lost
func (c Counter) Increment() {
    c.count++
}
```

### Valid

```golang
package main

type Counter struct {
    count int
}

// Pointer receiver: changes are preserved
func (c *Counter) Increment() {
    c.count++
}
```
