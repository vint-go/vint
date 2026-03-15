---
title: useMatchingConstant
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/useMatchingConstant`
- This rule is **not recommended**, meaning it is not enabled by default. It must be explicitly enabled in your configuration.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/useMatchingConstant:
    enabled: true
```

## Details

Detects string literals in the code that match the value of an existing named constant but are not referencing it.

When a constant has been defined for a particular string value, all usages of that value should reference the constant rather than using the raw string literal. Using the raw string literal instead of the constant defeats the purpose of having the constant in the first place: it creates a hidden dependency on the value that will not be updated if the constant's value changes.

This check is enabled via the `-match-constant` flag in goconst or `match-constant: true` in golangci-lint configuration.

Source: https://github.com/jgautheron/goconst

## Examples

### Invalid

```golang
const StatusActive = "active"

func IsActive(status string) bool {
    // "active" matches the existing constant StatusActive,
    // but the constant is not being used here.
    return status == "active"
}
```

```golang
const (
    RoleAdmin = "admin"
    RoleUser  = "user"
)

func CheckRole(role string) bool {
    // "admin" matches the constant RoleAdmin but is used as a raw literal.
    if role == "admin" {
        return true
    }
    // "user" matches the constant RoleUser but is used as a raw literal.
    if role == "user" {
        return true
    }
    return false
}
```

### Valid

```golang
const StatusActive = "active"

func IsActive(status string) bool {
    // Correctly references the named constant.
    return status == StatusActive
}
```

```golang
const (
    RoleAdmin = "admin"
    RoleUser  = "user"
)

func CheckRole(role string) bool {
    if role == RoleAdmin {
        return true
    }
    if role == RoleUser {
        return true
    }
    return false
}
```

```golang
// No constant exists for "pending", so a raw literal is fine.
func IsPending(status string) bool {
    return status == "pending"
}
```
