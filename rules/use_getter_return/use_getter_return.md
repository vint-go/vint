---
title: useGetterReturn
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useGetterReturn`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useGetterReturn:
    # no configuration options
```

## Details

Typically, functions with names prefixed with `Get` are supposed to return a value. This rule warns on getter-style functions that do not yield any result.

The rule identifies getters by the following criteria:
- The function name starts with "Get" (case-insensitive prefix match) followed by an uppercase letter.
- Functions named exactly "get" (without a suffix) are not considered getters.
- Functions whose first two parameters match the HTTP handler signature (`http.ResponseWriter`, `*http.Request`) are excluded, since the "Get" prefix in such cases typically refers to HTTP GET.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
// Function named with Get prefix but returns nothing
func GetBar(a, b int) {
}

// Method named with Get prefix but returns nothing
func (t *MyType) GetSaz(a string, b int) {
}
```

### Valid

```golang
// Getter that returns a value
func GetTaz(a string, b int) string {
    return ""
}

// HTTP handler - Get refers to HTTP GET, not a getter
func (b *MyType) GetInfo(w http.ResponseWriter, r *http.Request) {
}

// Not a getter - lowercase letter after "get"
func getfoo() {
}
```
