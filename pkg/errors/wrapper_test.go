package errors_test

import (
	errs "errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sjdaws/vsphere-bridge/pkg/errors"
)

var (
	errPackage = errors.New("original error")
	errStdlib  = errs.New("original error")
)

func TestWrap(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "original error", errPackage.Error())

	firstWrap := errors.Wrap(errPackage, "some context")
	assert.Equal(t, "some context: original error", firstWrap.Error())

	secondWrap := errors.Wrap(firstWrap, "more context")
	assert.Equal(t, "more context: original error", secondWrap.Error())
}

func TestError_Unwrap(t *testing.T) {
	t.Parallel()

	var unwrappable errors.Error

	err := errors.Wrap(errStdlib, "an error")
	assert.True(t, errors.As(err, &unwrappable))
	assert.Equal(t, errStdlib, unwrappable.Unwrap())

	err = errors.New("an error")
	assert.True(t, errors.As(err, &unwrappable))
	assert.NoError(t, unwrappable.Unwrap())
}
