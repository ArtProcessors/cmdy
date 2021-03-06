package cmdy

import (
	"errors"
	"fmt"
	"strings"
)

// FIXME: these are Unix codes, but other operating systems use
// different codes.
//
// On macOS/Linux, it looks like Go uses status code 2 for a panic,
// so it's probably a good idea to avoid that. Discussion will be
// on one or both of these threads:
//
//	- https://groups.google.com/forum/#!msg/golang-nuts/u9NgKibJsKI/XxCdDihFDAAJ
// 	- https://github.com/golang/go/issues/24284
// 	- https://manpage.me/index.cgi?apropos=0&q=sysexits&sektion=0&manpath=FreeBSD+12-RELEASE&arch=default&format=html
// 	- http://tldp.org/LDP/abs/html/exitcodes.html
// 	- https://man.openbsd.org/sysexits
//
const (
	ExitSuccess  = 0
	ExitFailure  = 1
	ExitUsage    = 64
	ExitInternal = 255
)

type Error interface {
	Code() int
	error
}

// QuietExit is an error you can return to prevent cmdy.Fatal() from printing an error
// message on exit, but still cause it to call os.Exit() with the status code it
// represents.
type QuietExit int

func (e QuietExit) Code() int     { return int(e) }
func (e QuietExit) Error() string { return fmt.Sprintf("exit code %d", e) }

// ErrWithCode allows you to tag an arbitrary error with a status code which will be used
// by cmdy.Fatal() as the exit code.
func ErrWithCode(code int, err error) error {
	if ee, ok := err.(*exitError); ok {
		ee.code = code
		return ee
	}
	return &exitError{err: err, code: code}
}

// HelpRequest returns an error that will instruct cmdy.Fatal() to print the full command
// help.
//
// You can test for this error type using IsHelpRequest(err).
func HelpRequest() error {
	return &usageError{helpRequest: true}
}

// UsageError wraps an existing error so that cmdy.Fatal() will print the full
// command usage above the error message.
func UsageError(err error) error {
	return &usageError{err: err}
}

// UsageErrorf formats an error message so that cmdy.Fatal() will print the full
// command usage above it.
func UsageErrorf(msg string, args ...interface{}) error {
	return &usageError{err: fmt.Errorf(msg, args...)}
}

// IsHelpRequest returns true if err (or the result of errors.Unwrap) is a help request.
func IsHelpRequest(err error) bool {
	var u *usageError
	if !errors.As(err, &u) {
		return false
	}
	return u.helpRequest
}

// IsUsageError returns true if err (or the result of errors.Unwrap) is a usage error.
// This includes help requests.
func IsUsageError(err error) bool {
	var u *usageError
	return errors.As(err, &u)
}

// ErrCode returns the error code associated with the error if it implements
// cmdy.Error, or ExitInternal if not.
func ErrCode(err error) (code int) {
	if err == nil {
		return ExitSuccess
	}
	e, ok := err.(Error)
	if !ok {
		return ExitInternal
	}
	return e.Code()
}

// FormatError builds the output which should be printed to the console.
//
// If the error is a usage error, the full help string will be assigned
// to msg, and if the usage error wraps another error, the text will be
// included at the end.
//
// If the error contains an 'Errors() []error' method, each individual
// error is printed in a list.
//
// If the error is a QuietExit, msg is empty but code will be set to the
// status code.
//
// Otherwise, msg will contain the result of calling Error().
//
func FormatError(err error) (msg string, code int) {
	if err == nil {
		return "", ExitSuccess
	}

	switch err := err.(type) {
	case QuietExit:
		// If we don't return here, a code of '0' will be interpreted as an
		// ExitFailure. In the case of QuietExit, it's a little bit less
		// natural to assume '0' means we want a non-zero exit status even
		// though we are technically returning an error.
		return "", err.Code()

	case *usageError:
		// usageError.usage is lazily populated in Runner.Run() before it is returned:
		msg = strings.TrimSpace(err.usage)
		code = err.Code()

		if err.err != nil {
			if msg != "" {
				msg += "\n\n"
			}
			msg += "error: " + err.err.Error()
		}
		return msg, code

	case Error:
		return err.Error(), err.Code()

	case errorGroup:
		withCode, _ := err.(interface{ Code() int })
		code = ExitFailure
		if withCode != nil {
			code = withCode.Code()
		}
		for i, e := range err.Errors() {
			if i != 0 {
				msg += "\n"
			}
			msg += "- " + e.Error()
		}
		return msg, code

	default:
		return err.Error(), ExitFailure
	}
}

type exitError struct {
	code int
	err  error
}

func (e *exitError) Code() int     { return e.code }
func (e *exitError) Unwrap() error { return e.err }
func (e *exitError) Error() string { return e.err.Error() }

type usageError struct {
	err         error
	usage       string
	helpRequest bool
}

func (u *usageError) Unwrap() error { return u.err }

func (u *usageError) Code() int {
	if u.helpRequest {
		return 0
	} else {
		return ExitUsage
	}
}

func (u *usageError) Error() string {
	if u.helpRequest {
		return "help requested"
	} else if u.err == nil {
		return "usage error"
	}
	return u.err.Error()
}

type errorGroup interface {
	Errors() []error
}
