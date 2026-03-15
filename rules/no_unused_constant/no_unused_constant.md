---
title: noUnusedConstant
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noUnusedConstant`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noUnusedConstant:
    # Controls whether constants in generated files are considered used.
    # generated-is-used: true
```

## Details

Detects constants that are declared but never referenced anywhere in the codebase.

The `unused` linter uses a graph-based analysis to determine whether constants are reachable from entry points. A constant is considered "used" if it is referenced in code that is itself reachable.

Key rules governing constant usage:
- Exported constants are always considered used (since external packages may reference them).
- If one constant in a `const` block (iota group) is used, all constants in that block are marked as used (rule 10.1). This prevents false positives when constants are part of an enumeration where not every value is explicitly referenced.
- Constants in generated files are considered used when the `generated-is-used` option is enabled (default).
- The blank identifier constant (`_`) is always considered used.

Source: https://github.com/dominikh/go-tools/tree/master/unused

## Examples

### Invalid

```golang
package mypackage

const maxRetries = 5 // const maxRetries is unused

func Process() error {
    return nil
}
```

```golang
package mypackage

const (
    internalTimeout = 30  // const internalTimeout is unused
    internalBuffer  = 256 // const internalBuffer is unused
)

func Run() {}
```

### Valid

```golang
package mypackage

const maxRetries = 5

func Process() error {
    for i := 0; i < maxRetries; i++ {
        // retry logic
    }
    return nil
}
```

```golang
package mypackage

// Exported constants are always considered used
const Version = "2.0.0"
```

```golang
package mypackage

// If one constant in a block is used, all are considered used
const (
    StatusPending  = iota
    StatusActive
    StatusInactive
    StatusArchived
)

func DefaultStatus() int {
    return StatusPending // only StatusPending is directly used,
                         // but all constants in the block are considered used
}
```

```golang
package mypackage

const _ = "blank identifier is always considered used"
```
