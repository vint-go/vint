package fixtures

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

// Invalid: calling (*sql.Tx).QueryRow without a context
func badTxQueryRow(db *sql.DB) {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	var name string
	err = tx.QueryRow("SELECT name FROM users WHERE id = $1", 1).Scan(&name) // MATCH /(*sql.Tx).QueryRow does not accept a context; use (*sql.Tx).QueryRowContext instead/
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	fmt.Println(name)
	tx.Commit()
}

// Valid: calling (*sql.Tx).QueryRowContext with a context
func goodTxQueryRowContext(db *sql.DB) {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	var name string
	err = tx.QueryRowContext(ctx, "SELECT name FROM users WHERE id = $1", 1).Scan(&name)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	fmt.Println(name)
	tx.Commit()
}
