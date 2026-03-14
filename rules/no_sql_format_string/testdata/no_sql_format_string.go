package fixtures

import (
	"database/sql"
	"fmt"
)

// Invalid: SQL injection via fmt.Sprintf in QueryRow
func getUser(db *sql.DB, name string) (*sql.Row, error) {
	query := fmt.Sprintf("SELECT * FROM users WHERE name = '%s'", name)
	return db.QueryRow(query), nil // MATCH /SQL query constructed using format string function/
}

// Invalid: SQL injection via fmt.Sprintf in Exec
func deleteUser(db *sql.DB, id string) error {
	query := fmt.Sprintf("DELETE FROM users WHERE id = %s", id)
	_, err := db.Exec(query) // MATCH /SQL query constructed using format string function/
	return err
}

// Invalid: SQL injection via inline fmt.Sprintf in QueryRow
func getUserInline(db *sql.DB, name string) (*sql.Row, error) {
	return db.QueryRow(fmt.Sprintf("SELECT * FROM users WHERE name = '%s'", name)), nil // MATCH /SQL query constructed using format string function/
}

// Valid: parameterized query
func getUserSafe(db *sql.DB, name string) (*sql.Row, error) {
	return db.QueryRow("SELECT * FROM users WHERE name = ?", name), nil
}

// Valid: parameterized query with $1 placeholder
func deleteUserSafe(db *sql.DB, id int64) error {
	_, err := db.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}

// Valid: using prepared statement
func getUsersPrepared(db *sql.DB, status string) (*sql.Rows, error) {
	stmt, err := db.Prepare("SELECT * FROM users WHERE status = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Query(status)
}

// Valid: constant string literal query
func getAllUsers(db *sql.DB) (*sql.Rows, error) {
	return db.Query("SELECT * FROM users")
}

// Valid: fmt.Sprintf with only constant arguments
func getConstQuery(db *sql.DB) (*sql.Rows, error) {
	query := fmt.Sprintf("SELECT * FROM %s", "users")
	return db.Query(query)
}
