---
title: useFuncDocPrefix
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useFuncDocPrefix`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useFuncDocPrefix:
    # rule options here
```

## Details

The documentation of an exported function should start with the function's name.

By Go convention, documentation comments for exported functions should begin with the function name, e.g., `// MyFunc does something.`

Source: https://staticcheck.dev/docs/checks/#ST1020

## Examples

### Invalid

```golang
package main

// This function processes data.
func ProcessData() {}
```

### Valid

```golang
package main

// ProcessData processes the given data and returns results.
func ProcessData() {}
```
