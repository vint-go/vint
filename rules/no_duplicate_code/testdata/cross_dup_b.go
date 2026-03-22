package fixtures

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

func processDuplicateAdmin(db *sql.DB, adminID int) (string, error) {
	row := db.QueryRow("SELECT id, name, email FROM admins WHERE id = ?", adminID)
	var id int
	var name string
	var email string
	err := row.Scan(&id, &name, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("admin not found: %d", adminID)
		}
		return "", fmt.Errorf("failed to query admin: %w", err)
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
