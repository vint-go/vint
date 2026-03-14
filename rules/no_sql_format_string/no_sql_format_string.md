---
title: noSqlFormatString
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noSqlFormatString`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noSqlFormatString:
    # rule options here
```

## Details

Detects SQL query construction using format string functions, which can lead to SQL injection vulnerabilities.

This rule identifies cases where SQL queries are built using `fmt.Sprintf`, `fmt.Sprint`, `fmt.Fprintf`, or similar formatting functions with unsanitized user input. When user-controlled data is interpolated into SQL query strings through format functions, attackers can inject arbitrary SQL commands, potentially reading, modifying, or deleting database contents.

The rule monitors calls to `database/sql` methods including `Exec`, `ExecContext`, `Query`, `QueryContext`, `QueryRow`, `QueryRowContext`, `Prepare`, and `PrepareContext`. It checks whether the query argument was constructed via string formatting with non-constant values.

SQL queries should always use parameterized queries (prepared statements) to safely incorporate user input.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
import (
    "database/sql"
    "fmt"
)

func getUser(db *sql.DB, name string) (*sql.Row, error) {
    // SQL injection via fmt.Sprintf
    query := fmt.Sprintf("SELECT * FROM users WHERE name = '%s'", name)
    return db.QueryRow(query), nil
}
```

```golang
import (
    "database/sql"
    "fmt"
)

func deleteUser(db *sql.DB, id string) error {
    query := fmt.Sprintf("DELETE FROM users WHERE id = %s", id)
    _, err := db.Exec(query)
    return err
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

func deleteUser(db *sql.DB, id int64) error {
    _, err := db.Exec("DELETE FROM users WHERE id = $1", id)
    return err
}
```

```golang
import "database/sql"

func getUsers(db *sql.DB, status string) (*sql.Rows, error) {
    // Using prepared statement
    stmt, err := db.Prepare("SELECT * FROM users WHERE status = ?")
    if err != nil {
        return nil, err
    }
    defer stmt.Close()
    return stmt.Query(status)
}
```
