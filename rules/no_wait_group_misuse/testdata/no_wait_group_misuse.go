package fixtures

import "sync"

func waitGroupAddInsideGoroutine() {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		go func() {
			wg.Add(1) // MATCH /WaitGroup.Add called inside goroutine, causing a race with Wait/
			defer wg.Done()
			doSomeWork()
		}()
	}
	wg.Wait()
}

func waitGroupAddBeforeGoroutine() {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			doSomeWork()
		}()
	}
	wg.Wait()
}

func waitGroupPointerAddInsideGoroutine(wg *sync.WaitGroup) {
	for i := 0; i < 10; i++ {
		go func() {
			wg.Add(1) // MATCH /WaitGroup.Add called inside goroutine, causing a race with Wait/
			defer wg.Done()
			doSomeWork()
		}()
	}
	wg.Wait()
}

func waitGroupPointerAddBeforeGoroutine(wg *sync.WaitGroup) {
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			doSomeWork()
		}()
	}
	wg.Wait()
}

func waitGroupDoneInsideGoroutineIsOk() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		doSomeWork()
	}()
	wg.Wait()
}

func noWaitGroupAtAll() {
	go func() {
		doSomeWork()
	}()
}

func doSomeWork() {}
