// ErrorUtils facilitates error handling and reporting
//
// The package provides a custom error type that can be used to add additional information to an error message. The recommendation is to use the NewReport function to create a new error report. The lineRef is only used in debug mode. The exitcode is only used for terminating errors.
// Return Errors that are not related to bugs should be generated with fmt.Errorf and not with NewReport.
package errorutils

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// Details is a custom error type that can be used to add additional information to an error message
type Details struct {
	lineRef  string //will only print if debug is enabled, use a random string instead of the line number
	msg      string //encouraged to be a single line
	exitcode int    // only set for terminating errors
}

func (e *Details) Error() string {
	m := fmt.Sprintf("error: %s", e.msg)
	if logrus.IsLevelEnabled(logrus.DebugLevel) {
		m += fmt.Sprintf("\nLineRef: %s\nExit Code: %d", e.lineRef, e.exitcode)
	}
	return m
}

type Option func(*Details)

func WithLineRef(lineRef string) Option {
	return func(e *Details) {
		e.lineRef = lineRef
	}
}

func WithExitCode(exitcode int) Option {
	return func(e *Details) {
		e.exitcode = exitcode
	}
}

func WithMsg(msg string) Option {
	return func(e *Details) {
		e.msg = e.msg + "\n\t" + msg
	}
}

type Handler func() error

func (fn Handler) Handle() error {
	return fn()
}
