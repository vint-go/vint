package fixtures

import (
	"context"
	"database/sql"
	"log"
)

// Invalid: calling (*sql.Tx).Prepare without a context
func badTxPrepare(db *sql.DB) {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("SELECT * FROM users WHERE id = $1") // MATCH /(*sql.Tx).Prepare does not accept a context; use (*sql.Tx).PrepareContext instead/
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	defer stmt.Close()
	tx.Commit()
}

// Valid: calling (*sql.Tx).PrepareContext with a context
func goodTxPrepareContext(db *sql.DB) {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.PrepareContext(ctx, "SELECT * FROM users WHERE id = $1")
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	defer stmt.Close()
	tx.Commit()
}
