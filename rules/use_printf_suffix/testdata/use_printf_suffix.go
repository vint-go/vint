package fixtures

import (
	"fmt"
	"log"
)

// Invalid: Function has printf-like signature but name does not end with 'f'
func myLog(format string, args ...interface{}) { // MATCH /printf-like formatting function 'myLog' should be named 'myLogf'/
	const prefix = "[my] "
	log.Printf(prefix+format, args...)
}

// Invalid: Custom error logging function missing the 'f' suffix
func logError(format string, args ...any) { // MATCH /printf-like formatting function 'logError' should be named 'logErrorf'/
	log.Printf("[ERROR] "+format, args...)
}

// Invalid: Wrapper around fmt.Sprintf-style formatting without 'f' suffix
func customPrint(format string, args ...interface{}) { // MATCH /printf-like formatting function 'customPrint' should be named 'customPrintf'/
	fmt.Printf("[custom] "+format, args...)
}

// Valid: Function name correctly ends with 'f'
func myLogf(format string, args ...interface{}) {
	const prefix = "[my] "
	log.Printf(prefix+format, args...)
}

// Valid: Properly named custom error logging function
func logErrorf(format string, args ...any) {
	log.Printf("[ERROR] "+format, args...)
}

// Valid: Properly named wrapper function
func customPrintf(format string, args ...interface{}) {
	fmt.Printf("[custom] "+format, args...)
}

// Valid: Not a printf-like function (has return value) - not flagged
func formatMessage(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

// Valid: Not a printf-like function (parameter not named 'format') - not flagged
func logMessage(msg string, args ...interface{}) {
	log.Println(msg)
}

// Valid: Not a printf-like function (only one parameter) - not flagged
func printLine(msg string) {
	fmt.Println(msg)
}
