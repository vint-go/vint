package fixtures

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func noBlankErrorAssignmentBadOpen() {
	f, _ := os.Open("file.txt") // MATCH /error assigned to blank identifier/
	_ = f.Close()               // MATCH /error assigned to blank identifier/
}

func noBlankErrorAssignmentBadEncode() {
	enc := json.NewEncoder(os.Stdout)
	_ = enc.Encode(map[string]string{"key": "value"}) // MATCH /error assigned to blank identifier/
}

func noBlankErrorAssignmentBadCloseDB(db *sql.DB) {
	_ = db.Close() // MATCH /error assigned to blank identifier/
}

func noBlankErrorAssignmentBadFprintf() {
	_, _ = fmt.Fprintf(os.Stdout, "hello") // MATCH /error assigned to blank identifier/
}

func noBlankErrorAssignmentGoodOpen() {
	f, err := os.Open("file.txt")
	if err != nil {
		log.Fatal(err)
	}
	err = f.Close()
	if err != nil {
		log.Printf("failed to close file: %v", err)
	}
}

func noBlankErrorAssignmentGoodEncode() {
	enc := json.NewEncoder(os.Stdout)
	if err := enc.Encode(map[string]string{"key": "value"}); err != nil {
		log.Fatal(err)
	}
}

func noBlankErrorAssignmentGoodCloseDB(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Printf("failed to close database: %v", err)
	}
}

func noBlankErrorAssignmentGoodFprintf() {
	_, err := fmt.Fprintf(os.Stdout, "hello")
	if err != nil {
		log.Fatal(err)
	}
}

func noBlankErrorAssignmentBlankNonError() {
	// Assigning non-error to blank identifier is fine
	_ = 42
	_ = "hello"
}
