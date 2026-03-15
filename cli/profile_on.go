//go:build profile

package cli

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
)

const profileEnabled = true

func startProfile() func() {
	var closers []func()

	if cpuProf := os.Getenv("VINT_CPU_PROFILE"); cpuProf != "" {
		f, err := os.Create(cpuProf)
		if err != nil {
			fail(fmt.Sprintf("could not create CPU profile: %v", err))
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			f.Close()
			fail(fmt.Sprintf("could not start CPU profile: %v", err))
		}
		closers = append(closers, func() {
			pprof.StopCPUProfile()
			f.Close()
		})
	}

	memProf := os.Getenv("VINT_MEM_PROFILE")

	return func() {
		for i := len(closers) - 1; i >= 0; i-- {
			closers[i]()
		}
		if memProf != "" {
			f, err := os.Create(memProf)
			if err != nil {
				fmt.Fprintf(os.Stderr, "could not create memory profile: %v\n", err)
				return
			}
			runtime.GC()
			if err := pprof.WriteHeapProfile(f); err != nil {
				fmt.Fprintf(os.Stderr, "could not write memory profile: %v\n", err)
			}
			f.Close()
		}
	}
}
