---
title: useImportAliasNaming
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useImportAliasNaming`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useImportAliasNaming:
    arguments:
      - allowRegex: "^[a-z][a-z0-9]{0,}$"
        denyRegex: "^v\\d+$"
```

## Details

Aligns with Go's naming conventions, as outlined in the official
[blog post](https://go.dev/blog/package-names). It enforces clear and lowercase import alias names, echoing
the principles of good package naming. Users can follow these guidelines by default or define a custom regex rule.
Importantly, aliases with underscores ("_") are always allowed.

The rule accepts either a single string argument (used as the allow regex) or a map with `allowRegex` and/or `denyRegex` keys. When no arguments are provided, the default allow regex `^[a-z][a-z0-9]{0,}$` is used.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
import (
    bar_foo "strings" // import name (bar_foo) must match the regular expression: ^[a-z][a-z0-9]{0,}$
    fooBAR "strings"  // import name (fooBAR) must match the regular expression: ^[a-z][a-z0-9]{0,}$
)
```

### Valid

```golang
import (
    magical "magic/hat"
    _ "strings"       // _ aliases are ignored
    . "dotimport"     // . aliases are ignored
    v1 "strings"      // matches default regex
)
```
