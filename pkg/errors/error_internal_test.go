package errors

import (
	"errors"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	errErrorPackage = New("original error")
	errErrorStdlib  = errors.New("original error")
)

func TestNew(t *testing.T) {
	t.Parallel()

	_, file, _, ok := runtime.Caller(0)
	require.True(t, ok)

	expected := Error{
		file:     file,
		line:     13,
		message:  "original error",
		previous: nil,
	}
	assert.Equal(t, expected, errErrorPackage)
}

func TestError_Error(t *testing.T) {
	t.Parallel()

	err := Error{
		file:     "",
		line:     0,
		message:  "an error",
		previous: nil,
	}
	assert.Equal(t, "an error", err.Error())

	err = Error{
		file:    "",
		line:    0,
		message: "an error",
		previous: Error{
			file:     "",
			line:     0,
			message:  "original message",
			previous: nil,
		},
	}
	assert.Equal(t, "an error: original message", err.Error())

	err = Error{
		file:     "",
		line:     0,
		message:  "an error",
		previous: errErrorStdlib,
	}
	assert.Equal(t, "an error: original error", err.Error())
}

func TestError_Trace(t *testing.T) {
	t.Parallel()

	err := Error{
		file:    "file",
		line:    1,
		message: "an error",
		previous: Error{
			file:    "file",
			line:    2,
			message: "some context",
			previous: Error{
				file:    "",
				line:    0,
				message: "more context",
				previous: Error{
					file:     "file",
					line:     3,
					message:  "even more context",
					previous: errErrorStdlib,
				},
			},
		},
	}
	expected := "an error: original error\n" +
		"- file:1: an error\n" +
		"- file:2: some context\n" +
		"- more context\n" +
		"- file:3: even more context"

	assert.Equal(t, expected, err.Trace())
}
