package fixtures

import (
	"fmt"
	"strings"
)

// Config is a placeholder type for the examples.
type Config struct {
	Name string
}

func defaultConfig() Config {
	return Config{Name: "default"}
}

func use(_ Config) {}

// Invalid: The error return value is always nil.
func parseConfig(path string) (cfg Config, err error) { // MATCH /result err is always nil/
	data := defaultConfig()
	_ = path
	return data, nil
}

func callParseConfig() {
	cfg, err := parseConfig("/etc/app.conf")
	if err != nil {
		fmt.Println(err)
	}
	use(cfg)
}

// Invalid: The boolean return value is always true.
func validate(name string) (cleaned string, ok bool) { // MATCH /result ok is always true/
	cleaned = strings.TrimSpace(name)
	return cleaned, true
}

func callValidate() {
	s, ok := validate("hello")
	fmt.Println(s, ok)
}

// Invalid: The second return value "count" is always 0.
func process(items []string) (result []string, count int) { // MATCH /result count is always 0/
	for _, item := range items {
		result = append(result, strings.ToUpper(item))
	}
	return result, 0
}

func callProcess() {
	r, c := process([]string{"a", "b"})
	fmt.Println(r, c)
}

// Valid: The error return value varies depending on the execution path.
func parseConfigValid(path string) (cfg Config, err error) {
	if path == "" {
		return Config{}, fmt.Errorf("empty path")
	}
	data := defaultConfig()
	return data, nil
}

func callParseConfigValid() {
	cfg, err := parseConfigValid("/etc/app.conf")
	if err != nil {
		fmt.Println(err)
	}
	use(cfg)
}

// Valid: The boolean return value varies depending on the input.
func validateValid(name string) (cleaned string, ok bool) {
	cleaned = strings.TrimSpace(name)
	if cleaned == "" {
		return "", false
	}
	return cleaned, true
}

func callValidateValid() {
	s, ok := validateValid("hello")
	fmt.Println(s, ok)
}

// Valid: The count return value reflects actual processing.
func processValid(items []string) (result []string, count int) {
	for _, item := range items {
		if item != "" {
			result = append(result, strings.ToUpper(item))
			count++
		}
	}
	return result, count
}

func callProcessValid() {
	r, c := processValid([]string{"a", "b"})
	fmt.Println(r, c)
}

// Valid: Exported function is not checked by default.
func ExportedAlwaysNil() (result string, err error) {
	return "hello", nil
}

// Valid: No named results -- unnamed results are not checked.
func noNamedResults() (string, error) {
	return "hello", nil
}

// Invalid: bare return with constant assignment (err = nil, same as zero value).
func bareReturnNil(path string) (data []byte, err error) { // MATCH /result err is always nil/
	_ = path
	data = []byte("ok")
	err = nil
	return
}

// Invalid: bare return with no assignment — zero value (nil for error).
func bareReturnZeroValue() (err error) { // MATCH /result err is always nil/
	return
}

// Invalid: returning named result by name, always assigned same constant.
func returnByName() (ok bool) { // MATCH /result ok is always true/
	ok = true
	return ok
}

// Invalid: mix of explicit constant return and bare return (both nil).
func mixedNilReturn(x int) (err error) { // MATCH /result err is always nil/
	if x > 0 {
		return nil
	}
	err = nil
	return
}

// Invalid: bare return with zero value for bool (false).
func bareReturnFalse() (ok bool) { // MATCH /result ok is always false/
	return
}

// Invalid: bare return with zero value for int.
func bareReturnZeroInt() (count int) { // MATCH /result count is always 0/
	return
}

// Valid: bare return with non-constant assignment.
func bareReturn(x int) (result int) {
	result = x * 2
	return
}

func callBareReturn() {
	fmt.Println(bareReturn(5))
}

// Valid: non-constant return value (variable) -- although err is always nil,
// the result return value varies, so the err IS constant. But we need to check
// that the non-constant result does NOT trigger.
func nonConstReturn(x int) (int, error) {
	return x, nil
}

func callNonConstReturn() {
	r, e := nonConstReturn(5)
	fmt.Println(r, e)
}

// Valid: bare return with different assignments.
func differentAssignments(x int) (ok bool) {
	if x > 0 {
		ok = true
	} else {
		ok = false
	}
	return
}

// Valid: bare return with conditional non-zero assignment (result could be 42 or 0).
func conditionalNonZero(x int) (result int) {
	if x > 0 {
		result = 42
	}
	return
}

// Valid: named result modified by increment.
func incrementResult() (count int) {
	count++
	return count
}

// Valid: named result captured by closure.
func closureCapture() (err error) {
	go func() {
		err = fmt.Errorf("oops")
	}()
	return
}

// Valid: named result address taken.
func addressTaken() (err error) {
	setErr(&err)
	return
}

func setErr(e *error) {
	*e = fmt.Errorf("set")
}

// Valid: named result from multi-value assignment.
func multiValueAssign() (err error) {
	_, err = fmt.Println("hi")
	return
}
