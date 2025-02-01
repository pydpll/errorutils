package errorutils

import "github.com/urfave/cli/v3"

// Create instance of custom error type or or add info to an existing one.
func New(err error, o ...Option) *Details {
	if val, _ := err.(*Details); val == nil {
		return nil
	}
	if _, ok := err.(cli.ExitCoder); ok {
		o = append(o, WithExitCode(err.(cli.ExitCoder).ExitCode()))
	}
	e := &Details{}
	//if error is of type Details apply options and return
	if dtl, ok := err.(*Details); ok {
		e = dtl
		return e
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
