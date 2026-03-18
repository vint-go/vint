---
title: useEarlyContinue
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/complexity/useEarlyContinue`
- This rule is not recommended (experimental, opinionated).
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/complexity/useEarlyContinue:
    bodyWidth: 5  # minimum body statements to trigger (default: 5)
```

## Details

Finds where nesting level could be reduced in loop bodies. This checker identifies loop constructs (for and range statements) containing a single if statement with a body of at least a configurable number of statements (default: 5) and no else clause. It suggests inverting the condition and using `continue` to skip unnecessary nesting, reducing code indentation depth.

This rule is distinct from `noHighCyclomaticComplexity` (which measures the number of independent execution paths), `noLongFunctions` (which measures line count), and `noExcessiveStatements` (which measures statement count). This rule specifically targets nesting depth inside loops and suggests a structural refactoring pattern rather than measuring a numeric threshold.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
for _, item := range items {
    if item.IsValid() {
        process(item)
        transform(item)
        save(item)
        log(item)
        notify(item)
    }
}
```

### Valid

```golang
for _, item := range items {
    if !item.IsValid() {
        continue
    }
    process(item)
    transform(item)
    save(item)
    log(item)
    notify(item)
}
```
