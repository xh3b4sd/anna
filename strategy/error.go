package strategy

import (
	"fmt"

	"github.com/juju/errgo"
)

var (
	maskAny = errgo.MaskFunc(errgo.Any)
)

func maskAnyf(err error, f string, v ...interface{}) error {
	if err == nil {
		return nil
	}

	f = fmt.Sprintf("%s: %s", err.Error(), f)
	newErr := errgo.WithCausef(nil, errgo.Cause(err), f, v...)
	newErr.(*errgo.Err).SetLocation(1)

	return newErr
}

var invalidConfigError = errgo.New("invalid config")

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return errgo.Cause(err) == invalidConfigError
}

var invalidInterfaceError = errgo.New("invalid interface")

// IsInvalidInterface asserts invalidInterfaceError.
func IsInvalidInterface(err error) bool {
	return errgo.Cause(err) == invalidInterfaceError
}

var indexOutOfRangeError = errgo.New("index out of range")

// IsIndexOutOfRange asserts indexOutOfRangeError.
func IsIndexOutOfRange(err error) bool {
	return errgo.Cause(err) == indexOutOfRangeError
}

var notRemovableError = errgo.New("not removable")

// IsNotRemovable asserts notRemovableError.
func IsNotRemovable(err error) bool {
	return errgo.Cause(err) == notRemovableError
}

var notSettableError = errgo.New("not settable")

// IsNotSettable asserts notSettableError.
func IsNotSettable(err error) bool {
	return errgo.Cause(err) == notSettableError
}
