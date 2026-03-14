---
title: noExcessiveStatements
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/complexity/noExcessiveStatements`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/complexity/noExcessiveStatements:
    statements: 40
```

## Details

Checks that functions do not exceed a maximum number of statements. Unlike a simple line count, the statement count measures the actual complexity of a function by counting executable statements, providing a metric that is independent of code formatting style.

The default maximum is 40 statements per function. Statements are counted recursively: the body of control flow structures (`for`, `range`, `if`, `switch`, `type switch`, `select`) and `case` clauses all contribute to the total count. Inline function literals used in assignments, `go` statements, and `defer` statements are also recursively parsed, and their statements are included in the parent function's total.

This check is evaluated before the line length check. If a function exceeds the statement limit, the line length check is skipped for that function.

The statement limit can be set to `-1` to disable this check entirely.

Source: https://github.com/ultraware/funlen

## Examples

### Invalid

```golang
// This function has too many statements (exceeds 40).
func handleRequest(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    logger := log.FromContext(ctx)
    logger.Info("handling request")

    userID := r.Header.Get("X-User-ID")
    token := r.Header.Get("Authorization")
    contentType := r.Header.Get("Content-Type")
    requestID := r.Header.Get("X-Request-ID")

    if userID == "" {
        http.Error(w, "missing user ID", http.StatusBadRequest)
        return
    }
    if token == "" {
        http.Error(w, "missing token", http.StatusUnauthorized)
        return
    }

    user, err := fetchUser(ctx, userID)
    if err != nil {
        logger.Error("fetch user failed", "error", err)
        http.Error(w, "internal error", http.StatusInternalServerError)
        return
    }

    valid, err := validateToken(ctx, token, user)
    if err != nil {
        logger.Error("token validation failed", "error", err)
        http.Error(w, "internal error", http.StatusInternalServerError)
        return
    }
    if !valid {
        http.Error(w, "invalid token", http.StatusForbidden)
        return
    }

    body, err := io.ReadAll(r.Body)
    if err != nil {
        logger.Error("read body failed", "error", err)
        http.Error(w, "bad request", http.StatusBadRequest)
        return
    }

    var req RequestPayload
    err = json.Unmarshal(body, &req)
    if err != nil {
        logger.Error("unmarshal failed", "error", err)
        http.Error(w, "bad request", http.StatusBadRequest)
        return
    }

    result, err := processPayload(ctx, user, &req)
    if err != nil {
        logger.Error("process failed", "error", err)
        http.Error(w, "internal error", http.StatusInternalServerError)
        return
    }

    response, err := json.Marshal(result)
    if err != nil {
        logger.Error("marshal failed", "error", err)
        http.Error(w, "internal error", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("X-Request-ID", requestID)
    w.WriteHeader(http.StatusOK)
    w.Write(response)
}
```

```golang
// Inline function literals also contribute to the parent's statement count.
func setup() {
    a := 1
    b := 2
    c := 3
    // ... many more statements ...
    go func() {
        // These statements count toward the parent function's total.
        x := doWork()
        log.Println(x)
        cleanup(x)
    }()
    defer func() {
        // These statements also count toward the parent function's total.
        saveState()
        closeConnections()
    }()
    // ... more statements that push the total over 40 ...
}
```

### Valid

```golang
// This function has a small number of statements.
func handleRequest(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    user, err := authenticateRequest(r)
    if err != nil {
        writeError(w, err)
        return
    }

    payload, err := parsePayload(r)
    if err != nil {
        writeError(w, err)
        return
    }

    result, err := processPayload(ctx, user, payload)
    if err != nil {
        writeError(w, err)
        return
    }

    writeJSON(w, http.StatusOK, result)
}
```

```golang
// Breaking logic into smaller functions keeps statement count low.
func authenticateRequest(r *http.Request) (*User, error) {
    userID := r.Header.Get("X-User-ID")
    if userID == "" {
        return nil, ErrMissingUserID
    }

    token := r.Header.Get("Authorization")
    if token == "" {
        return nil, ErrMissingToken
    }

    user, err := fetchUser(r.Context(), userID)
    if err != nil {
        return nil, fmt.Errorf("fetch user: %w", err)
    }

    if err := validateToken(r.Context(), token, user); err != nil {
        return nil, fmt.Errorf("validate token: %w", err)
    }

    return user, nil
}
```
