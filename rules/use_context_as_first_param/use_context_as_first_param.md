---
title: useContextAsFirstParam
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useContextAsFirstParam`
- This rule is recommended, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useContextAsFirstParam:
    arguments:
      - allowTypesBefore: "*testing.T,*github.com/user/repo/testing.Harness"
```

## Details

By [convention](https://go.dev/wiki/CodeReviewComments#contexts), `context.Context` should be the first parameter of a function. This rule spots function declarations that do not follow the convention.

An optional `allowTypesBefore` configuration parameter accepts a comma-separated list of types that are permitted to appear before `context.Context` in a function signature. The option name is case-insensitive and accepts both camelCase (`allowTypesBefore`) and kebab-case (`allow-types-before`).

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
// context.Context is not the first parameter
func process(s string, ctx context.Context) {
    // ...
}
```

```golang
// context.Context appears after multiple other parameters
func handle(s string, r int, ctx context.Context, x int) {
    // ...
}
```

### Valid

```golang
// context.Context is the first parameter
func process(ctx context.Context, s string) {
    // ...
}
```

```golang
// Only one parameter, no ordering issue
func process(ctx context.Context) {
    // ...
}
```
