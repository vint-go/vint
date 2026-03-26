---
title: noSliceBoundsOutOfRange
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noSliceBoundsOutOfRange`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noSliceBoundsOutOfRange:
    # rule options here
```

## Details

Detects possible slice bounds out of range errors.

This rule uses SSA (Static Single Assignment) analysis to identify code patterns where slice access could panic due to out-of-bounds indexing. Accessing a slice element or creating a sub-slice with an index that exceeds the slice's length causes a runtime panic, which can lead to denial of service in server applications.

The rule checks for patterns where a slice is accessed with an index that is not properly validated against the slice's length. This includes direct indexing, sub-slicing, and cases where length checks are missing or insufficient.

The rule recognizes the following safe patterns and does not report them:

- **`len()` bounds checks**: An `if` statement whose condition involves `len(slice)` that encloses the access.
- **`for i := range slice` loops**: The range index variable `i` is inherently bounded by the slice length, so `slice[i]` inside such a loop is safe.
- **`sort.Slice` / `sort.SliceStable` callbacks**: The callback parameters `i` and `j` are guaranteed valid indices into the sorted slice.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
func getElement(s []int, idx int) int {
    // No bounds check before access
    return s[idx]
}
```

```golang
func getFirst(s []string) string {
    // May panic if slice is empty
    return s[0]
}
```

```golang
func subSlice(s []int) []int {
    // May panic if slice has fewer than 5 elements
    return s[2:5]
}
```

### Valid

```golang
func getElement(s []int, idx int) (int, error) {
    // Bounds check before access
    if idx < 0 || idx >= len(s) {
        return 0, fmt.Errorf("index %d out of range [0, %d)", idx, len(s))
    }
    return s[idx], nil
}
```

```golang
func getFirst(s []string) (string, bool) {
    // Check if slice is non-empty
    if len(s) == 0 {
        return "", false
    }
    return s[0], true
}
```

```golang
func subSlice(s []int) []int {
    // Validate bounds before sub-slicing
    if len(s) < 5 {
        return nil
    }
    return s[2:5]
}
```

```golang
func sumElements(s []int) int {
    // Range index is inherently bounded
    sum := 0
    for i := range s {
        sum += s[i]
    }
    return sum
}
```

```golang
func sortEndpoints(endpoints []string) {
    // sort.Slice callback parameters are safe indices
    sort.Slice(endpoints, func(i, j int) bool {
        return endpoints[i] < endpoints[j]
    })
}
```
