package fixtures

import (
	"context"
	"database/sql"
	"log"
)

// Invalid: calling (*sql.Stmt).Query without a context
func badStmtQuery(db *sql.DB) {
	ctx := context.Background()
	stmt, err := db.PrepareContext(ctx, "SELECT id, name FROM users WHERE active = $1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(true) // MATCH /(*sql.Stmt).Query does not accept a context; use (*sql.Stmt).QueryContext instead/
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
}

// Valid: calling (*sql.Stmt).QueryContext with a context
func goodStmtQueryContext(db *sql.DB) {
	ctx := context.Background()
	stmt, err := db.PrepareContext(ctx, "SELECT id, name FROM users WHERE active = $1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, true)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
}
