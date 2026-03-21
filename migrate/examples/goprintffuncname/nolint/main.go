package example

import "log"

func myLog(format string, args ...interface{}) { //nolint:goprintffuncname
	const prefix = "[my] "
	log.Printf(prefix+format, args...)
}
