---
title: useTypeDefFirst
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useTypeDefFirst`
- This rule is not enabled by default and must be explicitly enabled.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useTypeDefFirst:
    # no additional options
```

## Details

Detects method declarations preceding the type definition itself. This checker identifies violations of Go code organization where methods are declared before their associated type is defined. It recommends that type definitions should precede their method implementations for better code readability.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
func (s MyStruct) Process() error {
    return nil
}

type MyStruct struct { // type defined after its methods
    Name string
}
```

### Valid

```golang
type MyStruct struct { // type defined first
    Name string
}

func (s MyStruct) Process() error {
    return nil
}
```
