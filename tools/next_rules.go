package main

import (
	"fmt"
	"os"
	"regexp"
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

	count := 0
	for _, m := range matches {
		fullPath := m[1]
		ruleName := m[2]
		snaked := toSnakeCase(ruleName)
		dir := "rules/" + snaked

		info, err := os.Stat(dir)
		if err == nil && info.IsDir() {
			continue
		}

		count++
		fmt.Println(fullPath)
		if count >= 10 {
			break
		}
	}

	if count == 0 {
		fmt.Println("All rules are implemented!")
	}
}
