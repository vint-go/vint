---
title: useAny
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useAny`
- This rule is **not recommended**, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useAny:
    # no configuration options
```

## Details

Since Go 1.18, `interface{}` has an alias: `any`. This rule proposes to replace instances of `interface{}` with `any`.

Using `any` instead of `interface{}` improves readability and aligns with modern Go idioms. The `any` type alias was introduced specifically to make code cleaner and more concise.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
var i interface{}
```

```golang
func foo(a interface{}) {}
```

```golang
m := map[interface{}]string{}
```

```golang
s := []interface{}{}
```

### Valid

```golang
var i any
```

```golang
func foo(a any) {}
```

```golang
m := map[any]string{}
```

```golang
// Non-empty interfaces are not flagged
type Closer interface{ Close() }
```
