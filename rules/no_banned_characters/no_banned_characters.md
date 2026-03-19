---
title: noBannedCharacters
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noBannedCharacters`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noBannedCharacters:
    arguments:
      - "Ω"
      - "Σ"
      - "σ"
```

## Details

Checks identifiers (constants, variables, functions, etc.) for the presence of banned characters.

This rule is useful for enforcing naming conventions that prohibit certain Unicode or special characters in identifier names. Characters appearing in comments or string literals are not flagged. Each banned character is specified as a string in the rule arguments.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
// With banned characters configured as ["Ω", "Σ", "σ"]

const Ω = "Omega" // error: banned character found: Ω

func funcΣ() error { // error: banned character found: Σ
    var charσhid string // error: banned character found: σ
    return nil
}
```

### Valid

```golang
// With banned characters configured as ["Ω", "Σ", "σ"]

const Omega = "Omega" // OK: no banned characters in identifier

func funcSigma() error { // OK: no banned characters in identifier
    var charHid string // OK: no banned characters in identifier
    return nil
}

// Banned characters in comments are OK: Ω Σ σ
```
