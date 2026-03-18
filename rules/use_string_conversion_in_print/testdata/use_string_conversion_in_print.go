package fixtures

import (
	"fmt"
	"os"
)

func useStringConversionInPrintInvalid() {
	b := []byte("hello")

	// Passing []byte variable to fmt print functions
	fmt.Println(b)                   // MATCH /convert []byte to string before passing to fmt.Println/
	fmt.Print(b)                     // MATCH /convert []byte to string before passing to fmt.Print/
	fmt.Fprint(os.Stdout, b)        // MATCH /convert []byte to string before passing to fmt.Fprint/
	fmt.Fprintln(os.Stdout, b)      // MATCH /convert []byte to string before passing to fmt.Fprintln/
	_ = fmt.Sprint(b)               // MATCH /convert []byte to string before passing to fmt.Sprint/
	_ = fmt.Sprintln(b)             // MATCH /convert []byte to string before passing to fmt.Sprintln/

	// Explicit []byte conversion passed directly
	fmt.Println([]byte("world"))    // MATCH /convert []byte to string before passing to fmt.Println/
}

func useStringConversionInPrintValid() {
	b := []byte("hello")

	// Already converted to string
	fmt.Println(string(b))
	fmt.Print(string(b))
	fmt.Fprint(os.Stdout, string(b))
	fmt.Fprintln(os.Stdout, string(b))
	_ = fmt.Sprint(string(b))
	_ = fmt.Sprintln(string(b))

	// Passing a regular string is fine
	fmt.Println("hello")
	fmt.Print("world")

	// Passing non-[]byte types is fine
	fmt.Println(42)
	fmt.Println(true)
	fmt.Println(3.14)

	// Using format functions is fine (user controls the format)
	fmt.Printf("%s", b)
	fmt.Sprintf("%s", b)
	fmt.Fprintf(os.Stdout, "%s", b)
}
