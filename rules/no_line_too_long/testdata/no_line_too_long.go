package fixtures

import (
	"github.com/some/very/long/package/path/that/would/normally/exceed/the/line/length/limit/but/is/excluded/from/checking"
)

// Short function signature that fits within the line length limit
func processUser(userID string) (string, error) {
	return "", nil
}

//go:generate stringer -type=MyType -output=mytype_string.go -trimprefix=MyType -this-is-a-very-long-directive-that-exceeds-limit

func processUserAccountDetailsAndReturnFormattedOutputWithExtendedValidationAndErrorHandlingForAllCasesX(userID string) (string, error) { // MATCH /line is 193 characters, exceeds limit of 120/
	return "", nil
}

var myMap = map[string]string{"key1": "value1", "key2": "value2", "key3": "value3", "key4": "value4", "key5": "value5", "key6": "value6"} // MATCH /line is 193 characters, exceeds limit of 120/

var tooLongAssignment = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" // MATCH /line is 200 characters, exceeds limit of 120/
