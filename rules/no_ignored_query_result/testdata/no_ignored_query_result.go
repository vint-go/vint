package fixtures

import (
	"context"
	"database/sql"
)

func ignoredQueryOnDB(db *sql.DB, id int) error {
	_, err := db.Query("UPDATE users SET active = true WHERE id = ?", id) // MATCH /query result rows are ignored; use Exec() instead of Query() when rows are not needed/
	return err
}

func ignoredQueryContextOnDB(db *sql.DB, ctx context.Context, id int) error {
	_, err := db.QueryContext(ctx, "UPDATE users SET active = true WHERE id = ?", id) // MATCH /query result rows are ignored; use Exec() instead of QueryContext() when rows are not needed/
	return err
}

func ignoredQueryOnTx(tx *sql.Tx, id int) error {
	_, err := tx.Query("DELETE FROM users WHERE id = ?", id) // MATCH /query result rows are ignored; use Exec() instead of Query() when rows are not needed/
	return err
}

func ignoredQueryContextOnTx(tx *sql.Tx, ctx context.Context, id int) error {
	_, err := tx.QueryContext(ctx, "DELETE FROM users WHERE id = ?", id) // MATCH /query result rows are ignored; use Exec() instead of QueryContext() when rows are not needed/
	return err
}

func ignoredQueryOnStmt(stmt *sql.Stmt, id int) error {
	_, err := stmt.Query(id) // MATCH /query result rows are ignored; use Exec() instead of Query() when rows are not needed/
	return err
}

func ignoredQueryContextOnStmt(stmt *sql.Stmt, ctx context.Context, id int) error {
	_, err := stmt.QueryContext(ctx, id) // MATCH /query result rows are ignored; use Exec() instead of QueryContext() when rows are not needed/
	return err
}

// Valid: using Exec instead of Query
func validExecOnDB(db *sql.DB, id int) error {
	_, err := db.Exec("UPDATE users SET active = true WHERE id = ?", id)
	return err
}

// Valid: Query result is actually used
func validQueryUsed(db *sql.DB) error {
	rows, err := db.Query("SELECT * FROM users WHERE active = true")
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

// Valid: QueryContext result is actually used
func validQueryContextUsed(db *sql.DB, ctx context.Context) error {
	rows, err := db.QueryContext(ctx, "SELECT * FROM users WHERE active = true")
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

// Valid: ExecContext is fine
func validExecContext(db *sql.DB, ctx context.Context, id int) error {
	_, err := db.ExecContext(ctx, "UPDATE users SET active = true WHERE id = ?", id)
	return err
}
