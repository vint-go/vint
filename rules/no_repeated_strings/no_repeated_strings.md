---
title: noRepeatedStrings
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noRepeatedStrings`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noRepeatedStrings:
    min-occurrences: 3
    min-length: 3
    ignore-strings: ""
    ignore-tests: true
    ignore: ""
```

## Details

Detects repeated string literals across Go source files that could be replaced by a named constant.

When the same string literal appears multiple times in the codebase, it becomes a maintenance burden. If the value needs to change, every occurrence must be found and updated, which is error-prone. Extracting the string into a constant provides a single source of truth, makes the intent clearer, and reduces the risk of typos or inconsistent values.

By default, goconst reports strings that appear at least 2 times and are at least 3 characters long. These thresholds can be configured via `min-occurrences` and `min-length` options respectively. You can also exclude specific strings using the `ignore-strings` regex option, exclude test files with `ignore-tests`, or exclude files matching a pattern with `ignore`.

Source: https://github.com/jgautheron/goconst

## Examples

### Invalid

```golang
// The string "user_status" appears multiple times and should be a constant.
func GetUser(db *sql.DB) (*User, error) {
    row := db.QueryRow("SELECT * FROM users WHERE status = 'user_status'")
    // ...
    return nil, nil
}

func UpdateUser(db *sql.DB) error {
    _, err := db.Exec("UPDATE users SET status = 'user_status'")
    // ...
    return err
}

func DeleteUser(db *sql.DB) error {
    _, err := db.Exec("DELETE FROM users WHERE status = 'user_status'")
    // ...
    return err
}
```

```golang
// The string "application/json" appears multiple times.
func handleRequest(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    // ...
}

func handleResponse(w http.ResponseWriter) {
    w.Header().Set("Content-Type", "application/json")
    // ...
}

func handleError(w http.ResponseWriter) {
    w.Header().Set("Content-Type", "application/json")
    // ...
}
```

### Valid

```golang
// The repeated string is extracted into a constant.
const contentTypeJSON = "application/json"

func handleRequest(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", contentTypeJSON)
    // ...
}

func handleResponse(w http.ResponseWriter) {
    w.Header().Set("Content-Type", contentTypeJSON)
    // ...
}

func handleError(w http.ResponseWriter) {
    w.Header().Set("Content-Type", contentTypeJSON)
    // ...
}
```

```golang
// Short strings below the minimum length threshold are acceptable.
func example() {
    a := "ok"
    b := "ok"
    c := "ok"
}
```

```golang
// A string that only appears once does not trigger the rule.
func greet() string {
    return "Hello, welcome to the application!"
}
```
