---
title: useTypeDocPrefix
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useTypeDocPrefix`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useTypeDocPrefix:
    # rule options here
```

## Details

The documentation of an exported type should start with the type's name.

By Go convention, documentation comments for exported types should begin with the type name, e.g., `// MyType represents ...`

Source: https://staticcheck.dev/docs/checks/#ST1021

## Examples

### Invalid

```golang
package main

// Represents a server configuration.
type Config struct {
    Host string
    Port int
}
```

### Valid

```golang
package main

// Config represents a server configuration.
type Config struct {
    Host string
    Port int
}
```
