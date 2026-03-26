---
title: noHugeParam
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/performance/noHugeParam`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/performance/noHugeParam:
    sizeThreshold: 80  # size in bytes at which a parameter is considered huge (default: 80)
```

## Details

Detects params that incur excessive amount of copying. This checker identifies function parameters that are greater than or equal to a specified byte threshold, as passing large structures by value causes unnecessary copying overhead. The checker suggests using a pointer instead. The `String() string` method is automatically excluded to avoid flagging Stringer interface implementations. Unnamed receivers and parameters (e.g., `func (MyType) Method()`) are silently skipped, matching gocritic behavior.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
func processData(data [1024]int) {
    // data is copied on each call
}
```

### Valid

```golang
func processData(data *[1024]int) {
    // data is passed by pointer, no copying
}
```

```golang
func (s MyStruct) String() string {
    // String() method is excluded from this check
    return fmt.Sprintf("%v", s)
}
```

```golang
func (Application) TableName() string {
    // Unnamed receivers are skipped
    return "applications"
}
```
