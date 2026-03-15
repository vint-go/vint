---
title: noUnusedWrite
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noUnusedWrite`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noUnusedWrite:
    # rule options here
```

## Details

Checks for unused writes to elements of a struct or array. This analyzer detects writes to fields of a struct or elements of an array where the struct or array is never used after the write. This often happens when a function copies a struct by value, modifies a field of the copy, and then never uses the copy again. The modification has no effect on the original value.

Source: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/unusedwrite

## Examples

### Invalid

```golang
type Point struct {
    X, Y int
}

func example() Point {
    p := Point{X: 1, Y: 2}
    q := p
    q.X = 10 // Bad: q is a copy and is never used after this write
    return p  // returns original, not modified copy
}
```

### Valid

```golang
type Point struct {
    X, Y int
}

func example() Point {
    p := Point{X: 1, Y: 2}
    p.X = 10 // Good: writing to p which is then returned
    return p
}
```

```golang
type Point struct {
    X, Y int
}

func example() Point {
    p := Point{X: 1, Y: 2}
    q := p
    q.X = 10
    return q // Good: modified copy is used
}
```
