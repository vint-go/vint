//go:build !profile

package cli

const profileEnabled = false

func startProfile() func() {
	return func() {}
}
