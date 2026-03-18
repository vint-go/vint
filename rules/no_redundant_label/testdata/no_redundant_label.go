package fixtures

func doSomething(i int) {}

var n, m int
var condition bool

// Invalid: label not referenced by any break or continue.
func redundantLabelOnFor() {
	loop: // MATCH /label "loop" is redundant/
	for i := 0; i < n; i++ {
		doSomething(i)
	}
}

// Invalid: label not referenced by any break or continue in range loop.
func redundantLabelOnRange() {
	items := []int{1, 2, 3}
	myLabel: // MATCH /label "myLabel" is redundant/
	for _, v := range items {
		doSomething(v)
	}
}

// Invalid: label not referenced on switch.
func redundantLabelOnSwitch() {
	sw: // MATCH /label "sw" is redundant/
	switch n {
	case 1:
		doSomething(1)
	case 2:
		doSomething(2)
	}
}

// Invalid: label not referenced on select.
func redundantLabelOnSelect() {
	ch := make(chan int, 1)
	ch <- 1
	sel: // MATCH /label "sel" is redundant/
	select {
	case v := <-ch:
		doSomething(v)
	}
}

// Valid: label used by break in nested loop.
func validLabelBreak() {
	outer:
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if condition {
				break outer
			}
		}
	}
}

// Valid: label used by continue in nested loop.
func validLabelContinue() {
	outer:
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if condition {
				continue outer
			}
		}
	}
}

// Valid: label used by break in switch inside loop.
func validLabelBreakFromSwitch() {
	loop:
	for i := 0; i < n; i++ {
		switch i {
		case 5:
			break loop
		}
	}
}

// Valid: no label at all.
func noLabel() {
	for i := 0; i < n; i++ {
		doSomething(i)
	}
}
