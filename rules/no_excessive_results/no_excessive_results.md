---
title: noExcessiveResults
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/complexity/noExcessiveResults`
- This rule is not recommended (experimental, opinionated).
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/complexity/noExcessiveResults:
    maxResults: 5  # maximum number of return values (default: 5)
```

## Details

Detects functions with too many results. When a function returns more values than the configured threshold, this checker suggests simplifying the function, potentially by grouping results into a struct. Functions with many return values are harder to use correctly and indicate the function may be doing too much.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
func getDetails() (string, int, bool, error, []byte, map[string]string) {
    // too many return values
}
```

### Valid

```golang
type Details struct {
    Name    string
    Count   int
    Active  bool
    Data    []byte
    Meta    map[string]string
}

func getDetails() (Details, error) {
    // grouped into a struct
}
```
