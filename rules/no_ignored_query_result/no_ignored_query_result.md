---
title: noIgnoredQueryResult
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noIgnoredQueryResult`
- This rule is not recommended (experimental).
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noIgnoredQueryResult:
    # no additional options
```

## Details

Detects issues in `Query()` and `Exec()` calls. This checker identifies problematic SQL query patterns where `Query()` or `QueryContext()` methods are called but their results are ignored (assigned to `_`). This can cause:

1. **Connection leak risk** -- When `Query()` returns result rows that are never read or closed, it can cause connection leaks in the database pool.
2. **Wrong method usage** -- When the underlying type has an `Exec()` method available, `Exec()` should be used instead of `Query()` for operations that don't need row results (like UPDATE or DELETE statements).

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
_, err := db.Query("UPDATE users SET active = true WHERE id = ?", id)
```

### Valid

```golang
_, err := db.Exec("UPDATE users SET active = true WHERE id = ?", id)
```

```golang
rows, err := db.Query("SELECT * FROM users WHERE active = true")
if err != nil {
    return err
}
defer rows.Close()
```
