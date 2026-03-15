package fixtures

import (
	"context"
	"database/sql"
	"log"
)

// Invalid: calling (*sql.DB).Exec without a context
func badExec(db *sql.DB) {
	_, err := db.Exec("INSERT INTO users (name) VALUES ($1)", "Alice") // MATCH /(*sql.DB).Exec does not accept a context; use (*sql.DB).ExecContext instead/
	if err != nil {
		log.Fatal(err)
	}
}

// Valid: calling (*sql.DB).ExecContext with a context
func goodExecContext(db *sql.DB) {
	ctx := context.Background()
	_, err := db.ExecContext(ctx, "INSERT INTO users (name) VALUES ($1)", "Alice")
	if err != nil {
		log.Fatal(err)
	}
}
