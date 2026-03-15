package fixtures

import (
	"context"
	"time"
)

func noLostCancelDiscarded(ctx context.Context) {
	ctx, _ = context.WithCancel(ctx) // MATCH /the cancel function returned by context.With* must be called, not discarded/
	doWork(ctx)
}

func noLostCancelNotDeferred(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second) // MATCH /the cancel function should be deferred immediately to avoid a context leak/
	result, err := doWorkWithResult(ctx)
	if err != nil {
		return err
	}
	cancel()
	_ = result
	return nil
}

func noLostCancelWithDefer(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	doWork(ctx)
}

func noLostCancelWithDeferTimeout(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	result, err := doWorkWithResult(ctx)
	if err != nil {
		return err
	}
	_ = result
	return nil
}

func noLostCancelWithDeferDeadline(ctx context.Context) {
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Minute))
	defer cancel()
	doWork(ctx)
}

func noLostCancelDeadlineDiscarded(ctx context.Context) {
	ctx, _ = context.WithDeadline(ctx, time.Now().Add(time.Minute)) // MATCH /the cancel function returned by context.With* must be called, not discarded/
	doWork(ctx)
}

func noLostCancelTimeoutDiscarded(ctx context.Context) {
	ctx, _ = context.WithTimeout(ctx, time.Second) // MATCH /the cancel function returned by context.With* must be called, not discarded/
	doWork(ctx)
}

func doWork(ctx context.Context)                              {}
func doWorkWithResult(ctx context.Context) (string, error)    { return "", nil }
