---
title: useErrorf
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useErrorf`
- This rule is recommended, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useErrorf:
    # no configuration options
```

## Details

It is possible to get a simpler program by replacing `errors.New(fmt.Sprintf())` with `fmt.Errorf()`. This rule spots that kind of simplification opportunities.

The rule also detects `t.Error(fmt.Sprintf(...))` in tests and suggests replacing it with `t.Errorf(...)`.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
// Using errors.New with fmt.Sprintf is unnecessarily verbose
err := errors.New(fmt.Sprintf("error: %s", val))
```

```golang
// Using t.Error with fmt.Sprintf in tests
t.Error(fmt.Sprintf("unexpected value: %d", got))
```

### Valid

```golang
// Use fmt.Errorf directly
err := fmt.Errorf("error: %s", val)
```

```golang
// Use t.Errorf directly in tests
t.Errorf("unexpected value: %d", got)
```
