package fixtures

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

// Invalid: calling (*sql.Tx).Query without a context
func badTxQuery(db *sql.DB) {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := tx.Query("SELECT id, name FROM users") // MATCH /(*sql.Tx).Query does not accept a context; use (*sql.Tx).QueryContext instead/
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		fmt.Println(id, name)
	}
	tx.Commit()
}

// Valid: calling (*sql.Tx).QueryContext with a context
func goodTxQueryContext(db *sql.DB) {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := tx.QueryContext(ctx, "SELECT id, name FROM users")
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		fmt.Println(id, name)
	}
	tx.Commit()
}
