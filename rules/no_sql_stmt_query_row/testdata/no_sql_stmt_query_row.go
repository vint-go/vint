package fixtures

import (
	"context"
	"database/sql"
	"log"
)

// Invalid: calling (*sql.Stmt).QueryRow without a context
func badStmtQueryRow(db *sql.DB) {
	ctx := context.Background()
	stmt, err := db.PrepareContext(ctx, "SELECT name FROM users WHERE id = $1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var name string
	err = stmt.QueryRow(1).Scan(&name) // MATCH /(*sql.Stmt).QueryRow does not accept a context; use (*sql.Stmt).QueryRowContext instead/
	if err != nil {
		log.Fatal(err)
	}
}

// Valid: calling (*sql.Stmt).QueryRowContext with a context
func goodStmtQueryRowContext(db *sql.DB) {
	ctx := context.Background()
	stmt, err := db.PrepareContext(ctx, "SELECT name FROM users WHERE id = $1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var name string
	err = stmt.QueryRowContext(ctx, 1).Scan(&name)
	if err != nil {
		log.Fatal(err)
	}
}
