---
title: noSqlInjectionTaint
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noSqlInjectionTaint`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noSqlInjectionTaint:
    # rule options here
```

## Details

Detects SQL injection vulnerabilities via taint analysis.

This rule uses taint analysis to trace the flow of untrusted data from sources (such as HTTP request parameters, user input, or external data) through the program to SQL query execution sinks. Unlike noSqlFormatString and noSqlConcatenation which use AST pattern matching, this rule performs interprocedural data flow analysis to detect SQL injection vulnerabilities that span multiple functions or packages.

Taint analysis can identify cases where user input flows through multiple transformations and function calls before reaching a SQL query, even when the direct construction of the query string is not immediately obvious from the source code.

SQL injection is one of the most critical web application vulnerabilities. All user-supplied data must be passed as parameterized query arguments rather than concatenated or formatted into query strings.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
func getUser(db *sql.DB, r *http.Request) (*User, error) {
    name := r.URL.Query().Get("name")
    // Taint flows from HTTP request to SQL query
    query := buildQuery(name)
    row := db.QueryRow(query)
    // ...
}

func buildQuery(name string) string {
    return "SELECT * FROM users WHERE name = '" + name + "'"
}
```

```golang
func handler(db *sql.DB, r *http.Request) {
    id := r.FormValue("id")
    result := fetchRecord(db, id)
    // ...
}

func fetchRecord(db *sql.DB, id string) *sql.Row {
    return db.QueryRow(fmt.Sprintf("SELECT * FROM records WHERE id = %s", id))
}
```

### Valid

```golang
func getUser(db *sql.DB, r *http.Request) (*User, error) {
    name := r.URL.Query().Get("name")
    // Parameterized query prevents SQL injection
    row := db.QueryRow("SELECT * FROM users WHERE name = $1", name)
    // ...
}
```

```golang
func handler(db *sql.DB, r *http.Request) {
    id := r.FormValue("id")
    // Using prepared statement
    stmt, _ := db.Prepare("SELECT * FROM records WHERE id = ?")
    row := stmt.QueryRow(id)
    // ...
}
```
