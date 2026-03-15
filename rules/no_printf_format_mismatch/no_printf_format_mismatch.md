---
title: noPrintfFormatMismatch
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noPrintfFormatMismatch`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noPrintfFormatMismatch:
    # rule options here
```

## Details

Checks consistency of Printf format strings and arguments. This analyzer verifies that calls to `fmt.Printf`, `fmt.Sprintf`, `fmt.Fprintf`, `log.Printf`, and similar formatting functions have format strings that match the provided arguments. It detects:

- Wrong number of arguments for the format string
- Argument types that do not match the format verb (e.g., `%d` with a string)
- Invalid or unknown format verbs
- Printf-style functions called with arguments but no format string

Source: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/printf

## Examples

### Invalid

```golang
import "fmt"

func example() {
    name := "world"
    // Bad: %d expects an integer, got a string
    fmt.Printf("Hello %d", name)
}
```

```golang
import "fmt"

func example() {
    // Bad: too few arguments for format string
    fmt.Printf("Name: %s, Age: %d", "Alice")
}
```

```golang
import "fmt"

func example() {
    // Bad: Println does not support format verbs
    fmt.Println("Hello %s", "world")
}
```

### Valid

```golang
import "fmt"

func example() {
    name := "world"
    fmt.Printf("Hello %s", name)
}
```

```golang
import "fmt"

func example() {
    fmt.Printf("Name: %s, Age: %d", "Alice", 30)
}
```

```golang
import "fmt"

func example() {
    fmt.Println("Hello", "world")
}
```
