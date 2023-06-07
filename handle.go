package errorutils

import (
	"errors"
	"os"

	"github.com/sirupsen/logrus"
)

func HandleFailure(err error, handleFn Handler, o ...Option) error {
	var err2 error
	if err != nil {
		logrus.Error(New(err, o...))
		err2 = handleFn.Handle()
		if err2 != nil {
			logrus.Error(New(err2, o...))
		}
	}
	return errors.Join(err, err2)
}

// Failure has been detected, log it
func LogFailures(err error, o ...Option) {
	cErr := &Details{}
	if !errors.As(err, &cErr) {
		logrus.Error(New(err, o...))
	} else {
		logrus.Error(cErr)
	}
}

func WarnOnFail(err error, o ...Option) {
	if err != nil {
		logrus.Warn(New(err, o...))
	}
}

func PanicOnFail(err error, o ...Option) {
	if err != nil {
		logrus.Fatal(New(err, o...))
	}
}

// SafeClose closes a file and enables error handling on defer
// https://wstrm.dev/posts/errors-join-heart-defer/
func SafeClose(file *os.File, origErr *error) {
	*origErr = errors.Join(*origErr, file.Close())
}

// NotifyClose visibilizes errors on defer
func NotifyClose(file *os.File) {
	err := file.Close()
	if err != nil {
		LogFailures(err)
	}
}

// a generic function x for either error or booleans
