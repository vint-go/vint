package fixtures

import (
	"context"
	"database/sql"
	"log"
)

// Invalid: calling (*sql.Tx).Exec without a context
func badTxExec(db *sql.DB) {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	_, err = tx.Exec("INSERT INTO users (name) VALUES ($1)", "Alice") // MATCH /(*sql.Tx).Exec does not accept a context; use (*sql.Tx).ExecContext instead/
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	tx.Commit()
}

// Valid: calling (*sql.Tx).ExecContext with a context
func goodTxExecContext(db *sql.DB) {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	_, err = tx.ExecContext(ctx, "INSERT INTO users (name) VALUES ($1)", "Alice")
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	tx.Commit()
}
