package fixtures

import "fmt"

func rangeLoopCaptureInGoroutine() {
	for _, v := range []int{1, 2, 3} {
		go func() {
			fmt.Println(v) // MATCH /loop variable v captured by closure/
		}()
	}
}

func forLoopCaptureInDefer() {
	for i := 0; i < 3; i++ {
		defer func() {
			fmt.Println(i) // MATCH /loop variable i captured by closure/
		}()
	}
}

func rangeLoopShadowed() {
	for _, v := range []int{1, 2, 3} {
		v := v // Good: create a new variable for each iteration
		go func() {
			fmt.Println(v)
		}()
	}
}

func rangeLoopPassedAsParam() {
	for _, v := range []int{1, 2, 3} {
		go func(val int) {
			fmt.Println(val)
		}(v)
	}
}

func rangeLoopKeyCaptured() {
	for k := range []int{1, 2, 3} {
		go func() {
			fmt.Println(k) // MATCH /loop variable k captured by closure/
		}()
	}
}

func forLoopDeferCapture() {
	for i := 0; i < 10; i++ {
		defer func() {
			_ = i // MATCH /loop variable i captured by closure/
		}()
	}
}

func noClosureInLoop() {
	for _, v := range []int{1, 2, 3} {
		fmt.Println(v)
	}
}

func closureWithoutLoopVar() {
	x := 42
	for range []int{1, 2, 3} {
		go func() {
			fmt.Println(x) // x is not a loop variable
		}()
	}
}
