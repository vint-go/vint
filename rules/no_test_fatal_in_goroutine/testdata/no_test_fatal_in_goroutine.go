package fixtures

import "testing"

// Invalid: t.Fatal called from a goroutine

func TestFatalInGoroutine(t *testing.T) {
	go func() {
		t.Fatal("something failed") // MATCH /t.Fatal must not be called from a non-test goroutine/
	}()
}

func TestFatalfInGoroutine(t *testing.T) {
	done := make(chan bool)
	go func() {
		if err := doWork(); err != nil {
			t.Fatalf("work failed: %v", err) // MATCH /t.Fatalf must not be called from a non-test goroutine/
		}
		done <- true
	}()
	<-done
}

func TestFailNowInGoroutine(t *testing.T) {
	go func() {
		t.FailNow() // MATCH /t.FailNow must not be called from a non-test goroutine/
	}()
}

// Invalid: b.Fatal called from a goroutine in a benchmark

func BenchmarkFatalInGoroutine(b *testing.B) {
	go func() {
		b.Fatal("benchmark failed") // MATCH /b.Fatal must not be called from a non-test goroutine/
	}()
}

// Valid: t.Fatal called from the test goroutine

func TestFatalInTestGoroutine(t *testing.T) {
	if err := doWork(); err != nil {
		t.Fatal(err)
	}
}

// Valid: send error back to the test goroutine

func TestErrorViaChan(t *testing.T) {
	done := make(chan error)
	go func() {
		done <- doWork()
	}()
	if err := <-done; err != nil {
		t.Fatal(err)
	}
}

// Valid: not a test function

func helperFunc(t *testing.T) {
	go func() {
		t.Fatal("this is a helper, not a test function")
	}()
}

// helper stubs to make the file parse
func doWork() error { return nil }
