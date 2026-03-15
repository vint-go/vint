package fixtures

import (
	"context"
	"database/sql"
	"log"
)

// Invalid: calling (*sql.Tx).Stmt without a context
func badTxStmt(db *sql.DB) {
	ctx := context.Background()
	stmt, err := db.PrepareContext(ctx, "SELECT * FROM users WHERE id = $1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	txStmt := tx.Stmt(stmt) // MATCH /(*sql.Tx).Stmt does not accept a context; use (*sql.Tx).StmtContext instead/
	_ = txStmt
	tx.Commit()
}

// Valid: calling (*sql.Tx).StmtContext with a context
func goodTxStmtContext(db *sql.DB) {
	ctx := context.Background()
	stmt, err := db.PrepareContext(ctx, "SELECT * FROM users WHERE id = $1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	txStmt := tx.StmtContext(ctx, stmt)
	_ = txStmt
	tx.Commit()
}
