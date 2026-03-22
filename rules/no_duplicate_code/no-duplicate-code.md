---
title: noDuplicateCode
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/complexity/noDuplicateCode`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/complexity/noDuplicateCode:
    threshold: 150
```

## Details

Detects duplicate fragments of code across Go source files.

The `dupl` linter uses suffix tree analysis on serialized abstract syntax trees (ASTs) to find structurally identical code blocks, regardless of specific variable names or literal values. For example, `if a == 13 {}` and `if x == 100 {}` are considered structurally identical because they share the same AST shape.

The `threshold` option controls the minimum number of tokens a duplicated code sequence must contain before it is reported. The default threshold in golangci-lint is `150`. Lowering this value will report smaller duplicated fragments, while raising it will only flag larger blocks of copied code.

Because the analysis ignores concrete values and focuses on structure, it may produce false positives -- matches that a developer would not consider true duplicates due to their small size or differing literal values. Manual review of reported duplicates is recommended.

When duplicate code is detected, refactoring the common logic into a shared function or method is typically the correct approach. This reduces maintenance burden, decreases the risk of divergent bug fixes, and improves overall code quality.

Source: https://github.com/mibk/dupl

## Examples

### Invalid

```golang
// File: handlers.go
// The following two functions contain duplicate logic that should be
// extracted into a shared helper.

func ProcessUser(db *sql.DB, userID int) (*User, error) {
    row := db.QueryRow("SELECT id, name, email FROM users WHERE id = ?", userID)
    var user User
    err := row.Scan(&user.ID, &user.Name, &user.Email)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("user not found: %d", userID)
        }
        return nil, fmt.Errorf("failed to query user: %w", err)
    }
    user.CreatedAt = time.Now()
    user.UpdatedAt = time.Now()
    return &user, nil
}

func ProcessAdmin(db *sql.DB, adminID int) (*Admin, error) {
    row := db.QueryRow("SELECT id, name, email FROM admins WHERE id = ?", adminID)
    var admin Admin
    err := row.Scan(&admin.ID, &admin.Name, &admin.Email)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("admin not found: %d", adminID)
        }
        return nil, fmt.Errorf("failed to query admin: %w", err)
    }
    admin.CreatedAt = time.Now()
    admin.UpdatedAt = time.Now()
    return &admin, nil
}
```

```golang
// File: validation.go
// These two validation functions share the same structure and logic,
// differing only in field names and error messages.

func ValidateCreateRequest(req CreateRequest) error {
    if req.Name == "" {
        return errors.New("name is required")
    }
    if len(req.Name) > 255 {
        return errors.New("name must be less than 255 characters")
    }
    if req.Email == "" {
        return errors.New("email is required")
    }
    if !strings.Contains(req.Email, "@") {
        return errors.New("email must be valid")
    }
    if req.Age < 0 || req.Age > 150 {
        return errors.New("age must be between 0 and 150")
    }
    return nil
}

func ValidateUpdateRequest(req UpdateRequest) error {
    if req.Name == "" {
        return errors.New("name is required")
    }
    if len(req.Name) > 255 {
        return errors.New("name must be less than 255 characters")
    }
    if req.Email == "" {
        return errors.New("email is required")
    }
    if !strings.Contains(req.Email, "@") {
        return errors.New("email must be valid")
    }
    if req.Age < 0 || req.Age > 150 {
        return errors.New("age must be between 0 and 150")
    }
    return nil
}
```

### Valid

```golang
// Shared helper eliminates code duplication.

func queryEntity(db *sql.DB, query string, id int, dest ...interface{}) error {
    row := db.QueryRow(query, id)
    err := row.Scan(dest...)
    if err != nil {
        if err == sql.ErrNoRows {
            return fmt.Errorf("entity not found: %d", id)
        }
        return fmt.Errorf("failed to query entity: %w", err)
    }
    return nil
}

func ProcessUser(db *sql.DB, userID int) (*User, error) {
    var user User
    err := queryEntity(db, "SELECT id, name, email FROM users WHERE id = ?", userID,
        &user.ID, &user.Name, &user.Email)
    if err != nil {
        return nil, err
    }
    user.CreatedAt = time.Now()
    user.UpdatedAt = time.Now()
    return &user, nil
}

func ProcessAdmin(db *sql.DB, adminID int) (*Admin, error) {
    var admin Admin
    err := queryEntity(db, "SELECT id, name, email FROM admins WHERE id = ?", adminID,
        &admin.ID, &admin.Name, &admin.Email)
    if err != nil {
        return nil, err
    }
    admin.CreatedAt = time.Now()
    admin.UpdatedAt = time.Now()
    return &admin, nil
}
```

```golang
// Generic validation function eliminates structural duplication.

type Validatable interface {
    GetName() string
    GetEmail() string
    GetAge() int
}

func validateCommonFields(v Validatable) error {
    if v.GetName() == "" {
        return errors.New("name is required")
    }
    if len(v.GetName()) > 255 {
        return errors.New("name must be less than 255 characters")
    }
    if v.GetEmail() == "" {
        return errors.New("email is required")
    }
    if !strings.Contains(v.GetEmail(), "@") {
        return errors.New("email must be valid")
    }
    if v.GetAge() < 0 || v.GetAge() > 150 {
        return errors.New("age must be between 0 and 150")
    }
    return nil
}

func ValidateCreateRequest(req CreateRequest) error {
    return validateCommonFields(req)
}

func ValidateUpdateRequest(req UpdateRequest) error {
    return validateCommonFields(req)
}
```
