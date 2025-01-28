package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	errWrapperPackage = New("original error")
	errWrapperStdlib  = errors.New("original error")
)

func Test_original(t *testing.T) {
	t.Parallel()

	var unwrappable Error

	wrapped := Wrap(errWrapperPackage, "some context")
	require.ErrorAs(t, wrapped, &unwrappable)
	assert.Equal(t, errWrapperPackage, unwrappable.original())

	wrapped = Wrap(errWrapperStdlib, "some context")
	require.ErrorAs(t, wrapped, &unwrappable)
	assert.Equal(t, errWrapperStdlib, unwrappable.original())

	doubleWrap := Wrap(wrapped, "more context")
	tripleWrap := Wrap(doubleWrap, "even more context")
	require.ErrorAs(t, tripleWrap, &unwrappable)
	assert.Equal(t, errWrapperStdlib, unwrappable.original())
}
