---
title: noMalformedTestFunction
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noMalformedTestFunction`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noMalformedTestFunction:
    # rule options here
```

## Details

Checks for common mistaken usages of tests and examples. This analyzer validates that test, benchmark, fuzz, and example functions in `_test.go` files follow the correct naming conventions and signatures required by the `go test` command. It detects:

- Test functions with wrong signatures (e.g., `TestFoo(t *testing.B)` instead of `TestFoo(t *testing.T)`)
- Example functions that reference non-existent identifiers
- Benchmark functions with wrong signatures
- Fuzz functions with wrong signatures

Source: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/tests

## Examples

### Invalid

```golang
import "testing"

// Bad: test function has wrong parameter type
func TestSomething(b *testing.B) {
    // ...
}
```

```golang
import "testing"

// Bad: benchmark function has wrong parameter type
func BenchmarkSomething(t *testing.T) {
    // ...
}
```

```golang
// Bad: example function references non-existent function
func ExampleNonExistentFunction() {
    // Output: hello
}
```

### Valid

```golang
import "testing"

// Good: correct test function signature
func TestSomething(t *testing.T) {
    // ...
}
```

```golang
import "testing"

// Good: correct benchmark function signature
func BenchmarkSomething(b *testing.B) {
    for i := 0; i < b.N; i++ {
        // ...
    }
}
```
