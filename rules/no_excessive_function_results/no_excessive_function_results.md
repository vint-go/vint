---
title: noExcessiveFunctionResults
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/complexity/noExcessiveFunctionResults`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/complexity/noExcessiveFunctionResults:
    arguments:
      - 3
```

## Details

Functions returning too many results can be hard to understand/use. This rule limits the maximum number of return results a function can have.

The argument is an integer specifying the maximum allowed number of return values. The default is `3`.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
// With max set to 3:
func foo() (a, b, c, d int) { // 4 results exceed the limit of 3
    var a, b, c, d int
}

func qux() (string, string, int, string, int) { // 5 results exceed the limit of 3
}
```

### Valid

```golang
// With max set to 3:
func bar(a, b int) { // no results
}

func baz() (string, error) { // 2 results, within the limit
    return "", nil
}

func qux() (int, string, error) { // 3 results, within the limit
    return 0, "", nil
}
```
