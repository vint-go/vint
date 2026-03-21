# dupl: duplicate detection limited to single file

## Affected rule
`lint/complexity/noDuplicateCode` (maps to golangci-lint `dupl`)

## Behavior in golangci-lint
The original `dupl` linter (github.com/mibk/dupl) uses suffix tree analysis across all files in a package to detect structurally identical code blocks. It can find duplicated functions even when they appear in different files.

## Behavior in vint
The `noDuplicateCode` rule's `Apply` method processes a single `*lint.File` at a time and compares function declarations only within that file. It uses a rolling hash approach on serialized AST tokens.

## Gap
Cross-file duplicate detection is not supported. If two identical functions exist in separate files (e.g. `handlers.go` and `utils.go`), the original `dupl` linter would flag them, but vint's `noDuplicateCode` rule will not.

## Example
```go
// file: handlers.go
func ProcessUser(db *sql.DB, id int) (*User, error) {
    row := db.QueryRow("SELECT * FROM users WHERE id = ?", id)
    var u User
    err := row.Scan(&u.ID, &u.Name)
    if err != nil { return nil, err }
    return &u, nil
}

// file: admin.go
func ProcessAdmin(db *sql.DB, id int) (*Admin, error) {
    row := db.QueryRow("SELECT * FROM admins WHERE id = ?", id)
    var a Admin
    err := row.Scan(&a.ID, &a.Name)
    if err != nil { return nil, err }
    return &a, nil
}
```

golangci-lint's `dupl` would flag these as duplicates. Vint's `noDuplicateCode` would not, since they are in different files.

## Impact on migration
Users who rely on `dupl` to catch cross-file code duplication will lose that coverage after migration. Only duplicates within the same file will be detected by vint.
