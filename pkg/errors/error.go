package errors

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

// Error root error type.
type Error struct {
	file     string
	line     int
	message  string
	previous error
}

// As attempts to set an error to target.
func As(err error, target any) bool {
	//goland:noinspection GoErrorsAs
	return errors.As(err, target)
}

// Is unwraps an error to determine if the error type is target.
func Is(err error, target error) bool {
	return errors.Is(err, target)
}

// New creates a new error message.
func New(message any, replacements ...any) error {
	_, file, line, _ := runtime.Caller(1)

	return Error{
		file:     file,
		line:     line,
		message:  fmt.Sprintf(fmt.Sprintf("%v", message), replacements...),
		previous: nil,
	}
}

// Error provides the last context message and original error message.
func (e Error) Error() string {
	original := e.original()

	if errors.Is(original, e) {
		return e.message
	}

	return fmt.Sprintf("%s: %s", e.message, original.Error())
}

// Trace returns an error message with caller information.
func (e Error) Trace() string {
	stack := make([]string, 0)

	err := e
	previous := true

	for previous {
		var line string
		if err.file != "" {
			line = fmt.Sprintf("%s:%d: ", err.file, err.line)
		}

		stack = append(stack, fmt.Sprintf("%s%s", line, err.message))

		previous = errors.As(err.previous, &err)
	}

	return fmt.Sprintf("%s\n- %s", e.Error(), strings.Join(stack, "\n- "))
}
