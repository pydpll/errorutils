package errorutils

// Create instance of custom error type or or add info to an existing one.
func New(err error, o ...Option) error {
	//if error is of type Details apply options and return
	if dtl, ok := err.(*Details); ok {
		for _, opt := range o {
			opt(dtl)
		}
		return dtl
	}
	//else create a new Details and apply options
	e := &Details{
		msg: err.Error(),
	}
	for _, opt := range o {
		opt(e)
	}
	return e
}

// Succinct syntax for new error report
//
// Recommended for errors that could potentially be bugs, otherwise use fmt.Errorf. LineRef is only used in debug mode
func NewReport(msg string, lineRef string) *Details {
	return &Details{
		msg:      msg,
		lineRef:  lineRef,
		exitcode: 1, //default exit code
	}
}
