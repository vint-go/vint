package fixtures

import (
	"context"
	"net/http"
)

func handlerNoContext(ctx context.Context) {
	// Goroutine does not receive the parent context
	go func() { // MATCH /goroutine does not propagate the parent context/
		result := longRunningOperation()
		processResult(result)
	}()
}

func fetchAllNoContext(ctx context.Context, urls []string) {
	for _, url := range urls {
		// Context not passed to HTTP request
		go func(u string) { // MATCH /goroutine does not propagate the parent context/
			resp, _ := http.Get(u)
			defer resp.Body.Close()
		}(url)
	}
}

func handlerWithContext(ctx context.Context) {
	// Goroutine receives and respects the parent context
	go func(ctx context.Context) {
		select {
		case <-ctx.Done():
			return
		case result := <-doWork(ctx):
			processResult(result)
		}
	}(ctx)
}

func fetchAllWithContext(ctx context.Context, urls []string) {
	for _, url := range urls {
		go func(ctx context.Context, u string) {
			req, _ := http.NewRequestWithContext(ctx, "GET", u, nil)
			resp, _ := http.DefaultClient.Do(req)
			defer resp.Body.Close()
		}(ctx, url)
	}
}

func handlerWithClosure(ctx context.Context) {
	// Goroutine captures context via closure — valid
	go func() {
		select {
		case <-ctx.Done():
			return
		default:
		}
	}()
}

func noContextParam() {
	// Function has no context param — should not trigger
	go func() {
		longRunningOperation()
	}()
}

func handlerGoNamedFunc(ctx context.Context) {
	// go with named function call without context — should trigger
	go doSomethingWithoutContext() // MATCH /goroutine does not propagate the parent context/
}

func handlerGoNamedFuncWithCtx(ctx context.Context) {
	// go with named function call with context — should not trigger
	go doSomethingWithContext(ctx)
}

func longRunningOperation() interface{} { return nil }
func processResult(interface{})          {}
func doWork(ctx context.Context) <-chan interface{} {
	return nil
}
func doSomethingWithoutContext()              {}
func doSomethingWithContext(ctx context.Context) {}
