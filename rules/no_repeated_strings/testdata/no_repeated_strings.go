package fixtures

import "fmt"

// Invalid: The string "user_status" appears 3 times and should be a constant.
func GetStatus() string {
	return "user_status" // MATCH /string literal "user_status" appears 3 times, consider extracting it into a named constant/
}

func SetStatus(s string) string {
	if s == "user_status" {
		return "user_status"
	}
	return s
}

// Invalid: The string "error_code_not_found" appears 3 times.
func CheckError1() string {
	return "error_code_not_found" // MATCH /string literal "error_code_not_found" appears 3 times, consider extracting it into a named constant/
}

func CheckError2() string {
	return "error_code_not_found"
}

func CheckError3() string {
	return "error_code_not_found"
}

// Valid: The repeated string is extracted into a constant.
const statusActive = "active_status"

func GetActive1() string {
	return statusActive
}

func GetActive2() string {
	return statusActive
}

func GetActive3() string {
	return statusActive
}

// Valid: Short strings below the minimum length threshold are acceptable.
func example() {
	a := "ok"
	b := "ok"
	c := "ok"
	_ = a
	_ = b
	_ = c
}

// Valid: A string that only appears once does not trigger the rule.
func greet() string {
	return "Hello, welcome to the application!"
}

// Valid: A string that appears only twice is below the default threshold of 3.
func twoTimes1() string {
	return "only_twice"
}

func twoTimes2() string {
	return "only_twice"
}

// Valid: Strings that appear only in function call arguments are ignored
// when ignore-calls is true (the default).
func callsOnly() {
	fmt.Println("call_only_value")
	fmt.Println("call_only_value")
	fmt.Println("call_only_value")
}

// Invalid: String appears in both call arguments and non-call contexts.
func mixedContexts() {
	x := "mixed_context" // MATCH /string literal "mixed_context" appears 3 times, consider extracting it into a named constant/
	fmt.Println("mixed_context")
	_ = x == "mixed_context"
}

// Valid: Strings in const declarations should not be counted as occurrences.
// Even though "postgres" appears 3 times below, all are in const blocks,
// so it should NOT be flagged.
const dbDriver = "postgres"
const dbDriver2 = "postgres"
const dbDriver3 = "postgres"

// Valid: Strings in const blocks should not contribute to occurrence count.
// "config_key" appears once in a const and twice in code — only 2 non-const
// occurrences, which is below the threshold of 3.
const configConst = "config_key"

func useConfigKey1() string {
	return "config_key"
}

func useConfigKey2() string {
	return "config_key"
}

// Valid: Strings in composite literals (struct/slice/map values) are not counted.
type Config struct {
	Driver string
	Host   string
}

func getConfigs() []Config {
	return []Config{
		{Driver: "comp_lit_str", Host: "comp_lit_str"},
		{Driver: "comp_lit_str", Host: "comp_lit_str"},
	}
}

// Valid: Strings appearing only in map literals are not counted.
func getMap() map[string]string {
	return map[string]string{
		"map_lit_key": "map_lit_val",
		"map_lit_key": "map_lit_val",
		"map_lit_key": "map_lit_val",
	}
}
