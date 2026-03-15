package fixtures

import (
	"os"
	"os/signal"
	"syscall"
)

func unbufferedSignalChannelInvalid() {
	// Bad: unbuffered channel may miss signals
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT) // MATCH /unbuffered os.Signal channel passed to signal.Notify; use a buffered channel with at least capacity 1/
}

func unbufferedSignalChannelInlineMake() {
	// Bad: direct unbuffered make in the call
	signal.Notify(make(chan os.Signal), syscall.SIGINT) // MATCH /unbuffered os.Signal channel passed to signal.Notify; use a buffered channel with at least capacity 1/
}

func bufferedSignalChannelValid() {
	// Good: buffered channel with capacity of 1
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
}

func bufferedSignalChannelLarger() {
	// Good: buffered channel with larger capacity
	c := make(chan os.Signal, 10)
	signal.Notify(c, syscall.SIGINT)
}

func noSignalNotifyCall() {
	// Good: no signal.Notify call, just making a channel
	_ = make(chan os.Signal)
}

func notSignalPackage() {
	// Good: not calling signal.Notify
	c := make(chan os.Signal, 1)
	_ = c
}
