package fixtures

// Invalid: empty import declaration block.

import () // MATCH /empty import declaration block should be removed/

// Invalid: empty var, const, type declaration blocks.

var () // MATCH /empty var declaration block should be removed/

const () // MATCH /empty const declaration block should be removed/

type () // MATCH /empty type declaration block should be removed/

// Valid: non-empty declaration blocks are fine.

var (
	x int
)

const (
	a = 1
)

type (
	MyStruct struct{}
)

// Valid: single-line declarations without parens are fine.

var y int

const b = 2

type AnotherStruct struct{}
