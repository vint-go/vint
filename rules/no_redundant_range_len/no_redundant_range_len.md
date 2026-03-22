---
title: noRedundantRangeLen
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noRedundantRangeLen`
- This rule is **recommended**, meaning it is enabled by default.
- This rule has a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noRedundantRangeLen:
    # no additional options
```

## Details

Detects `for i := range len(slice)` patterns that can be simplified to `for i := range slice`.

When ranging over `len()` of a slice or array, the `len()` call is unnecessary because Go can range directly over the slice or array. Removing the `len()` call makes the code more concise and idiomatic.

This check only applies to confirmed slice and array types. If the loop variable is unused, the rule suggests `for range slice` instead.

Source: https://github.com/ckaznocha/intrange

## Examples

### Invalid

```golang
// Unnecessary len() in range
s := []int{1, 2, 3}
for i := range len(s) {
    fmt.Println(i, s[i])
}
```

```golang
// Unused loop variable with len()
s := []string{"a", "b", "c"}
for _ = range len(s) {
    fmt.Println("item")
}
```

### Valid

```golang
// Ranging directly over the slice
s := []int{1, 2, 3}
for i := range s {
    fmt.Println(i, s[i])
}
```

```golang
// Ranging directly without variable
s := []string{"a", "b", "c"}
for range s {
    fmt.Println("item")
}
```
