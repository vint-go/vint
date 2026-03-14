package fixtures

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

// The following two functions contain duplicate logic that should be
// extracted into a shared helper. They are structurally identical.

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

func processDuplicateAdmin(db *sql.DB, adminID int) (string, error) { // MATCH /duplicate code detected: processDuplicateUser and processDuplicateAdmin share 153 tokens of identical structure/
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

// These two functions are different enough that they should not trigger the rule.

func uniqueFunction1() string {
	x := 1
	y := 2
	return fmt.Sprintf("%d+%d=%d", x, y, x+y)
}

func uniqueFunction2(a, b int) error {
	if a > b {
		return errors.New("a must be less than b")
	}
	return nil
}
