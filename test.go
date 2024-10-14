package myLinter

import "log"

// TODO: do something	// want "TODO comment has no author"
func myLog(format string, args ...interface{}) {
	const prefix = "[my] "
	log.Printf(prefix+format, args...)
}
