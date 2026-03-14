package fixtures

import "database/sql"

// Invalid: SQL injection via string concatenation in QueryRow
func getUser(db *sql.DB, name string) (*sql.Row, error) {
	query := "SELECT * FROM users WHERE name = '" + name + "'"
	return db.QueryRow(query), nil // MATCH /SQL query constructed using string concatenation/
}

// Invalid: SQL injection via += concatenation in Query
func searchUsers(db *sql.DB, filter string) (*sql.Rows, error) {
	query := "SELECT * FROM users WHERE "
	query += filter
	return db.Query(query) // MATCH /SQL query constructed using string concatenation/
}

// Invalid: SQL injection via inline concatenation in QueryRow
func getByID(db *sql.DB, table string, id string) (*sql.Row, error) {
	query := "SELECT * FROM " + table + " WHERE id = " + id
	return db.QueryRow(query), nil // MATCH /SQL query constructed using string concatenation/
}

// Valid: parameterized query
func getUserSafe(db *sql.DB, name string) (*sql.Row, error) {
	return db.QueryRow("SELECT * FROM users WHERE name = ?", name), nil
}

// Valid: parameterized query with $1 placeholder
func getByIDSafe(db *sql.DB, id int64) (*sql.Row, error) {
	return db.QueryRow("SELECT * FROM users WHERE id = $1", id), nil
}

// Valid: constant string literal query
func getAllUsers(db *sql.DB) (*sql.Rows, error) {
	return db.Query("SELECT * FROM users")
}

// Valid: constant string concatenation (both sides are literals)
func getConstQuery(db *sql.DB) (*sql.Rows, error) {
	return db.Query("SELECT * FROM " + "users")
}
