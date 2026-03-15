package fixtures

import (
	"context"
	"database/sql"
	"log"
)

// Invalid: calling (*sql.DB).Begin without a context
func badBegin(db *sql.DB) {
	tx, err := db.Begin() // MATCH /(*sql.DB).Begin does not accept a context; use (*sql.DB).BeginTx instead/
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()
}

// Valid: calling (*sql.DB).BeginTx with a context
func goodBeginTx(db *sql.DB) {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()
}
