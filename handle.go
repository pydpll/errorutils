package errorutils

import (
	"errors"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

// HandleFailure handles an error by logging it and then calling the handler function
// If the handler function returns an error, it is logged as well
func HandleFailure(err error, handleFn Handler, o ...Option) (err2 *Details) {
	if err != nil {
		logrus.Error(New(err, o...))
		err2 = handleFn.Handle()
		LogFailures(New(err2, o...))
	}
	return err2
}

// Failure has been detected, log it
func LogFailures(err error, o ...Option) {
	if err != nil {
		logrus.Error(New(err, o...))
	}
}

// Failure has been detected, log it
// formatted string should contain one '%s' for the error message
func LogFailuresf(err error, format string, o ...Option) {
	if err != nil {
		logrus.Errorf(format, New(err, o...))
	}
}

func WarnOnFail(err error, o ...Option) {
	if err != nil {
		logrus.Warn(New(err, o...))
	}
}

// formatted string should contain one '%s' for the error message
func WarnOnFailf(err error, format string, o ...Option) {
	if err != nil {
		logrus.Warnf(format, New(err, o...))
	}
}

// irrecoverable programming error
func ExitOnFail(err error, o ...Option) {
	if err != nil {
		optErr := New(err, o...)
		std := logrus.StandardLogger()
		std.Log(logrus.FatalLevel, optErr)
		os.Exit(optErr.ExitCode())
	}
}

// depriecated
// legacy panic, replace with ExitOnFail or use panic() instead
func PanicOnFail(err error, o ...Option) {
	if err != nil {
		optErr := New(err, o...)
		std := logrus.StandardLogger()
		std.Log(logrus.PanicLevel, optErr)
		os.Exit(optErr.ExitCode())
	}
}

// SafeClose closes a file and appends any errors to the error that a function is supposed to return
// https://wstrm.dev/posts/errors-join-heart-defer/
func SafeClose(file io.Closer, origErr *error) {
	*origErr = errors.Join(*origErr, file.Close())
}

// NotifyClose visibilizes errors on defer for functions that do not return an error
func NotifyClose(file io.Closer) {
	err := file.Close()
	if err != nil {
		LogFailures(err)
	}
}
