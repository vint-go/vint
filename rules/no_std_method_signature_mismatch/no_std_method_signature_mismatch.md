---
title: noStdMethodSignatureMismatch
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noStdMethodSignatureMismatch`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noStdMethodSignatureMismatch:
    # rule options here
```

## Details

Checks for misspellings in the signatures of methods similar to well-known interfaces. If a type has a method whose name matches a well-known interface method (such as `Read`, `Write`, `String`, `MarshalJSON`, `UnmarshalJSON`, etc.) but with a different signature, it is likely a mistake. The method will not satisfy the interface, and the programmer probably intended to implement the interface correctly.

Well-known methods checked include those from `io.Reader`, `io.Writer`, `fmt.Stringer`, `encoding/json.Marshaler`, `encoding/json.Unmarshaler`, and many more.

Source: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/stdmethods

## Examples

### Invalid

```golang
type MyType struct{}

// Bad: wrong signature for fmt.Stringer interface
// Should return string, not error
func (m MyType) String() error {
    return nil
}
```

```golang
type MyType struct{}

// Bad: wrong signature for encoding/json.Marshaler
// Should be MarshalJSON() ([]byte, error)
func (m MyType) MarshalJSON() error {
    return nil
}
```

### Valid

```golang
type MyType struct{}

// Good: correct signature for fmt.Stringer
func (m MyType) String() string {
    return "MyType"
}
```

```golang
type MyType struct{}

// Good: correct signature for encoding/json.Marshaler
func (m MyType) MarshalJSON() ([]byte, error) {
    return []byte(`"mytype"`), nil
}
```
