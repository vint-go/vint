package fixtures

import (
	"image"
	"image/color"
)

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
