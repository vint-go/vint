package fixtures

import "sync"

type Server struct {
	sync.Mutex
}

func (s *Server) Process() {
	s.Mutex.Lock()         // MATCH /embedded field Mutex can be omitted from selector/
	defer s.Mutex.Unlock() // MATCH /embedded field Mutex can be omitted from selector/
}

func (s *Server) ProcessValid() {
	s.Lock()
	defer s.Unlock()
}

type Inner struct {
	Value int
}

type Outer struct {
	Inner
}

func example() {
	o := Outer{}
	_ = o.Inner.Value // MATCH /embedded field Inner can be omitted from selector/
}

func exampleValid() {
	o := Outer{}
	_ = o.Value
}

// Ambiguous case: both embedded fields have a method with the same name.
// The simplified selector would be ambiguous, so no warning.
type A struct{}

func (A) Foo() {}

type B struct{}

func (B) Foo() {}

type Ambiguous struct {
	A
	B
}

func ambiguousExample() {
	a := Ambiguous{}
	a.A.Foo() // no match — simplification would be ambiguous
	a.B.Foo() // no match — simplification would be ambiguous
}

// Non-embedded named field — should not trigger
type Named struct {
	Mu sync.Mutex
}

func namedExample() {
	n := Named{}
	n.Mu.Lock() // no match — Mu is not embedded
}
