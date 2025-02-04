// ErrorUtils facilitates error handling and reporting
//
// The package provides a custom error type that can be used to add additional information to an error message. The recommendation is to use the NewReport function to create a new error report. The lineRef is only used in debug mode. The exitcode is only used for terminating errors.
// Return simple errors such as user errors with fmt.Errorf and not with NewReport to avoid verbosity.
// Details type is capable of wrapping and extends urfave/cli functionality to retain information on exit codes.
package errorutils

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/sirupsen/logrus"
)

// Details is a custom error type that can be used to add additional information to an error message
type Details struct {
	lineRef  string //will only print if debug is enabled, use a random string instead of the line number
	msg      string //encouraged to be a single line
	exitcode int    // only set for terminating errors
	altPrint string
	inner    error
}

func (e *Details) Error() string {
	if isNilDetail(e) {
		return "<nil error>"
	}
	m := fmt.Sprintf("error: %s", e.msg)
	if e.inner != nil {
		m += fmt.Sprintf("\t%v", e.inner)
	}
	if logrus.IsLevelEnabled(logrus.DebugLevel) {
		m += fmt.Sprintf("\tLineRef: %s\tExit Code: %d", e.lineRef, e.exitcode)
	}
	return m
}

func (e *Details) Unwrap() error {
	return e.inner
}

func (e *Details) ExitCode() int {
	return e.exitcode

}

func (e *Details) HasAltprint() bool {
	if isNilDetail(e) {
		return false
	}
	return e.altPrint != ""
}

func (e *Details) Striptype() error {
	if isNilDetail(e) {
		return nil
	}
	return errors.New(e.Error())
}

func isNilDetail(e error) bool {
	if e == nil {
		return true
	}
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

func compareOptions(target, template Option) bool {
	return reflect.ValueOf(target).Pointer() == reflect.ValueOf(template).Pointer()
	// vx := reflect.ValueOf(target)
	// vy := reflect.ValueOf(template)
	// px := vx.Pointer()
	// logrus.Debug("px is ", px)
	// py := vy.Pointer()
	// logrus.Debug("py is ", py)
	// return px == py
}

// WithLineRef is an option to add a line reference to the error message. use a random string instead of the line number. Calling this option on an existing Details value will append the new lineref to the existing one with a '_' separator.
//
//go:noinline
func WithLineRef(lineRef string) Option {
	return func(e *Details) {
		e.lineRef = e.lineRef + "_" + lineRef
	}
}

//go:noinline
func WithExitCode(exitcode int) Option {
	return func(e *Details) {
		e.exitcode = exitcode
	}
}

//go:noinline
func WithMsg(msg string) Option {
	return func(e *Details) {
		e.msg = e.msg + "\n\t" + msg
	}
}

//go:noinline
func WithInner(err error) Option {
	return func(e *Details) {
		e.inner = err
	}
}

type Handler func() *Details

func (fn Handler) Handle() *Details {
	return fn()
}
