---
title: noIncorrectTimeFormat
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noIncorrectTimeFormat`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noIncorrectTimeFormat:
    # rule options here
```

## Details

Checks for the use of incorrect time format strings in calls to `time.Format` and `time.Parse`. Go uses a reference time (`Mon Jan 2 15:04:05 MST 2006`) for defining time layouts, which is different from the common `YYYY-MM-DD` style used in other languages.

A common mistake is using formats like `2006-01-02 15:04:05` with incorrect reference values, such as using `01` for the day instead of `02`, or `04` for the hour instead of `15`.

This analyzer also checks for the use of `2006-02-01` instead of `2006-01-02` (swapped month and day).

Source: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/timeformat

## Examples

### Invalid

```golang
import "time"

func example() {
    t := time.Now()
    // Bad: month and day are swapped (should be 01-02, not 02-01)
    s := t.Format("2006-02-01")
    _ = s
}
```

### Valid

```golang
import "time"

func example() {
    t := time.Now()
    // Good: correct reference time layout
    s := t.Format("2006-01-02 15:04:05")
    _ = s
}
```

```golang
import "time"

func example() {
    // Good: using predefined time constants
    s := time.Now().Format(time.RFC3339)
    _ = s
}
```
