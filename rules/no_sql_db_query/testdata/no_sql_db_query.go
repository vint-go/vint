package fixtures

import (
	"context"
	"database/sql"
	"log"
)

// Invalid: calling (*sql.DB).Query without a context
func badQuery(db *sql.DB) {
	rows, err := db.Query("SELECT id, name FROM users") // MATCH /(*sql.DB).Query does not accept a context; use (*sql.DB).QueryContext instead/
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
}

// Valid: calling (*sql.DB).QueryContext with a context
func goodQueryContext(db *sql.DB) {
	ctx := context.Background()
	rows, err := db.QueryContext(ctx, "SELECT id, name FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
}
