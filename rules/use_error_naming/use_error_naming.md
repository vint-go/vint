---
title: useErrorNaming
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useErrorNaming`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useErrorNaming:
    # no configuration options
```

## Details

By convention, for the sake of readability, variables of type `error` must be named with the prefix `err` (for unexported variables) or `Err` (for exported variables). This rule checks package-level variable declarations that are initialized with `errors.New` or `fmt.Errorf` and ensures they follow this naming convention.

The rule does not check local variables inside functions, only package-level declarations. It also allows the blank identifier `_` to be used without triggering a warning.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
package foo

import "errors"

// Unexported error variable without "err" prefix
var unexp = errors.New("some unexported error")

// Exported error variable without "Err" prefix
var Exp = errors.New("some exported error")
```

### Valid

```golang
package foo

import "errors"

// Unexported error variable with correct "err" prefix
var errUnexp = errors.New("some unexported error")

// Exported error variable with correct "Err" prefix
var ErrExp = errors.New("some exported error")
```
