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
    extraInitialisms:   # additional initialisms beyond defaults
      - GRPC
      - AMQP
    excludeInitialisms: # default initialisms to skip checking
      - ID
```

### Options

| Option | Type | Description |
|---|---|---|
| `extraInitialisms` | `[]string` | Additional initialisms to enforce beyond the built-in defaults (e.g. `GRPC`, `AMQP`). Names containing these words will be flagged if not fully capitalized. |
| `excludeInitialisms` | `[]string` | Built-in default initialisms to exclude from checking (e.g. `ID`, `URL`). Names containing these words will not be flagged for initialism casing. |

The built-in default initialisms are: `ACL`, `API`, `ASCII`, `CPU`, `CSS`, `DNS`, `EOF`, `GUID`, `HTML`, `HTTP`, `HTTPS`, `ID`, `IDS`, `IP`, `JSON`, `LHS`, `QPS`, `RAM`, `RHS`, `RPC`, `SLA`, `SMTP`, `SQL`, `SSH`, `TCP`, `TLS`, `TTL`, `UDP`, `UI`, `UID`, `UUID`, `URI`, `URL`, `UTF8`, `VM`, `XML`, `XMPP`, `XSRF`, `XSS`.

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
