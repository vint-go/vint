package fixtures

import (
	"context"
	"database/sql"
	"log"
)

// Invalid: calling (*sql.DB).Ping without a context
func badPing(db *sql.DB) {
	if err := db.Ping(); err != nil { // MATCH /(*sql.DB).Ping does not accept a context; use (*sql.DB).PingContext instead/
		log.Fatal(err)
	}
}

// Valid: calling (*sql.DB).PingContext with a context
func goodPingContext(db *sql.DB) {
	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		log.Fatal(err)
	}
}
