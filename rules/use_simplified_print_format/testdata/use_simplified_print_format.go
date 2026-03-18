package fixtures

import (
	"fmt"
	"os"
)

func simplifiedPrintInvalid() {
	// fmt.Print wrapping fmt.Sprintf
	fmt.Print(fmt.Sprintf("Hello, %s!", "world")) // MATCH /fmt.Print(fmt.Sprintf(...)) can be simplified to fmt.Printf(...)/

	// fmt.Println wrapping fmt.Sprintf
	fmt.Println(fmt.Sprintf("Hello, %s!", "world")) // MATCH /fmt.Println(fmt.Sprintf(...)) can be simplified to fmt.Printf(...)/

	// fmt.Fprint wrapping fmt.Sprintf
	fmt.Fprint(os.Stdout, fmt.Sprintf("Hello, %s!", "world")) // MATCH /fmt.Fprint(fmt.Sprintf(...)) can be simplified to fmt.Fprintf(...)/

	// fmt.Fprintln wrapping fmt.Sprintf
	fmt.Fprintln(os.Stderr, fmt.Sprintf("Hello, %s!", "world")) // MATCH /fmt.Fprintln(fmt.Sprintf(...)) can be simplified to fmt.Fprintf(...)/
}

func simplifiedPrintValid() {
	// Direct use of Printf is fine
	fmt.Printf("Hello, %s!\n", "world")

	// Fprintf is fine
	fmt.Fprintf(os.Stdout, "Hello, %s!\n", "world")

	// Print with a non-Sprintf argument is fine
	fmt.Print("hello")

	// Println with a non-Sprintf argument is fine
	fmt.Println("hello")

	// Fprint with a non-Sprintf argument is fine
	fmt.Fprint(os.Stdout, "hello")

	// Print with multiple arguments is fine
	fmt.Print("a", "b")

	// Println with multiple arguments is fine
	fmt.Println("a", "b")

	// Fprint with multiple non-Sprintf arguments is fine
	fmt.Fprint(os.Stdout, "a", "b")

	// Sprint is not a print function (it returns a string)
	_ = fmt.Sprint(fmt.Sprintf("Hello, %s!", "world"))

	// Sprintf wrapping Sprintf is not what this rule checks
	_ = fmt.Sprintf("%s", fmt.Sprintf("Hello, %s!", "world"))
}
