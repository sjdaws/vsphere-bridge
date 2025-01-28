package errors

import (
	"errors"
	"fmt"
	"runtime"
)

// Wrap wraps an error with additional context.
func Wrap(err error, context any, replacements ...any) error {
	_, file, line, _ := runtime.Caller(1)

	return Error{
		file:     file,
		line:     line,
		message:  fmt.Sprintf(fmt.Sprintf("%v", context), replacements...),
		previous: err,
	}
}

// Unwrap returns the next error in the error stack.
func (e Error) Unwrap() error {
	return e.previous
}

// original continuously unwraps an error until the original error is found.
func (e Error) original() error {
	next := e.Unwrap()

	if next == nil {
		return e
	}

	var err Error
	if errors.As(next, &err) {
		return err.original()
	}

	return next
}
