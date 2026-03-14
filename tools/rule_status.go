package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode"
)

func toSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}

func main() {
	data, err := os.ReadFile("architecture/RULES_INDEX.md")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading RULES_INDEX.md: %v\n", err)
		os.Exit(1)
	}

	re := regexp.MustCompile(`\[(lint/\w+/(\w+))\]`)
	matches := re.FindAllStringSubmatch(string(data), -1)

	if len(matches) == 0 {
		fmt.Println("No rules found in RULES_INDEX.md")
		os.Exit(0)
	}

	var implemented, pending []string

	for _, m := range matches {
		fullPath := m[1]  // e.g. lint/complexity/noHighCyclomaticComplexity
		ruleName := m[2]  // e.g. noHighCyclomaticComplexity
		snaked := toSnakeCase(ruleName) // e.g. no_high_cyclomatic_complexity
		dir := "rules/" + snaked

		info, err := os.Stat(dir)
		if err == nil && info.IsDir() {
			implemented = append(implemented, fmt.Sprintf("  %-50s  ->  %s/", fullPath, dir))
		} else {
			pending = append(pending, fmt.Sprintf("  %s", fullPath))
		}
	}

	total := len(implemented) + len(pending)

	fmt.Printf("Implemented (%d/%d):\n", len(implemented), total)
	fmt.Println(strings.Join(implemented, "\n"))

	fmt.Printf("\nPending (%d/%d):\n", len(pending), total)
	if len(pending) > 0 {
		fmt.Println(strings.Join(pending, "\n"))
	} else {
		fmt.Println("  (none)")
	}
}
