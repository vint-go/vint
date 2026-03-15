package fixtures

import (
	"fmt"
	"net/http"
	"time"
)

// Invalid: Magic number 100 in struct literal field value.
func createOrder() interface{} {
	s := struct {
		Amount int
	}{
		Amount: 100, // MATCH /magic number: 100, in <assign> detected/
	}
	return s
}

// Invalid: Magic number 5 in binary expression in struct field value.
func createClient() *http.Client {
	return &http.Client{
		Timeout: 5 * time.Second, // MATCH /magic number: 5, in <assign> detected/
	}
}

// Invalid: Magic number 12 in unary expression in assignment.
func compute() {
	res := -12 // MATCH /magic number: 12, in <assign> detected/
	fmt.Println(res)
}

// Invalid: Magic number 12 in binary expression with unary on assignment RHS.
func compute2() {
	var x int32
	res := x + -12 // MATCH /magic number: 12, in <assign> detected/
	fmt.Println(res)
}

// Invalid: Magic number 12 in binary expression on assignment RHS.
func compute3() {
	var x int32
	res := 12 + x // MATCH /magic number: 12, in <assign> detected/
	fmt.Println(res)
}

// Valid: Using a named constant in struct literal field.
const defaultAmount = 100

func createOrderGood() interface{} {
	s := struct {
		Amount int
	}{
		Amount: defaultAmount,
	}
	return s
}

// Valid: Using a named constant in struct field with binary expression.
const clientTimeout = 5

func createClientGood() *http.Client {
	return &http.Client{
		Timeout: clientTimeout * time.Second,
	}
}

// Valid: Using a named constant with unary in binary expression.
const offset = 12

func computeGood() {
	var x int32
	res := x + -offset
	fmt.Println(res)
}

// Valid: 0 and 1.0 are excluded by default.
func defaultValues() {
	x := 0   // 0 is excluded by default
	y := 1.0 // 1.0 is excluded by default
	_ = x
	_ = y
}
