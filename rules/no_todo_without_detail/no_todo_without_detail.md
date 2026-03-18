---
title: noTodoWithoutDetail
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noTodoWithoutDetail`
- This rule is not recommended (experimental, opinionated).
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noTodoWithoutDetail:
    # no additional options
```

## Details

Detects TODO comments without detail/assignee. This checker identifies TODO, FIX, FIXME, and BUG comments that lack additional context. It suggests adding clarifying information such as assignee names or explanatory details to make the comments more actionable for developers reviewing the code.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
// TODO
func processData() {}
```

```golang
// FIXME
func brokenFunction() {}
```

### Valid

```golang
// TODO(jsmith): implement error handling for edge cases
func processData() {}
```

```golang
// FIXME(team): this function panics when input is empty, needs guard clause
func brokenFunction() {}
```
