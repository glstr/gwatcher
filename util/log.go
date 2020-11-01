package util

import "log"

func Notice(format string, v ...interface{}) {
	log.Printf(format, v...)
}
