---
title: useVarConstDocPrefix
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useVarConstDocPrefix`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useVarConstDocPrefix:
    # rule options here
```

## Details

The documentation of an exported variable or constant should start with the name of the variable or constant.

By Go convention, documentation comments for exported variables and constants should begin with their name.

Source: https://staticcheck.dev/docs/checks/#ST1022

## Examples

### Invalid

```golang
package main

// The default timeout for HTTP requests.
var DefaultTimeout = 30
```

### Valid

```golang
package main

// DefaultTimeout is the default timeout for HTTP requests.
var DefaultTimeout = 30
```
