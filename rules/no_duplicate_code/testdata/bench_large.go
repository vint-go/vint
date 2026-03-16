package fixtures

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

// This file is used for benchmarking the no_duplicate_code rule.
// It contains 30 structurally-identical functions to stress the O(n^2 * m^2) algorithm.

func benchFunc01(db *sql.DB, id int) (string, error) {
	row := db.QueryRow("SELECT id, name, email FROM t WHERE id = ?", id)
	var rid int
	var name string
	var email string
	err := row.Scan(&rid, &name, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("not found: %d", id)
		}
		return "", fmt.Errorf("query failed: %w", err)
	}
	if name == "" {
		return "", errors.New("name required")
	}
	if len(name) > 255 {
		return "", errors.New("name too long")
	}
	if email == "" {
		return "", errors.New("email required")
	}
	if !strings.Contains(email, "@") {
		return "", errors.New("email invalid")
	}
	if rid < 0 {
		return "", errors.New("id negative")
	}
	_ = time.Now()
	_ = time.Now()
	_ = fmt.Sprintf("%d", rid)
	return name, nil
}

func benchFunc02(db *sql.DB, id int) (string, error) {
	row := db.QueryRow("SELECT id, name, email FROM t WHERE id = ?", id)
	var rid int
	var name string
	var email string
	err := row.Scan(&rid, &name, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("not found: %d", id)
		}
		return "", fmt.Errorf("query failed: %w", err)
	}
	if name == "" {
		return "", errors.New("name required")
	}
	if len(name) > 255 {
		return "", errors.New("name too long")
	}
	if email == "" {
		return "", errors.New("email required")
	}
	if !strings.Contains(email, "@") {
		return "", errors.New("email invalid")
	}
	if rid < 0 {
		return "", errors.New("id negative")
	}
	_ = time.Now()
	_ = time.Now()
	_ = fmt.Sprintf("%d", rid)
	return name, nil
}

func benchFunc03(db *sql.DB, id int) (string, error) {
	row := db.QueryRow("SELECT id, name, email FROM t WHERE id = ?", id)
	var rid int
	var name string
	var email string
	err := row.Scan(&rid, &name, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("not found: %d", id)
		}
		return "", fmt.Errorf("query failed: %w", err)
	}
	if name == "" {
		return "", errors.New("name required")
	}
	if len(name) > 255 {
		return "", errors.New("name too long")
	}
	if email == "" {
		return "", errors.New("email required")
	}
	if !strings.Contains(email, "@") {
		return "", errors.New("email invalid")
	}
	if rid < 0 {
		return "", errors.New("id negative")
	}
	_ = time.Now()
	_ = time.Now()
	_ = fmt.Sprintf("%d", rid)
	return name, nil
}

func benchFunc04(db *sql.DB, id int) (string, error) {
	row := db.QueryRow("SELECT id, name, email FROM t WHERE id = ?", id)
	var rid int
	var name string
	var email string
	err := row.Scan(&rid, &name, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("not found: %d", id)
		}
		return "", fmt.Errorf("query failed: %w", err)
	}
	if name == "" {
		return "", errors.New("name required")
	}
	if len(name) > 255 {
		return "", errors.New("name too long")
	}
	if email == "" {
		return "", errors.New("email required")
	}
	if !strings.Contains(email, "@") {
		return "", errors.New("email invalid")
	}
	if rid < 0 {
		return "", errors.New("id negative")
	}
	_ = time.Now()
	_ = time.Now()
	_ = fmt.Sprintf("%d", rid)
	return name, nil
}

func benchFunc05(db *sql.DB, id int) (string, error) {
	row := db.QueryRow("SELECT id, name, email FROM t WHERE id = ?", id)
	var rid int
	var name string
	var email string
	err := row.Scan(&rid, &name, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("not found: %d", id)
		}
		return "", fmt.Errorf("query failed: %w", err)
	}
	if name == "" {
		return "", errors.New("name required")
	}
	if len(name) > 255 {
		return "", errors.New("name too long")
	}
	if email == "" {
		return "", errors.New("email required")
	}
	if !strings.Contains(email, "@") {
		return "", errors.New("email invalid")
	}
	if rid < 0 {
		return "", errors.New("id negative")
	}
	_ = time.Now()
	_ = time.Now()
	_ = fmt.Sprintf("%d", rid)
	return name, nil
}

func benchFunc06(db *sql.DB, id int) (string, error) {
	row := db.QueryRow("SELECT id, name, email FROM t WHERE id = ?", id)
	var rid int
	var name string
	var email string
	err := row.Scan(&rid, &name, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("not found: %d", id)
		}
		return "", fmt.Errorf("query failed: %w", err)
	}
	if name == "" {
		return "", errors.New("name required")
	}
	if len(name) > 255 {
		return "", errors.New("name too long")
	}
	if email == "" {
		return "", errors.New("email required")
	}
	if !strings.Contains(email, "@") {
		return "", errors.New("email invalid")
	}
	if rid < 0 {
		return "", errors.New("id negative")
	}
	_ = time.Now()
	_ = time.Now()
	_ = fmt.Sprintf("%d", rid)
	return name, nil
}

func benchFunc07(db *sql.DB, id int) (string, error) {
	row := db.QueryRow("SELECT id, name, email FROM t WHERE id = ?", id)
	var rid int
	var name string
	var email string
	err := row.Scan(&rid, &name, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("not found: %d", id)
		}
		return "", fmt.Errorf("query failed: %w", err)
	}
	if name == "" {
		return "", errors.New("name required")
	}
	if len(name) > 255 {
		return "", errors.New("name too long")
	}
	if email == "" {
		return "", errors.New("email required")
	}
	if !strings.Contains(email, "@") {
		return "", errors.New("email invalid")
	}
	if rid < 0 {
		return "", errors.New("id negative")
	}
	_ = time.Now()
	_ = time.Now()
	_ = fmt.Sprintf("%d", rid)
	return name, nil
}

func benchFunc08(db *sql.DB, id int) (string, error) {
	row := db.QueryRow("SELECT id, name, email FROM t WHERE id = ?", id)
	var rid int
	var name string
	var email string
	err := row.Scan(&rid, &name, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("not found: %d", id)
		}
		return "", fmt.Errorf("query failed: %w", err)
	}
	if name == "" {
		return "", errors.New("name required")
	}
	if len(name) > 255 {
		return "", errors.New("name too long")
	}
	if email == "" {
		return "", errors.New("email required")
	}
	if !strings.Contains(email, "@") {
		return "", errors.New("email invalid")
	}
	if rid < 0 {
		return "", errors.New("id negative")
	}
	_ = time.Now()
	_ = time.Now()
	_ = fmt.Sprintf("%d", rid)
	return name, nil
}

func benchFunc09(db *sql.DB, id int) (string, error) {
	row := db.QueryRow("SELECT id, name, email FROM t WHERE id = ?", id)
	var rid int
	var name string
	var email string
	err := row.Scan(&rid, &name, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("not found: %d", id)
		}
		return "", fmt.Errorf("query failed: %w", err)
	}
	if name == "" {
		return "", errors.New("name required")
	}
	if len(name) > 255 {
		return "", errors.New("name too long")
	}
	if email == "" {
		return "", errors.New("email required")
	}
	if !strings.Contains(email, "@") {
		return "", errors.New("email invalid")
	}
	if rid < 0 {
		return "", errors.New("id negative")
	}
	_ = time.Now()
	_ = time.Now()
	_ = fmt.Sprintf("%d", rid)
	return name, nil
}

func benchFunc10(db *sql.DB, id int) (string, error) {
	row := db.QueryRow("SELECT id, name, email FROM t WHERE id = ?", id)
	var rid int
	var name string
	var email string
	err := row.Scan(&rid, &name, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("not found: %d", id)
		}
		return "", fmt.Errorf("query failed: %w", err)
	}
	if name == "" {
		return "", errors.New("name required")
	}
	if len(name) > 255 {
		return "", errors.New("name too long")
	}
	if email == "" {
		return "", errors.New("email required")
	}
	if !strings.Contains(email, "@") {
		return "", errors.New("email invalid")
	}
	if rid < 0 {
		return "", errors.New("id negative")
	}
	_ = time.Now()
	_ = time.Now()
	_ = fmt.Sprintf("%d", rid)
	return name, nil
}

func benchFunc11(db *sql.DB, id int) (string, error) {
	row := db.QueryRow("SELECT id, name, email FROM t WHERE id = ?", id)
	var rid int
	var name string
	var email string
	err := row.Scan(&rid, &name, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("not found: %d", id)
		}
		return "", fmt.Errorf("query failed: %w", err)
	}
	if name == "" {
		return "", errors.New("name required")
	}
	if len(name) > 255 {
		return "", errors.New("name too long")
	}
	if email == "" {
		return "", errors.New("email required")
	}
	if !strings.Contains(email, "@") {
		return "", errors.New("email invalid")
	}
	if rid < 0 {
		return "", errors.New("id negative")
	}
	_ = time.Now()
	_ = time.Now()
	_ = fmt.Sprintf("%d", rid)
	return name, nil
}

func benchFunc12(db *sql.DB, id int) (string, error) {
	row := db.QueryRow("SELECT id, name, email FROM t WHERE id = ?", id)
	var rid int
	var name string
	var email string
	err := row.Scan(&rid, &name, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("not found: %d", id)
		}
		return "", fmt.Errorf("query failed: %w", err)
	}
	if name == "" {
		return "", errors.New("name required")
	}
	if len(name) > 255 {
		return "", errors.New("name too long")
	}
	if email == "" {
		return "", errors.New("email required")
	}
	if !strings.Contains(email, "@") {
		return "", errors.New("email invalid")
	}
	if rid < 0 {
		return "", errors.New("id negative")
	}
	_ = time.Now()
	_ = time.Now()
	_ = fmt.Sprintf("%d", rid)
	return name, nil
}

func benchFunc13(db *sql.DB, id int) (string, error) {
	row := db.QueryRow("SELECT id, name, email FROM t WHERE id = ?", id)
	var rid int
	var name string
	var email string
	err := row.Scan(&rid, &name, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("not found: %d", id)
		}
		return "", fmt.Errorf("query failed: %w", err)
	}
	if name == "" {
		return "", errors.New("name required")
	}
	if len(name) > 255 {
		return "", errors.New("name too long")
	}
	if email == "" {
		return "", errors.New("email required")
	}
	if !strings.Contains(email, "@") {
		return "", errors.New("email invalid")
	}
	if rid < 0 {
		return "", errors.New("id negative")
	}
	_ = time.Now()
	_ = time.Now()
	_ = fmt.Sprintf("%d", rid)
	return name, nil
}

func benchFunc14(db *sql.DB, id int) (string, error) {
	row := db.QueryRow("SELECT id, name, email FROM t WHERE id = ?", id)
	var rid int
	var name string
	var email string
	err := row.Scan(&rid, &name, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("not found: %d", id)
		}
		return "", fmt.Errorf("query failed: %w", err)
	}
	if name == "" {
		return "", errors.New("name required")
	}
	if len(name) > 255 {
		return "", errors.New("name too long")
	}
	if email == "" {
		return "", errors.New("email required")
	}
	if !strings.Contains(email, "@") {
		return "", errors.New("email invalid")
	}
	if rid < 0 {
		return "", errors.New("id negative")
	}
	_ = time.Now()
	_ = time.Now()
	_ = fmt.Sprintf("%d", rid)
	return name, nil
}

func benchFunc15(db *sql.DB, id int) (string, error) {
	row := db.QueryRow("SELECT id, name, email FROM t WHERE id = ?", id)
	var rid int
	var name string
	var email string
	err := row.Scan(&rid, &name, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("not found: %d", id)
		}
		return "", fmt.Errorf("query failed: %w", err)
	}
	if name == "" {
		return "", errors.New("name required")
	}
	if len(name) > 255 {
		return "", errors.New("name too long")
	}
	if email == "" {
		return "", errors.New("email required")
	}
	if !strings.Contains(email, "@") {
		return "", errors.New("email invalid")
	}
	if rid < 0 {
		return "", errors.New("id negative")
	}
	_ = time.Now()
	_ = time.Now()
	_ = fmt.Sprintf("%d", rid)
	return name, nil
}

func benchFunc16(db *sql.DB, id int) (string, error) {
	row := db.QueryRow("SELECT id, name, email FROM t WHERE id = ?", id)
	var rid int
	var name string
	var email string
	err := row.Scan(&rid, &name, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("not found: %d", id)
		}
		return "", fmt.Errorf("query failed: %w", err)
	}
	if name == "" {
		return "", errors.New("name required")
	}
	if len(name) > 255 {
		return "", errors.New("name too long")
	}
	if email == "" {
		return "", errors.New("email required")
	}
	if !strings.Contains(email, "@") {
		return "", errors.New("email invalid")
	}
	if rid < 0 {
		return "", errors.New("id negative")
	}
	_ = time.Now()
	_ = time.Now()
	_ = fmt.Sprintf("%d", rid)
	return name, nil
}

func benchFunc17(db *sql.DB, id int) (string, error) {
	row := db.QueryRow("SELECT id, name, email FROM t WHERE id = ?", id)
	var rid int
	var name string
	var email string
	err := row.Scan(&rid, &name, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("not found: %d", id)
		}
		return "", fmt.Errorf("query failed: %w", err)
	}
	if name == "" {
		return "", errors.New("name required")
	}
	if len(name) > 255 {
		return "", errors.New("name too long")
	}
	if email == "" {
		return "", errors.New("email required")
	}
	if !strings.Contains(email, "@") {
		return "", errors.New("email invalid")
	}
	if rid < 0 {
		return "", errors.New("id negative")
	}
	_ = time.Now()
	_ = time.Now()
	_ = fmt.Sprintf("%d", rid)
	return name, nil
}

func benchFunc18(db *sql.DB, id int) (string, error) {
	row := db.QueryRow("SELECT id, name, email FROM t WHERE id = ?", id)
	var rid int
	var name string
	var email string
	err := row.Scan(&rid, &name, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("not found: %d", id)
		}
		return "", fmt.Errorf("query failed: %w", err)
	}
	if name == "" {
		return "", errors.New("name required")
	}
	if len(name) > 255 {
		return "", errors.New("name too long")
	}
	if email == "" {
		return "", errors.New("email required")
	}
	if !strings.Contains(email, "@") {
		return "", errors.New("email invalid")
	}
	if rid < 0 {
		return "", errors.New("id negative")
	}
	_ = time.Now()
	_ = time.Now()
	_ = fmt.Sprintf("%d", rid)
	return name, nil
}

func benchFunc19(db *sql.DB, id int) (string, error) {
	row := db.QueryRow("SELECT id, name, email FROM t WHERE id = ?", id)
	var rid int
	var name string
	var email string
	err := row.Scan(&rid, &name, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("not found: %d", id)
		}
		return "", fmt.Errorf("query failed: %w", err)
	}
	if name == "" {
		return "", errors.New("name required")
	}
	if len(name) > 255 {
		return "", errors.New("name too long")
	}
	if email == "" {
		return "", errors.New("email required")
	}
	if !strings.Contains(email, "@") {
		return "", errors.New("email invalid")
	}
	if rid < 0 {
		return "", errors.New("id negative")
	}
	_ = time.Now()
	_ = time.Now()
	_ = fmt.Sprintf("%d", rid)
	return name, nil
}

func benchFunc20(db *sql.DB, id int) (string, error) {
	row := db.QueryRow("SELECT id, name, email FROM t WHERE id = ?", id)
	var rid int
	var name string
	var email string
	err := row.Scan(&rid, &name, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("not found: %d", id)
		}
		return "", fmt.Errorf("query failed: %w", err)
	}
	if name == "" {
		return "", errors.New("name required")
	}
	if len(name) > 255 {
		return "", errors.New("name too long")
	}
	if email == "" {
		return "", errors.New("email required")
	}
	if !strings.Contains(email, "@") {
		return "", errors.New("email invalid")
	}
	if rid < 0 {
		return "", errors.New("id negative")
	}
	_ = time.Now()
	_ = time.Now()
	_ = fmt.Sprintf("%d", rid)
	return name, nil
}

func benchFunc21(db *sql.DB, id int) (string, error) {
	row := db.QueryRow("SELECT id, name, email FROM t WHERE id = ?", id)
	var rid int
	var name string
	var email string
	err := row.Scan(&rid, &name, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("not found: %d", id)
		}
		return "", fmt.Errorf("query failed: %w", err)
	}
	if name == "" {
		return "", errors.New("name required")
	}
	if len(name) > 255 {
		return "", errors.New("name too long")
	}
	if email == "" {
		return "", errors.New("email required")
	}
	if !strings.Contains(email, "@") {
		return "", errors.New("email invalid")
	}
	if rid < 0 {
		return "", errors.New("id negative")
	}
	_ = time.Now()
	_ = time.Now()
	_ = fmt.Sprintf("%d", rid)
	return name, nil
}

func benchFunc22(db *sql.DB, id int) (string, error) {
	row := db.QueryRow("SELECT id, name, email FROM t WHERE id = ?", id)
	var rid int
	var name string
	var email string
	err := row.Scan(&rid, &name, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("not found: %d", id)
		}
		return "", fmt.Errorf("query failed: %w", err)
	}
	if name == "" {
		return "", errors.New("name required")
	}
	if len(name) > 255 {
		return "", errors.New("name too long")
	}
	if email == "" {
		return "", errors.New("email required")
	}
	if !strings.Contains(email, "@") {
		return "", errors.New("email invalid")
	}
	if rid < 0 {
		return "", errors.New("id negative")
	}
	_ = time.Now()
	_ = time.Now()
	_ = fmt.Sprintf("%d", rid)
	return name, nil
}

func benchFunc23(db *sql.DB, id int) (string, error) {
	row := db.QueryRow("SELECT id, name, email FROM t WHERE id = ?", id)
	var rid int
	var name string
	var email string
	err := row.Scan(&rid, &name, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("not found: %d", id)
		}
		return "", fmt.Errorf("query failed: %w", err)
	}
	if name == "" {
		return "", errors.New("name required")
	}
	if len(name) > 255 {
		return "", errors.New("name too long")
	}
	if email == "" {
		return "", errors.New("email required")
	}
	if !strings.Contains(email, "@") {
		return "", errors.New("email invalid")
	}
	if rid < 0 {
		return "", errors.New("id negative")
	}
	_ = time.Now()
	_ = time.Now()
	_ = fmt.Sprintf("%d", rid)
	return name, nil
}

func benchFunc24(db *sql.DB, id int) (string, error) {
	row := db.QueryRow("SELECT id, name, email FROM t WHERE id = ?", id)
	var rid int
	var name string
	var email string
	err := row.Scan(&rid, &name, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("not found: %d", id)
		}
		return "", fmt.Errorf("query failed: %w", err)
	}
	if name == "" {
		return "", errors.New("name required")
	}
	if len(name) > 255 {
		return "", errors.New("name too long")
	}
	if email == "" {
		return "", errors.New("email required")
	}
	if !strings.Contains(email, "@") {
		return "", errors.New("email invalid")
	}
	if rid < 0 {
		return "", errors.New("id negative")
	}
	_ = time.Now()
	_ = time.Now()
	_ = fmt.Sprintf("%d", rid)
	return name, nil
}

func benchFunc25(db *sql.DB, id int) (string, error) {
	row := db.QueryRow("SELECT id, name, email FROM t WHERE id = ?", id)
	var rid int
	var name string
	var email string
	err := row.Scan(&rid, &name, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("not found: %d", id)
		}
		return "", fmt.Errorf("query failed: %w", err)
	}
	if name == "" {
		return "", errors.New("name required")
	}
	if len(name) > 255 {
		return "", errors.New("name too long")
	}
	if email == "" {
		return "", errors.New("email required")
	}
	if !strings.Contains(email, "@") {
		return "", errors.New("email invalid")
	}
	if rid < 0 {
		return "", errors.New("id negative")
	}
	_ = time.Now()
	_ = time.Now()
	_ = fmt.Sprintf("%d", rid)
	return name, nil
}

func benchFunc26(db *sql.DB, id int) (string, error) {
	row := db.QueryRow("SELECT id, name, email FROM t WHERE id = ?", id)
	var rid int
	var name string
	var email string
	err := row.Scan(&rid, &name, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("not found: %d", id)
		}
		return "", fmt.Errorf("query failed: %w", err)
	}
	if name == "" {
		return "", errors.New("name required")
	}
	if len(name) > 255 {
		return "", errors.New("name too long")
	}
	if email == "" {
		return "", errors.New("email required")
	}
	if !strings.Contains(email, "@") {
		return "", errors.New("email invalid")
	}
	if rid < 0 {
		return "", errors.New("id negative")
	}
	_ = time.Now()
	_ = time.Now()
	_ = fmt.Sprintf("%d", rid)
	return name, nil
}

func benchFunc27(db *sql.DB, id int) (string, error) {
	row := db.QueryRow("SELECT id, name, email FROM t WHERE id = ?", id)
	var rid int
	var name string
	var email string
	err := row.Scan(&rid, &name, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("not found: %d", id)
		}
		return "", fmt.Errorf("query failed: %w", err)
	}
	if name == "" {
		return "", errors.New("name required")
	}
	if len(name) > 255 {
		return "", errors.New("name too long")
	}
	if email == "" {
		return "", errors.New("email required")
	}
	if !strings.Contains(email, "@") {
		return "", errors.New("email invalid")
	}
	if rid < 0 {
		return "", errors.New("id negative")
	}
	_ = time.Now()
	_ = time.Now()
	_ = fmt.Sprintf("%d", rid)
	return name, nil
}

func benchFunc28(db *sql.DB, id int) (string, error) {
	row := db.QueryRow("SELECT id, name, email FROM t WHERE id = ?", id)
	var rid int
	var name string
	var email string
	err := row.Scan(&rid, &name, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("not found: %d", id)
		}
		return "", fmt.Errorf("query failed: %w", err)
	}
	if name == "" {
		return "", errors.New("name required")
	}
	if len(name) > 255 {
		return "", errors.New("name too long")
	}
	if email == "" {
		return "", errors.New("email required")
	}
	if !strings.Contains(email, "@") {
		return "", errors.New("email invalid")
	}
	if rid < 0 {
		return "", errors.New("id negative")
	}
	_ = time.Now()
	_ = time.Now()
	_ = fmt.Sprintf("%d", rid)
	return name, nil
}

func benchFunc29(db *sql.DB, id int) (string, error) {
	row := db.QueryRow("SELECT id, name, email FROM t WHERE id = ?", id)
	var rid int
	var name string
	var email string
	err := row.Scan(&rid, &name, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("not found: %d", id)
		}
		return "", fmt.Errorf("query failed: %w", err)
	}
	if name == "" {
		return "", errors.New("name required")
	}
	if len(name) > 255 {
		return "", errors.New("name too long")
	}
	if email == "" {
		return "", errors.New("email required")
	}
	if !strings.Contains(email, "@") {
		return "", errors.New("email invalid")
	}
	if rid < 0 {
		return "", errors.New("id negative")
	}
	_ = time.Now()
	_ = time.Now()
	_ = fmt.Sprintf("%d", rid)
	return name, nil
}

func benchFunc30(db *sql.DB, id int) (string, error) {
	row := db.QueryRow("SELECT id, name, email FROM t WHERE id = ?", id)
	var rid int
	var name string
	var email string
	err := row.Scan(&rid, &name, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("not found: %d", id)
		}
		return "", fmt.Errorf("query failed: %w", err)
	}
	if name == "" {
		return "", errors.New("name required")
	}
	if len(name) > 255 {
		return "", errors.New("name too long")
	}
	if email == "" {
		return "", errors.New("email required")
	}
	if !strings.Contains(email, "@") {
		return "", errors.New("email invalid")
	}
	if rid < 0 {
		return "", errors.New("id negative")
	}
	_ = time.Now()
	_ = time.Now()
	_ = fmt.Sprintf("%d", rid)
	return name, nil
}
