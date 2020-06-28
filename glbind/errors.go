package main

import "fmt"

var errors = []string{"", "unspecified error"}

// create a new error with a error-code that will be defined both C and Go side
func newError(format string, args ...interface{}) int {
	errors = append(errors, fmt.Sprintf(format, args...))
	return len(errors) - 1
}
