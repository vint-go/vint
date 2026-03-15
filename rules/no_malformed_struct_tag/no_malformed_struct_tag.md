---
title: noMalformedStructTag
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noMalformedStructTag`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noMalformedStructTag:
    # rule options here
```

## Details

Checks that struct field tags are well formed. This analyzer validates that struct field tags follow the correct syntax as defined by `reflect.StructTag`. It detects:

- Malformed tag syntax (e.g., missing quotes, invalid characters)
- Duplicate tag keys (e.g., two `json` keys on the same field)
- Misuse of tag options (e.g., unknown options for `json` or `xml` tags)
- Incorrectly structured tag key-value pairs

Source: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/structtag

## Examples

### Invalid

```golang
type User struct {
    // Bad: duplicate json tag key
    Name string `json:"name" json:"username"`
}
```

```golang
type User struct {
    // Bad: malformed tag syntax (missing closing quote)
    Name string `json:"name`
}
```

```golang
type User struct {
    // Bad: space in tag value without proper quoting
    Name string `json: "name"`
}
```

### Valid

```golang
type User struct {
    Name  string `json:"name"`
    Email string `json:"email,omitempty"`
    Age   int    `json:"age" xml:"age"`
}
```

```golang
type User struct {
    Name string `json:"name" validate:"required"`
    ID   int    `json:"-"` // excluded from JSON
}
```
