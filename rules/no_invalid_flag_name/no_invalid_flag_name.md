---
title: noInvalidFlagName
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noInvalidFlagName`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noInvalidFlagName:
    # no additional options
```

## Details

Detects suspicious flag names. This checker validates flag names passed to Go's `flag` package functions and warns when flag names:

1. Are empty strings
2. Start with a hyphen (`-`)
3. Contain equals signs (`=`)
4. Contain whitespace characters

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
flag.Bool("-verbose", false, "enable verbose mode")
```

```golang
flag.String("output=file", "", "output file path")
```

```golang
flag.Int("", 0, "some value")
```

### Valid

```golang
flag.Bool("verbose", false, "enable verbose mode")
```

```golang
flag.String("output", "", "output file path")
```
