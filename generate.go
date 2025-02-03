package errorutils

import (
	"errors"

	"github.com/urfave/cli/v3"
)

// Adds context to an error. Subtle nil check: returns nil if both error and withInner are nil. Bubbles up inner if err is nil. WithMessage does *not* guarantee a new error report (use NewReport); undefined behavior with multiple withInner options.
func New(err error, o ...Option) *Details {
	var wantsInner int = -1
	for i, opt := range o {
		if compareOptions(opt, WithInner(nil)) {
			wantsInner = i
		}
	}
	//must check if the inner error is nil, otherwise retain it
	val, isDtl := err.(*Details)
	if err == nil || ( isDtl && val == nil) {
		if wantsInner < 0 {
			return nil
		}
		err = New(errors.New("wrapper"), o[wantsInner]).Unwrap()
		if err == nil {
			return nil
		}
		val, isDtl = err.(*Details) //necessary juggling to extract from closure
		if val == nil {
			return nil //avoid unnecessary typed error
		}
	}
	e := &Details{}
	//if error is of type Details apply options and return
	if isDtl {
		e = val
		for _, opt := range o {
			opt(e)
		}
		return e
	}
	if _, ok := err.(cli.ExitCoder); ok {
		o = append(o, WithExitCode(err.(cli.ExitCoder).ExitCode()))
	}
	//test is one of the Options is withMsg
	hasmsg := false
	for _, opt := range o {
		if compareOptions(opt, WithMsg("")) {
			hasmsg = true
		}
	}
	if hasmsg {
		e = &Details{
			inner:    err,
			exitcode: 1,
		}
	} else {
		e = &Details{
			msg:      err.Error(),
			exitcode: 1,
		}
	}
	for _, opt := range o {
		opt(e)
	}
	return e
}

// Succinct syntax for new error report
//
// Recommended for errors that could potentially be bugs, otherwise use fmt.Errorf. LineRef is only used in debug mode
func NewReport(msg string, lineRef string, o ...Option) *Details {
	d := &Details{
		msg:      msg,
		lineRef:  lineRef,
		exitcode: 1, //default exit code
	}
	for _, opt := range o {
		opt(d)
	}
	return d
}
