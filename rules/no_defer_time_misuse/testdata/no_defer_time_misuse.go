package fixtures

import (
	"log"
	"time"
)

func badDeferTimeSince() {
	start := time.Now()
	defer log.Println(time.Since(start)) // MATCH /time.Since is evaluated immediately in defer, not when the deferred function runs/
}

func badDeferTimeSinceDirect() {
	start := time.Now()
	defer time.Since(start) // MATCH /time.Since is evaluated immediately in defer, not when the deferred function runs/
}

func goodDeferTimeSinceInClosure() {
	start := time.Now()
	defer func() {
		log.Println(time.Since(start))
	}()
}

func goodNoDeferTimeSince() {
	start := time.Now()
	elapsed := time.Since(start)
	log.Println(elapsed)
}
