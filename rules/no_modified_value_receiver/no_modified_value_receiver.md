---
title: noModifiedValueReceiver
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noModifiedValueReceiver`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noModifiedValueReceiver:
    # no options
```

## Details

A method that modifies its receiver value can have undesired behavior.
The modification can also be the root of a bug because the actual value receiver could be a copy of that used at the calling site.
This rule warns when a method modifies its receiver.

The rule skips methods with pointer receivers, anonymous receivers, and receivers whose underlying type is a slice or map (since those are reference types). It also allows modifications when the method returns the receiver value, as this pattern is commonly used for builder-style APIs.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
type data struct {
    num   int
    key   *string
    items map[string]bool
}

func (this data) vmethod() {
    this.num = 8                       // suspicious assignment to a by-value method receiver
    this.items = make(map[string]bool) // suspicious assignment to a by-value method receiver
}

func (this data) incrementDecrement() {
    this.num++ // suspicious assignment to a by-value method receiver
    this.num-- // suspicious assignment to a by-value method receiver
}
```

### Valid

```golang
// Pointer receiver - modifications are visible to the caller
func (this *data) vmethod() {
    this.num = 8
}

// Returning the receiver - builder pattern
func (b JailerCommandBuilder) WithBin(bin string) JailerCommandBuilder {
    b.bin = bin
    return b
}

// Returning a pointer to the receiver
func (a A) Foo() *A {
    a.whatever = true
    return &a
}
```
