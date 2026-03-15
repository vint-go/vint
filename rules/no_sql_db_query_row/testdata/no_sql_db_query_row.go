package fixtures

import (
	"context"
	"database/sql"
	"log"
)

// Invalid: calling (*sql.DB).QueryRow without a context
func badQueryRow(db *sql.DB) {
	var name string
	err := db.QueryRow("SELECT name FROM users WHERE id = $1", 1).Scan(&name) // MATCH /(*sql.DB).QueryRow does not accept a context; use (*sql.DB).QueryRowContext instead/
	if err != nil {
		log.Fatal(err)
	}
}

// Valid: calling (*sql.DB).QueryRowContext with a context
func goodQueryRowContext(db *sql.DB) {
	ctx := context.Background()
	var name string
	err := db.QueryRowContext(ctx, "SELECT name FROM users WHERE id = $1", 1).Scan(&name)
	if err != nil {
		log.Fatal(err)
	}
}
