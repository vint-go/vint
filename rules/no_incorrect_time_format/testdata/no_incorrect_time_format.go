package fixtures

import "time"

func badSwappedMonthDay() {
	t := time.Now()
	// Bad: month and day are swapped (should be 01-02, not 02-01)
	s := t.Format("2006-02-01") // MATCH /month and day are swapped in time format string: use "2006-01-02" not "2006-02-01"/
	_ = s
}

func badSwappedMonthDayInParse() {
	// Bad: month and day are swapped in time.Parse
	_, _ = time.Parse("2006-02-01", "2025-03-14") // MATCH /month and day are swapped in time format string: use "2006-01-02" not "2006-02-01"/
}

func badSwappedMonthDaySlash() {
	t := time.Now()
	// Bad: month and day are swapped with slash separator
	s := t.Format("2006/02/01") // MATCH /month and day are swapped in time format string: use "2006/01/02" not "2006/02/01"/
	_ = s
}

func goodCorrectFormat() {
	t := time.Now()
	// Good: correct reference time layout
	s := t.Format("2006-01-02 15:04:05")
	_ = s
}

func goodPredefinedConstant() {
	// Good: using predefined time constants
	s := time.Now().Format(time.RFC3339)
	_ = s
}

func goodCorrectSlashFormat() {
	t := time.Now()
	// Good: correct reference time layout with slashes
	s := t.Format("2006/01/02")
	_ = s
}

func goodCorrectParse() {
	// Good: correct reference time layout in Parse
	_, _ = time.Parse("2006-01-02", "2025-03-14")
}
