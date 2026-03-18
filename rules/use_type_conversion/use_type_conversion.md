---
title: useTypeConversion
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useTypeConversion`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useTypeConversion:
    # rule options here
```

## Details

Use a type conversion instead of manually copying struct fields.

When two struct types have identical field names and types, assigning field by field can be replaced with a type conversion, which is clearer and less error-prone.

Source: https://staticcheck.dev/docs/checks/#S1016

## Examples

### Invalid

```golang
type Point2D struct {
    X, Y int
}

type Coordinate struct {
    X, Y int
}

func convert(p Point2D) Coordinate {
    var c Coordinate
    c.X = p.X
    c.Y = p.Y
    return c
}
```

### Valid

```golang
type Point2D struct {
    X, Y int
}

type Coordinate struct {
    X, Y int
}

func convert(p Point2D) Coordinate {
    return Coordinate(p)
}
```
