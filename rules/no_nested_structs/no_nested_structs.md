---
title: noNestedStructs
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noNestedStructs`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noNestedStructs:
    # no configuration options
```

## Details

Packages declaring structs that contain other inline struct definitions can be hard to understand/read for other developers. This rule flags nested (inline) struct declarations within struct fields, encouraging the use of named types instead.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
type Foo struct {
	Bar struct { // nested struct not allowed
		Baz string
	}
}
```

```golang
type Bad struct {
	Field []struct{} // nested struct in slice not allowed
}
```

### Valid

```golang
type Bar struct {
	Baz string
}

type Foo struct {
	Bar Bar
}
```

```golang
type issue744 struct {
	c chan struct{} // empty struct in chan is allowed
}
```
