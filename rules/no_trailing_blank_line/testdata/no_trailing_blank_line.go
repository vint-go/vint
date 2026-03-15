package fixtures

import "fmt"

// Invalid: blank line at end of function body.
func badFunction(id int) error { // MATCH /unnecessary trailing blank line/
	_ = id
	return nil

}

// Invalid: blank line at end of if block.
func badIf() {
	err := fmt.Errorf("err")
	if err != nil { // MATCH /unnecessary trailing blank line/
		fmt.Println(err)

	}
}

// Invalid: blank line at end of for loop.
func badFor() {
	for i := 0; i < 10; i++ { // MATCH /unnecessary trailing blank line/
		fmt.Println(i)

	}
}

// Invalid: blank line at end of switch block.
func badSwitch(role string) {
	switch role { // MATCH /unnecessary trailing blank line/
	case "admin":
		fmt.Println("admin")
	default:
		fmt.Println("user")

	}
}

// Valid: no blank line at end of function body.
func goodFunction(id int) error {
	_ = id
	return nil
}

// Valid: no blank line at end of if block.
func goodIf() {
	err := fmt.Errorf("err")
	if err != nil {
		fmt.Println(err)
	}
}

// Valid: no blank line at end of for loop.
func goodFor() {
	for i := 0; i < 10; i++ {
		fmt.Println(i)
	}
}

// Valid: no blank line at end of switch block.
func goodSwitch(role string) {
	switch role {
	case "admin":
		fmt.Println("admin")
	default:
		fmt.Println("user")
	}
}

// Valid: empty block bodies.
func goodEmpty() {
	for {
	}
}
