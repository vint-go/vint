---
title: noSqlConcatenation
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noSqlConcatenation`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noSqlConcatenation:
    # rule options here
```

## Details

Detects SQL query construction using string concatenation, which can lead to SQL injection vulnerabilities.

This rule identifies cases where SQL queries are built by concatenating strings with the `+` operator or by using `+=` to append user-controlled values to query strings. When variable data is concatenated directly into SQL queries, attackers can inject arbitrary SQL commands.

The rule monitors calls to `database/sql` methods (`Exec`, `Query`, `Prepare`, and their context variants) and tracks variable mutations through assignment and concatenation operators to determine if the final query string contains non-constant parts.

SQL queries should always use parameterized queries (prepared statements) to safely incorporate dynamic values.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
import "database/sql"

func getUser(db *sql.DB, name string) (*sql.Row, error) {
    // SQL injection via string concatenation
    query := "SELECT * FROM users WHERE name = '" + name + "'"
    return db.QueryRow(query), nil
}
```

```golang
import "database/sql"

func searchUsers(db *sql.DB, filter string) (*sql.Rows, error) {
    query := "SELECT * FROM users WHERE "
    query += filter // User-controlled filter appended
    return db.Query(query)
}
```

```golang
import "database/sql"

func getByID(db *sql.DB, table string, id string) (*sql.Row, error) {
    query := "SELECT * FROM " + table + " WHERE id = " + id
    return db.QueryRow(query), nil
}
```

### Valid

```golang
import "database/sql"

func getUser(db *sql.DB, name string) (*sql.Row, error) {
    // Using parameterized query
    return db.QueryRow("SELECT * FROM users WHERE name = ?", name), nil
}
```

```golang
import "database/sql"

func getByID(db *sql.DB, id int64) (*sql.Row, error) {
    return db.QueryRow("SELECT * FROM users WHERE id = $1", id), nil
}
```

```golang
import (
    "database/sql"
    "github.com/lib/pq"
)

func getByColumn(db *sql.DB, column string, value string) (*sql.Row, error) {
    // Using pq.QuoteIdentifier for safe identifier quoting
    query := "SELECT * FROM users WHERE " + pq.QuoteIdentifier(column) + " = $1"
    return db.QueryRow(query, value), nil
}
```
