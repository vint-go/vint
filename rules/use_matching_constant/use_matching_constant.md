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

### Options

| Option            | Type | Default | Description                                                              |
|-------------------|------|---------|--------------------------------------------------------------------------|
| `min-occurrences` | int  | 3       | Minimum number of times a string literal must appear before it is reported. |

```yaml title="vint.yaml"
settings:
  lint/suspicious/useMatchingConstant:
    enabled: true
    arguments:
      - min-occurrences: 2
```

## Details

Detects string literals in the code that match the value of an existing named constant but are not referencing it.

When a constant has been defined for a particular string value, all usages of that value should reference the constant rather than using the raw string literal. Using the raw string literal instead of the constant defeats the purpose of having the constant in the first place: it creates a hidden dependency on the value that will not be updated if the constant's value changes.

The rule requires a string literal to appear at least `min-occurrences` times (default 3) before reporting, matching the behavior of goconst's `min-occurrences` setting.

String literals that appear as map keys in composite literals are not counted, since map keys are structurally different from general string usage and often coincidentally share the same text as a constant without being semantically related.

This check is enabled via the `-match-constant` flag in goconst or `match-constant: true` in golangci-lint configuration.

Source: https://github.com/jgautheron/goconst

## Examples

### Invalid

```golang
const StatusActive = "active"

// "active" appears 3+ times, matching the constant.
func IsActive(status string) bool {
    return status == "active"
}

func IsActive2(status string) bool {
    return status == "active"
}

func IsActive3(status string) bool {
    return status == "active"
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
const StatusActive = "active"

// Only one occurrence of the literal — below the default threshold.
func IsActive(status string) bool {
    return status == "active"
}
```

```golang
const RoleUser = "user"

// Map keys are not counted as occurrences.
func Defaults() map[string]string {
    return map[string]string{
        "user": "default",
    }
}
```

```golang
// No constant exists for "pending", so a raw literal is fine.
func IsPending(status string) bool {
    return status == "pending"
}
```
