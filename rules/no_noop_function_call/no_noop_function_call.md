---
title: noNoopFunctionCall
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noNoopFunctionCall`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noNoopFunctionCall:
    # no additional options
```

## Details

Detects suspicious function calls. This checker identifies function calls that are likely mistakes, such as calling `strings.Replace` with a count of 0 (which replaces nothing), or other standard library functions with arguments that make the call a no-op.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
strings.Replace(s, "old", "new", 0) // replaces 0 occurrences (no-op)
```

```golang
strings.SplitN(s, ",", 0) // returns nil
```

### Valid

```golang
strings.Replace(s, "old", "new", -1) // replaces all occurrences
```

```golang
strings.ReplaceAll(s, "old", "new")
```

```golang
strings.SplitN(s, ",", -1) // splits all
```
