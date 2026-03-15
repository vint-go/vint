package fixtures

import (
	"context"
	"database/sql"
	"log"
)

// Invalid: calling (*sql.DB).Prepare without a context
func badPrepare(db *sql.DB) {
	stmt, err := db.Prepare("SELECT * FROM users WHERE id = $1") // MATCH /(*sql.DB).Prepare does not accept a context; use (*sql.DB).PrepareContext instead/
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
}

// Valid: calling (*sql.DB).PrepareContext with a context
func goodPrepareContext(db *sql.DB) {
	ctx := context.Background()
	stmt, err := db.PrepareContext(ctx, "SELECT * FROM users WHERE id = $1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
}
