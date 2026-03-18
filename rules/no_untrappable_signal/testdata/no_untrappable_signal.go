package fixtures

import (
	"os"
	"os/signal"
	"syscall"
)

func untrappableSignalNotifyKill() {
	c := make(chan os.Signal, 1)
	// SIGKILL cannot be trapped
	signal.Notify(c, syscall.SIGKILL) // MATCH /SIGKILL cannot be trapped by signal.Notify/
}

func untrappableSignalNotifyStop() {
	c := make(chan os.Signal, 1)
	// SIGSTOP cannot be trapped
	signal.Notify(c, syscall.SIGSTOP) // MATCH /SIGSTOP cannot be trapped by signal.Notify/
}

func untrappableSignalIgnoreKill() {
	// SIGKILL cannot be ignored
	signal.Ignore(syscall.SIGKILL) // MATCH /SIGKILL cannot be trapped by signal.Ignore/
}

func untrappableSignalIgnoreStop() {
	// SIGSTOP cannot be ignored
	signal.Ignore(syscall.SIGSTOP) // MATCH /SIGSTOP cannot be trapped by signal.Ignore/
}

func untrappableSignalNotifyMixed() {
	c := make(chan os.Signal, 1)
	// One valid, one invalid
	signal.Notify(c, syscall.SIGTERM, syscall.SIGKILL) // MATCH /SIGKILL cannot be trapped by signal.Notify/
}

func validSignalNotify() {
	c := make(chan os.Signal, 1)
	// SIGTERM can be trapped
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
}

func validSignalIgnore() {
	// SIGINT can be ignored
	signal.Ignore(syscall.SIGINT)
}

func validSignalNotifyNoSignals() {
	c := make(chan os.Signal, 1)
	// Calling Notify with no signals is valid (catches all)
	signal.Notify(c)
}
