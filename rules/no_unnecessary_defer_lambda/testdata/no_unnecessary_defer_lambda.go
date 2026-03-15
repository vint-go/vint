package fixtures

import (
	"fmt"
	"os"
)

type myFile struct {
	f *os.File
}

func (m *myFile) Close() error {
	return m.f.Close()
}

func badSimpleClose() {
	f, _ := os.Open("file.txt")
	defer func() { f.Close() }() // MATCH /unnecessary defer lambda, call the function directly/
}

func badMethodCall() {
	m := &myFile{}
	defer func() { m.Close() }() // MATCH /unnecessary defer lambda, call the function directly/
}

func badFunctionCall() {
	defer func() { fmt.Println("done") }() // MATCH /unnecessary defer lambda, call the function directly/
}

// Valid: direct defer without lambda
func goodDirectDefer() {
	f, _ := os.Open("file.txt")
	defer f.Close()
}

// Valid: lambda body has more than one statement
func goodMultipleStatements() {
	f, _ := os.Open("file.txt")
	defer func() {
		fmt.Println("closing")
		f.Close()
	}()
}

// Valid: lambda body has no statements
func goodEmptyBody() {
	defer func() {}()
}

// Valid: lambda with parameters
func goodLambdaWithParams() {
	defer func(msg string) {
		fmt.Println(msg)
	}("done")
}

// Valid: lambda body is not a call expression (e.g., assignment)
func goodNotACallExpr() {
	x := 0
	defer func() {
		x = 1
	}()
	_ = x
}
