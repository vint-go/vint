package fixtures

import (
	"context"
	"database/sql"
	"log"
)

// Invalid: calling (*sql.Stmt).Exec without a context
func badStmtExec(db *sql.DB) {
	ctx := context.Background()
	stmt, err := db.PrepareContext(ctx, "INSERT INTO users (name) VALUES ($1)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec("Alice") // MATCH /(*sql.Stmt).Exec does not accept a context; use (*sql.Stmt).ExecContext instead/
	if err != nil {
		log.Fatal(err)
	}
}

// Valid: calling (*sql.Stmt).ExecContext with a context
func goodStmtExecContext(db *sql.DB) {
	ctx := context.Background()
	stmt, err := db.PrepareContext(ctx, "INSERT INTO users (name) VALUES ($1)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, "Alice")
	if err != nil {
		log.Fatal(err)
	}
}
