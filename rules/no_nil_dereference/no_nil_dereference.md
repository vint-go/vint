---
title: noNilDereference
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noNilDereference`
- This rule is **not recommended**, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noNilDereference:
    # rule options here
```

## Details

Reports nil pointer dereferences and degenerate nil comparisons. This analyzer uses SSA (Static Single Assignment) form to track the flow of nil values through a program and reports:

- Nil pointer dereferences (accessing a field or method on a nil pointer)
- Degenerate nil comparisons (comparisons against nil where the result is always the same)
- Unreachable code caused by nil checks that are always true or false

This analyzer provides deeper analysis than basic nil checks by performing data flow analysis across the entire function.

Source: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/nilness

## Examples

### Invalid

```golang
func example(p *int) int {
    if p == nil {
        // Bad: dereferencing p after confirming it is nil
        return *p
    }
    return 0
}
```

```golang
func example() {
    var p *int
    // Bad: p is always nil, this comparison is degenerate
    if p != nil {
        fmt.Println(*p) // unreachable
    }
}
```

### Valid

```golang
func example(p *int) int {
    if p == nil {
        return 0 // Good: return without dereferencing nil pointer
    }
    return *p // Good: p is guaranteed non-nil here
}
```
