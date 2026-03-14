package fixtures

import (
	"database/sql"
	"fmt"
	"net/http"
)

// Invalid: Taint flows from HTTP request through function call to SQL query
func getUser(w http.ResponseWriter, r *http.Request) {
	db := getDB()
	name := r.URL.Query().Get("name")
	// Taint flows from HTTP request to SQL query via function call
	query := buildQuery(name)
	db.QueryRow(query) // MATCH /potential SQL injection: tainted data from user input flows to SQL query/
}

func buildQuery(name string) string {
	return "SELECT * FROM users WHERE name = '" + name + "'"
}

// Invalid: Taint flows from FormValue through fmt.Sprintf to SQL query in another function
func handler(w http.ResponseWriter, r *http.Request) {
	db := getDB()
	id := r.FormValue("id")
	fetchRecord(db, id) // MATCH /potential SQL injection: tainted data from user input flows to SQL query/
}

func fetchRecord(db *sql.DB, id string) *sql.Row {
	return db.QueryRow(fmt.Sprintf("SELECT * FROM records WHERE id = %s", id))
}

// Invalid: Direct taint from HTTP request to SQL query via string concatenation
func handlerConcat(w http.ResponseWriter, r *http.Request) {
	db := getDB()
	name := r.URL.Query().Get("name")
	query := "SELECT * FROM users WHERE name = '" + name + "'"
	db.Query(query) // MATCH /potential SQL injection: tainted data from user input flows to SQL query/
}

// Invalid: Taint flows through fmt.Sprintf to SQL Exec
func handlerExec(w http.ResponseWriter, r *http.Request) {
	db := getDB()
	table := r.FormValue("table")
	query := fmt.Sprintf("DELETE FROM %s WHERE id = 1", table)
	db.Exec(query) // MATCH /potential SQL injection: tainted data from user input flows to SQL query/
}

// Valid: Parameterized query prevents SQL injection
func getUserSafe(w http.ResponseWriter, r *http.Request) {
	db := getDB()
	name := r.URL.Query().Get("name")
	db.QueryRow("SELECT * FROM users WHERE name = $1", name)
}

// Valid: Using prepared statement
func handlerSafe(w http.ResponseWriter, r *http.Request) {
	db := getDB()
	id := r.FormValue("id")
	stmt, _ := db.Prepare("SELECT * FROM records WHERE id = ?")
	stmt.QueryRow(id)
}

// Valid: Constant query string with no user input
func getAllUsers() {
	db := getDB()
	db.Query("SELECT * FROM users")
}

// Valid: No HTTP input, just normal function parameter
func getByName(db *sql.DB, name string) *sql.Row {
	return db.QueryRow("SELECT * FROM users WHERE name = $1", name)
}

func getDB() *sql.DB {
	return nil
}
