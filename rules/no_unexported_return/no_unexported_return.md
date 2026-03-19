---
title: noUnexportedReturn
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noUnexportedReturn`
- This rule is recommended, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noUnexportedReturn:
    # no configuration options
```

## Details

This rule warns when an exported function or method returns a value of an unexported type. Returning an unexported type from an exported function can be annoying to use because callers in other packages cannot refer to the return type by name, making it difficult to store, pass, or document such values.

The rule skips files that are not importable (main packages and test files), as their symbols cannot be used in other packages. It also skips exported methods on unexported receiver types, such as private implementations of well-known interfaces.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
package mylib

type myStruct struct{}

// Exported function returns unexported type
func NewMyStruct() myStruct {
    return myStruct{}
}
```

```golang
package mylib

type unexportedResult struct{}

type MyType struct{}

// Exported method returns unexported type
func (m MyType) GetResult() unexportedResult {
    return unexportedResult{}
}
```

### Valid

```golang
package mylib

type MyStruct struct{}

// Exported function returns exported type
func NewMyStruct() MyStruct {
    return MyStruct{}
}
```

```golang
package mylib

type myStruct struct{}

// Unexported function can return unexported type
func newMyStruct() myStruct {
    return myStruct{}
}
```

```golang
package main

type foo struct{}

// Main package symbols are not importable
func NewFoo() foo {
    return foo{}
}
```
