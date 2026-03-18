package fixtures

func busyWait() {
	for { // MATCH /empty for loop busy-waits and consumes 100% CPU; use select, channel, or runtime.Gosched()/
	}
}

func busyWaitNoBody() {
	for { // MATCH /empty for loop busy-waits and consumes 100% CPU; use select, channel, or runtime.Gosched()/
	}
}

func validForWithBody() {
	done := make(chan struct{})
	for {
		select {
		case <-done:
			return
		}
	}
}

func validForWithCondition() {
	x := true
	for x {
		x = false
	}
}

func validForWithInitCondPost() {
	for i := 0; i < 10; i++ {
		_ = i
	}
}

func validForRange() {
	s := []int{1, 2, 3}
	for _, v := range s {
		_ = v
	}
}

func validSelect() {
	done := make(chan struct{})
	select {
	case <-done:
	}
}
