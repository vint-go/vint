package fixtures

import (
	"image"
	"image/color"
)

// LocalStruct is a struct defined in the same package.
type LocalStruct struct {
	X int
	Y int
}

// NestedLocal is another same-package struct.
type NestedLocal struct {
	A string
	B int
	C bool
}

func unkeyedStructLiteral() {
	// Unkeyed composite literal - fragile if image.Point changes
	p := image.Point{1, 2} // MATCH /composite literal uses unkeyed fields/
	_ = p
}

func keyedStructLiteral() {
	// Keyed composite literal - robust against struct changes
	p := image.Point{X: 1, Y: 2}
	_ = p
}

func sliceLiteral() {
	// Unkeyed literals for basic types are acceptable
	s := []int{1, 2, 3}
	_ = s
}

func mapLiteral() {
	// Map literals are acceptable
	m := map[string]int{"a": 1, "b": 2}
	_ = m
}

func arrayLiteral() {
	// Array literals are acceptable
	a := [3]int{1, 2, 3}
	_ = a
}

func emptyStructLiteral() {
	// Empty struct literals are acceptable
	p := image.Point{}
	_ = p
}

func unkeyedColorRGBA() {
	c := color.RGBA{255, 0, 0, 255} // MATCH /composite literal uses unkeyed fields/
	_ = c
}

func keyedColorRGBA() {
	c := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	_ = c
}

func localStructUnkeyed() {
	// Same-package struct with unkeyed fields should NOT be flagged
	l := LocalStruct{1, 2}
	_ = l
}

func localStructKeyed() {
	// Same-package struct with keyed fields (always fine)
	l := LocalStruct{X: 1, Y: 2}
	_ = l
}

func nestedLocalUnkeyed() {
	// Another same-package struct with unkeyed fields should NOT be flagged
	n := NestedLocal{"hello", 42, true}
	_ = n
}

func anonymousStructUnkeyed() {
	// Anonymous struct with unkeyed fields should NOT be flagged
	a := struct {
		X int
		Y int
	}{1, 2}
	_ = a
}
