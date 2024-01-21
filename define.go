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
	altPrint string
}

func (e *Details) Error() string {
	m := fmt.Sprintf("error: %s", e.msg)
	if logrus.IsLevelEnabled(logrus.DebugLevel) {
		m += fmt.Sprintf("\nLineRef: %s\nExit Code: %d", e.lineRef, e.exitcode)
	}
	return m
}

func (e *Details) ExitCode() int {
	return e.exitcode

}

func (e *Details) HasAltprint() bool {
	if e == nil {
		return false
	}
	return e.altPrint != ""
}

func isNilDetail(e error) bool {
	val, _ := e.(*Details)
	return val == nil
}

// Alternative message for nil errors
func WithAltPrint(altM string) Option {
	return func(e *Details) {
		e.altPrint = altM
	}
}

type Option func(*Details)

// WithLineRef is an option to add a line reference to the error message
// use a random string instead of the line number. Calling this option on an existing Details value will append the new lineref to the existing one with a '_' separator.
func WithLineRef(lineRef string) Option {
	return func(e *Details) {
		e.lineRef = e.lineRef + "_" + lineRef
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

type Handler func() *Details

func (fn Handler) Handle() *Details {
	return fn()
}
