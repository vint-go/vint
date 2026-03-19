---
title: useErrorsNew
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useErrorsNew`
- This rule is **not recommended**, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useErrorsNew:
    # no configuration options
```

## Details

This rule identifies calls to `fmt.Errorf` that can be safely replaced by the more efficient `errors.New`. When `fmt.Errorf` is called with a single string argument and no formatting verbs, there is no reason to use the heavier formatting function. Using `errors.New` instead is more direct, avoids unnecessary overhead, and makes the intent clearer.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
// fmt.Errorf with a single string argument should use errors.New instead
err := fmt.Errorf("connection refused")
```

```golang
// fmt.Errorf with a concatenated string (still a single argument) should use errors.New
err := fmt.Errorf(msg + "something")
```

### Valid

```golang
// fmt.Errorf with formatting verbs is legitimate
err := fmt.Errorf("unable to load repo: %w", err)
```

```golang
// fmt.Errorf with multiple arguments is legitimate
err := fmt.Errorf("failed to get commit for %s: %w", branch, err)
```

```golang
// errors.New is the correct choice for plain error strings
err := errors.New("connection refused")
```
