---
title: useFilenameFormat
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useFilenameFormat`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useFilenameFormat:
    arguments:
      - "^[_a-z][_a-z0-9]*\\.go$"
```

## Details

Enforces conventions on source file names. By default, the rule enforces filenames of the form `^[_A-Za-z0-9][_A-Za-z0-9-]*\.go$`. Optionally, the rule can be configured with a custom regular expression to enforce other naming conventions. The rule also detects and reports non-ASCII characters found in filenames.

The rule accepts a single string argument: a regular expression that source filenames must match.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
// File: my file.go
// Filename "my file.go" does not match the default pattern.
package main
```

```golang
// File: filenamе_with_cyrillic.go (contains Cyrillic 'е')
// Non-ASCII characters in filenames are flagged.
package main
```

### Valid

```golang
// File: my_module.go
package main
```

```golang
// File: utils-helper.go
package main
```
