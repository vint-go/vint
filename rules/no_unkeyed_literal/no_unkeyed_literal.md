---
title: noUnkeyedLiteral
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noUnkeyedLiteral`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noUnkeyedLiteral:
    # rule options here
```

## Details

Checks for unkeyed composite literals. This analyzer flags composite literals of structs that do not use field names (keys). Unkeyed composite literals are fragile because adding a new field to the struct will cause compilation errors in all places that use unkeyed literals.

Using keyed literals makes code more readable and robust against changes to the struct definition. Matching the behavior of `go vet`'s composite analyzer, this rule only flags struct types from **external packages**. Same-package types and anonymous structs are exempt, since changes to those are within the developer's control.

Source: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/composite

## Examples

### Invalid

```golang
import "image"

func example() {
    // Unkeyed composite literal - fragile if image.Point changes
    p := image.Point{1, 2}
}
```

### Valid

```golang
import "image"

func example() {
    // Keyed composite literal - robust against struct changes
    p := image.Point{X: 1, Y: 2}
}
```

```golang
func example() {
    // Unkeyed literals for basic types are acceptable
    s := []int{1, 2, 3}
}
```

```golang
// Same-package types are exempt
type MyPoint struct {
    X, Y int
}

func example() {
    p := MyPoint{1, 2} // OK - same package
}
```

```golang
func example() {
    // Anonymous struct literals are exempt
    a := struct{ X, Y int }{1, 2} // OK
}
```
