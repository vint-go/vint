---
title: noSingleCaseSelect
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noSingleCaseSelect`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noSingleCaseSelect:
    # no additional options
```

## Details

Use plain channel send or receive instead of single-case select.

A `select` statement with a single case and no default is equivalent to a plain channel operation. The select adds unnecessary complexity.

Source: https://staticcheck.dev/docs/checks/#S1000

## Examples

### Invalid

```golang
select {
case v := <-ch:
    _ = v
}
```

```golang
select {
case ch <- 1:
}
```

### Valid

```golang
v := <-ch
_ = v
```

```golang
ch <- 1
```
