package fixtures

import (
	"fmt"
	"io"
)

// Custom interfaces and types for testing
type MyInterface interface {
	DoSomething()
}

type MyStruct struct{}

func (MyStruct) DoSomething() {}

type AnotherInterface interface {
	DoOtherThing()
}

type AnotherStruct struct{}

func (AnotherStruct) DoOtherThing() {}

// Invalid: concrete type after interface that matches it

func unreachableAfterInterface(x interface{}) {
	switch x.(type) {
	case error:
		fmt.Println("error")
	case *fmt.Stringer:
		// not unreachable, *fmt.Stringer is not a real type
	}
}

func unreachableConcreteAfterInterface(x MyInterface) {
	switch x.(type) {
	case MyInterface:
		fmt.Println("interface")
	case MyStruct: // MATCH /unreachable type case: fixtures.MyStruct is already matched by earlier interface case fixtures.MyInterface/
		fmt.Println("struct")
	}
}

func unreachableWithAssign(x io.Reader) {
	switch v := x.(type) {
	case io.Reader:
		_ = v
	case *io.LimitedReader: // MATCH /unreachable type case: *io.LimitedReader is already matched by earlier interface case io.Reader/
		_ = v
	}
}

// Valid: concrete type before interface

func validConcreteBeforeInterface(x MyInterface) {
	switch x.(type) {
	case MyStruct:
		fmt.Println("struct")
	case MyInterface:
		fmt.Println("interface")
	}
}

// Valid: unrelated types

func validUnrelatedTypes(x interface{}) {
	switch x.(type) {
	case int:
		fmt.Println("int")
	case string:
		fmt.Println("string")
	}
}

// Valid: different interfaces

func validDifferentInterfaces(x interface{}) {
	switch x.(type) {
	case MyInterface:
		fmt.Println("my interface")
	case AnotherInterface:
		fmt.Println("another interface")
	}
}

// Valid: no type switch
func validNoTypeSwitch(x int) {
	switch x {
	case 1:
		fmt.Println("one")
	case 2:
		fmt.Println("two")
	}
}
