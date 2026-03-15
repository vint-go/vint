---
title: noReflectValueCompare
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noReflectValueCompare`
- This rule is **not recommended**, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noReflectValueCompare:
    # rule options here
```

## Details

Checks for accidentally using `==` or `reflect.DeepEqual` to compare `reflect.Value` values. Comparing `reflect.Value` with `==` compares the reflect metadata, not the underlying values. Similarly, `reflect.DeepEqual` on `reflect.Value` may not produce the expected results.

To compare the underlying values, use the `reflect.Value.Equal` method or `Interface()` to extract the actual values first.

Source: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/reflectvaluecompare

## Examples

### Invalid

```golang
import "reflect"

func example(a, b reflect.Value) bool {
    // Bad: compares reflect metadata, not underlying values
    return a == b
}
```

```golang
import "reflect"

func example(a, b reflect.Value) bool {
    // Bad: DeepEqual on reflect.Value may not compare underlying values
    return reflect.DeepEqual(a, b)
}
```

### Valid

```golang
import "reflect"

func example(a, b reflect.Value) bool {
    // Good: compares underlying values
    return a.Equal(b)
}
```

```golang
import "reflect"

func example(a, b reflect.Value) bool {
    // Good: extract interface values and compare
    return reflect.DeepEqual(a.Interface(), b.Interface())
}
```
