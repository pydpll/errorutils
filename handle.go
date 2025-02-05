package errorutils

import (
	"errors"
	"io"
	"os"
)

// HandleFailure handles an error by logging it and then calling the handler function
// If the handler function returns an error, it is logged as well
func HandleFailure(err error, handleFn Handler, o ...Option) (err2 *Details) {
	dtls := New(err, o...)
	if !isNilDetail(dtls) {
		log.Error().Msg(dtls.Error())
		err2 = handleFn.Handle()
		LogFailures(New(err2, o...))
	} else if dtls.HasAltprint() {
		log.Info().Msg(dtls.altPrint)
	}
	return err2
}

// Failure has been detected, log it
func LogFailures(err error, o ...Option) {
	dtls := New(err, o...)
	if !isNilDetail(dtls) {
		log.Error().Msg(dtls.Error())
	} else if dtls.HasAltprint() {
		log.Info().Msg(dtls.altPrint)
	}
}

// Failure has been detected, log it
// formatted string should contain one '%s' for the error message
func LogFailuresf(err error, format string, o ...Option) {
	dtls := New(err, o...)
	if !isNilDetail(dtls) {
		log.Error().Msgf(format, dtls)
	} else if dtls.HasAltprint() {
		log.Info().Msg(dtls.altPrint)
	}
}

func WarnOnFail(err error, o ...Option) {
	dtls := New(err, o...)

	if !isNilDetail(dtls) {
		log.Warn().Msg(dtls.Error())
	} else if dtls.HasAltprint() {
		log.Info().Msg(dtls.altPrint)
	}
}

// formatted string should contain one '%s' for the error message
func WarnOnFailf(err error, format string, o ...Option) {
	dtls := New(err, o...)
	if !isNilDetail(dtls) {
		log.Warn().Msgf(format, dtls)
	} else if dtls.HasAltprint() {
		log.Info().Msg(dtls.altPrint)
	}
}

// irrecoverable programming error
func ExitOnFail(err error, o ...Option) {
	dtls := New(err, o...)
	if !isNilDetail(dtls) {
		log.Fatal().Msg(dtls.Error())
		os.Exit(dtls.ExitCode())
	} else if dtls.HasAltprint() {
		log.Info().Msg(dtls.altPrint)
	}
}

// SafeClose closes a file and appends any errors to the error that a function is supposed to return
// https://wstrm.dev/posts/errors-join-heart-defer/
func SafeClose(file io.Closer, origErr *error) {
	//avoid changing it to the dynamically alocated slice of errors errors.Join returns if possible
	if cerr := file.Close(); cerr != nil {
		*origErr = errors.Join(*origErr, cerr)
	}
}

// NotifyClose visibilizes errors on defer for functions that do not return an error
func NotifyClose(file io.Closer) {
	err := file.Close()
	if err != nil {
		LogFailures(err)
	}
}
