package fixtures

import "fmt"

// Invalid: blank line at start of function body.
func badFunction(id int) error { // MATCH /unnecessary leading blank line/

	_ = id
	return nil
}

// Invalid: blank line at start of if block.
func badIf() {
	err := fmt.Errorf("err")
	if err != nil { // MATCH /unnecessary leading blank line/

		fmt.Println(err)
	}
}

// Invalid: blank line at start of for loop.
func badFor() {
	for i := 0; i < 10; i++ { // MATCH /unnecessary leading blank line/

		fmt.Println(i)
	}
}

// Invalid: blank line at start of case clause.
func badSwitch(status string) {
	switch status {
	case "active": // MATCH /unnecessary leading blank line/

		fmt.Println("active")
	case "inactive":
		fmt.Println("inactive")
	}
}

// Valid: no blank line at start of function body.
func goodFunction(id int) error {
	_ = id
	return nil
}

// Valid: no blank line at start of if block.
func goodIf() {
	err := fmt.Errorf("err")
	if err != nil {
		fmt.Println(err)
	}
}

// Valid: no blank line at start of for loop.
func goodFor() {
	for i := 0; i < 10; i++ {
		fmt.Println(i)
	}
}

// Valid: no blank line at start of switch case clause.
func goodSwitch(status string) {
	switch status {
	case "active":
		fmt.Println("active")
	case "inactive":
		fmt.Println("inactive")
	}
}

// Valid: empty block bodies.
func goodEmpty() {
	for {
	}
}
