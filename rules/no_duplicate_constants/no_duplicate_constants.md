---
title: noDuplicateConstants
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noDuplicateConstants`
- This rule is **not recommended**, meaning it is not enabled by default. It must be explicitly enabled in your configuration.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noDuplicateConstants:
    enabled: true
```

## Details

Detects named constants that have the same value as another constant, indicating potential duplication or a missed opportunity for consolidation.

When two or more constants share the same underlying value, it often suggests that they represent the same concept and should be unified into a single constant. Duplicate constants can lead to confusion about which one should be used, make refactoring harder, and increase the risk of introducing inconsistencies if only one of the duplicates is updated.

This check is enabled via the `-find-duplicates` flag in goconst or `find-duplicates: true` in golangci-lint configuration.

Note that there are legitimate reasons for having separate constants with the same value -- for example, when the constants represent conceptually different things that happen to share the same value today but might diverge in the future. Use judgment when addressing findings from this rule.

Source: https://github.com/jgautheron/goconst

## Examples

### Invalid

```golang
const (
    StatusActive  = "active"
    UserActive    = "active"  // Duplicate value: same as StatusActive
)
```

```golang
const (
    DefaultTimeout  = 30
    RequestTimeout  = 30  // Duplicate value: same as DefaultTimeout
    MaxRetryWait    = 30  // Duplicate value: same as DefaultTimeout
)
```

```golang
const (
    ErrNotFound     = "not found"
    ErrMissing      = "not found"  // Duplicate value: same as ErrNotFound
)
```

### Valid

```golang
const (
    StatusActive   = "active"
    StatusInactive = "inactive"
    StatusPending  = "pending"
)
```

```golang
// If two constants must share the same value, reference one from the other.
const (
    DefaultTimeout = 30
    RequestTimeout = DefaultTimeout
)
```

```golang
// Conceptually different constants with distinct values.
const (
    MaxRetries     = 3
    DefaultTimeout = 30
    BufferSize     = 1024
)
```
