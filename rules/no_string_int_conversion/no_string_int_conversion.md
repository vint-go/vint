---
title: noStringIntConversion
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noStringIntConversion`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noStringIntConversion:
    # rule options here
```

## Details

Flags type conversions from integers to strings. In Go, `string(n)` where `n` is an integer does not produce the decimal string representation of the number. Instead, it produces the string containing the Unicode code point corresponding to that integer value. For example, `string(65)` produces `"A"`, not `"65"`.

To convert an integer to its decimal string representation, use `strconv.Itoa` or `fmt.Sprint`.

Source: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/stringintconv

## Examples

### Invalid

```golang
func example() {
    n := 65
    // Bad: produces "A" (Unicode code point 65), not "65"
    s := string(n)
    _ = s
}
```

```golang
func example(code int) string {
    // Bad: converts integer to Unicode code point, not decimal string
    return string(code)
}
```

### Valid

```golang
import "strconv"

func example() {
    n := 65
    // Good: produces "65"
    s := strconv.Itoa(n)
    _ = s
}
```

```golang
import "fmt"

func example(code int) string {
    // Good: produces decimal string representation
    return fmt.Sprint(code)
}
```

```golang
func example() {
    // Good: explicit rune-to-string conversion is sometimes intentional
    s := string(rune(65)) // "A"
    _ = s
}
```
