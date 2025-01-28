package errors_test

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sjdaws/vsphere-bridge/pkg/errors"
)

var errErrorPackage = errors.New("original error")

func TestAs(t *testing.T) {
	t.Parallel()

	var errType errors.Error

	firstWrap := errors.Wrap(errErrorPackage, "some context")
	require.True(t, errors.As(firstWrap, &errType))

	secondWrap := errors.Wrap(firstWrap, "more context")
	require.True(t, errors.As(secondWrap, &errType))
}

func TestIs(t *testing.T) {
	t.Parallel()

	firstWrap := errors.Wrap(errErrorPackage, "some context")
	assert.True(t, errors.Is(firstWrap, errErrorPackage))

	secondWrap := errors.Wrap(firstWrap, "more context")
	assert.True(t, errors.Is(secondWrap, errErrorPackage))
}

func TestNew(t *testing.T) {
	t.Parallel()

	err := errors.New("an error")
	assert.Equal(t, "an error", err.Error())
}

func TestError_Error(t *testing.T) {
	t.Parallel()

	err := errors.New("%s occurred", "an error")
	assert.Equal(t, "an error occurred", err.Error())
}

func TestError_Trace(t *testing.T) {
	t.Parallel()

	_, file, _, ok := runtime.Caller(0)
	require.True(t, ok)

	var errType errors.Error

	err := errors.New("%s occurred", "an error")
	require.ErrorAs(t, err, &errType)

	assert.Equal(t, fmt.Sprintf("an error occurred\n- %s:60: an error occurred", file), errType.Trace())
}
