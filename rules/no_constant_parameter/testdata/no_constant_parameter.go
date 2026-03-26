package fixtures

import (
	"fmt"
	"strings"
)

// Invalid: every call to "repeat" passes 3 as the "times" argument.
func repeat(s string, times int) string { // MATCH /times always receives 3/
	var result string
	for i := 0; i < times; i++ {
		result += s
	}
	return result
}

func callRepeat() {
	fmt.Println(repeat("ha", 3))
	fmt.Println(repeat("ho", 3))
	fmt.Println(repeat("he", 3))
	fmt.Println(repeat("hee", 3))
}

// Invalid: every caller passes true for "verbose".
func logMessage(msg string, verbose bool) { // MATCH /verbose always receives true/
	if verbose {
		fmt.Println("[VERBOSE]", msg)
	} else {
		fmt.Println(msg)
	}
}

func callLogMessage() {
	logMessage("starting", true)
	logMessage("processing", true)
	logMessage("done", true)
	logMessage("finishing", true)
}

// Invalid: the "separator" parameter is always ",".
func joinStrings(parts []string, separator string) string { // MATCH /separator always receives ","/
	return strings.Join(parts, separator)
}

func callJoinStrings(rows [][]string) string {
	var lines []string
	for _, row := range rows {
		lines = append(lines, joinStrings(row, ","))
	}
	lines = append(lines, joinStrings(nil, ","))
	lines = append(lines, joinStrings(nil, ","))
	lines = append(lines, joinStrings(nil, ","))
	return strings.Join(lines, "\n")
}

// Valid: the "times" parameter receives different values at different call sites.
func repeatValid(s string, times int) string {
	var result string
	for i := 0; i < times; i++ {
		result += s
	}
	return result
}

func callRepeatValid() {
	fmt.Println(repeatValid("ha", 3))
	fmt.Println(repeatValid("ho", 5))
	fmt.Println(repeatValid("he", 1))
}

// Valid: the "verbose" parameter varies across call sites.
func logMessageValid(msg string, verbose bool) {
	if verbose {
		fmt.Println("[VERBOSE]", msg)
	} else {
		fmt.Println(msg)
	}
}

func callLogMessageValid() {
	logMessageValid("starting", true)
	logMessageValid("processing", false)
	logMessageValid("done", true)
}

// Valid: the separator varies depending on the output format.
func joinStringsValid(parts []string, separator string) string {
	return strings.Join(parts, separator)
}

func callJoinStringsValid(rows [][]string, format string) string {
	var sep string
	if format == "csv" {
		sep = ","
	} else {
		sep = "\t"
	}
	var lines []string
	for _, row := range rows {
		lines = append(lines, joinStringsValid(row, sep))
	}
	return strings.Join(lines, "\n")
}

// Valid: exported function is not checked by default.
func ExportedFunc(unused string) string {
	return "hello"
}

func callExportedFunc() {
	fmt.Println(ExportedFunc("always"))
	fmt.Println(ExportedFunc("always"))
}

// Valid: function with no callers in this file.
func uncalledFunc(x int) int {
	return x + 1
}

// Valid: function with no parameters.
func noParams() {
	fmt.Println("hello")
}

// Valid: non-constant argument (variable).
func takesVar(x int) int {
	return x * 2
}

func callTakesVar(a int) {
	fmt.Println(takesVar(a))
	fmt.Println(takesVar(a))
}
