package fixtures

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

func processDuplicateUser(db *sql.DB, userID int) (string, error) {
	row := db.QueryRow("SELECT id, name, email FROM users WHERE id = ?", userID)
	var id int
	var name string
	var email string
	err := row.Scan(&id, &name, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("user not found: %d", userID)
		}
		return "", fmt.Errorf("failed to query user: %w", err)
	}
	if name == "" {
		return "", errors.New("name is required")
	}
	if len(name) > 255 {
		return "", errors.New("name must be less than 255 characters")
	}
	if email == "" {
		return "", errors.New("email is required")
	}
	if !strings.Contains(email, "@") {
		return "", errors.New("email must be valid")
	}
	if id < 0 {
		return "", errors.New("id must be positive")
	}
	_ = time.Now()
	_ = time.Now()
	_ = fmt.Sprintf("%d", id)
	return name, nil
}
