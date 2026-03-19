---
title: noInefficientMapLookup
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/performance/noInefficientMapLookup`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/performance/noInefficientMapLookup:
    # no options
```

## Details

This rule identifies code that iteratively searches for a key in a map.

This inefficiency is usually introduced when refactoring code from using a slice to a map.
For example if during refactoring the `elements` slice is transformed into a map, and then
a loop over `elements` is changed in an obvious but inefficient way: iterating over map keys
and comparing each key to a target value, instead of performing a direct map lookup.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
aMap := map[string]bool{}
target := "foo"

// Inefficient: iterating over map keys to find a match
for k := range aMap {
  if k == target {
    // do something
  }
}
```

```golang
aMap := map[string]bool{}
target := "foo"

// Inefficient: iterating and skipping non-matching keys
for k := range aMap {
  if k != target {
    continue
  }
  // do something with k
}
```

### Valid

```golang
aMap := map[string]bool{}
target := "foo"

// Efficient: direct map lookup
if _, ok := aMap[target]; ok {
  // do something
}
```
