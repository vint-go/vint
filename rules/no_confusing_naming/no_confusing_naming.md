---
title: noConfusingNaming
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noConfusingNaming`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noConfusingNaming:
    # no configuration options
```

## Details

Methods or fields of `struct` that have names different only by capitalization could be confusing.

This rule warns when:
- Two methods of the same struct differ only by capitalization (e.g., `aFoo` and `AFoo`).
- Two functions in the same package differ only by capitalization (e.g., `aGlobal` and `AGlobal`).
- Two fields of the same struct differ only by capitalization (e.g., `asd` and `aSd`).

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
type foo struct{}

func (t foo) aFoo() {}
func (t *foo) AFoo() {} // method 'AFoo' differs only by capitalization to method 'aFoo'
```

```golang
func aGlobal() {}
func AGlobal() {} // function 'AGlobal' differs only by capitalization to function 'aGlobal'
```

```golang
type tFoo struct {
	asd string
	aSd int // field 'aSd' differs only by capitalization to other field in the struct type tFoo
}
```

### Valid

```golang
type foo struct{}

func (t foo) Foo() {}
func (t foo) Bar() {}
```

```golang
type tFoo struct {
	name    string
	address string
}
```
