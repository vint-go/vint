---
title: useIdiomaticNaming
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useIdiomaticNaming`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useIdiomaticNaming:
    # rule options here
```

## Details

Poorly chosen identifier.

Go has specific naming conventions. Names should use MixedCaps or mixedCaps (camelCase), not underscores. Acronyms should be all caps (e.g., `HTTP`, `ID`, `URL`). This check flags identifiers that don't follow these conventions.

Source: https://staticcheck.dev/docs/checks/#ST1003

## Examples

### Invalid

```golang
package main

// snake_case names
var my_variable int

// Incorrect casing for acronyms
type HttpClient struct{}
type XmlParser struct{}
func GetUrlPath() string { return "" }
```

### Valid

```golang
package main

// camelCase names
var myVariable int

// Correct casing for acronyms
type HTTPClient struct{}
type XMLParser struct{}
func GetURLPath() string { return "" }
```
