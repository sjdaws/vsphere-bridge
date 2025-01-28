package logging_test

import (
	errs "errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sjdaws/vsphere-bridge/pkg/errors"
	"github.com/sjdaws/vsphere-bridge/pkg/logging"
)

var (
	errPackage = errors.New("error")
	errStdlib  = errs.New("error")
)

func TestDefault(t *testing.T) {
	t.Parallel()

	logger := logging.Default()
	assert.Implements(t, (*logging.Logger)(nil), logger)
}

func TestNew(t *testing.T) {
	t.Parallel()

	logger, err := logging.New(logging.Info, &stdout{data: []byte{}}, 0)
	require.NoError(t, err)
	assert.Implements(t, (*logging.Logger)(nil), logger)
}

func TestNew_ErrNilWriter(t *testing.T) {
	t.Parallel()

	logger, err := logging.New(logging.Info, nil, 0)

	require.Error(t, err)
	require.EqualError(t, err, "nil writer specified, to suppress logs use io.Discard")
	assert.Nil(t, logger)
}

func TestLog_Debug(t *testing.T) {
	t.Parallel()

	testLog(t, testcase{
		expected: expected{
			logLevel: "[debug] test",
			verbosity: map[logging.Verbosity]string{
				logging.Debug: "[debug] test",
				logging.Error: "[error] test",
				logging.Info:  "[info] test",
				logging.Warn:  "[warn] test",
			},
		},
		function:  logging.Logger.Debug,
		message:   "test",
		verbosity: logging.Debug,
	})
}

func TestLog_Error(t *testing.T) {
	t.Parallel()

	testLog(t, testcase{
		expected: expected{
			logLevel: "[error] test",
			verbosity: map[logging.Verbosity]string{
				logging.Error: "[error] test",
			},
		},
		function:  logging.Logger.Error,
		message:   "test",
		verbosity: logging.Error,
	})
}

func TestLog_Fatal(t *testing.T) {
	t.Parallel()

	testLog(t, testcase{
		expected: expected{
			logLevel: "[fatal] test",
			verbosity: map[logging.Verbosity]string{
				logging.Error: "[error] test",
			},
		},
		function:  logging.Logger.Fatal,
		message:   "test",
		verbosity: logging.Error,
	})
}

func TestLog_Info(t *testing.T) {
	t.Parallel()

	testLog(t, testcase{
		expected: expected{
			logLevel: "[info] test",
			verbosity: map[logging.Verbosity]string{
				logging.Error: "[error] test",
				logging.Info:  "[info] test",
				logging.Warn:  "[warn] test",
			},
		},
		function:  logging.Logger.Info,
		message:   "test",
		verbosity: logging.Info,
	})
}

func TestLog_SetDepth(t *testing.T) {
	t.Parallel()

	logger := logging.Default()
	verboseLogger := logger.SetDepth(100)
	assert.Implements(t, (*logging.Logger)(nil), verboseLogger)
	assert.NotEqual(t, logger, verboseLogger)
}

func TestLog_SetVerbosity(t *testing.T) {
	t.Parallel()

	writer := &stdout{data: []byte{}}
	logger, err := logging.New(logging.Error, writer, 0)
	require.NoError(t, err)

	logger.Debug("test")
	assert.Equal(t, []byte{}, writer.data)

	infoLogger := logger.SetVerbosity(logging.Info)
	assert.Implements(t, (*logging.Logger)(nil), infoLogger)
	assert.NotEqual(t, logger, infoLogger)

	infoLogger.Debug("test")
	assert.Equal(t, []byte{}, writer.data)

	debugLogger := logger.SetVerbosity(logging.Debug)
	assert.Implements(t, (*logging.Logger)(nil), debugLogger)
	assert.NotEqual(t, logger, debugLogger)

	debugLogger.Debug("test")
	assert.Contains(t, string(writer.data), "[debug] test")
}

func TestLog_Warn(t *testing.T) {
	t.Parallel()

	testLog(t, testcase{
		expected: expected{
			logLevel: "[warn] test",
			verbosity: map[logging.Verbosity]string{
				logging.Error: "[error] test",
				logging.Warn:  "[warn] test",
			},
		},
		function:  logging.Logger.Warn,
		message:   "test",
		verbosity: logging.Warn,
	})
}

func TestLog_format(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		expected string
		message  any
	}{
		"bool": {
			expected: "true",
			message:  true,
		},
		"int": {
			expected: "5",
			message:  5,
		},
		"pkg error": {
			expected: "error",
			message:  errPackage,
		},
		"stdlib error": {
			expected: "error",
			message:  errStdlib,
		},
		"string": {
			expected: "error",
			message:  "error",
		},
	}

	for name, test := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			writer := &stdout{data: []byte{}}
			logger, err := logging.New(logging.Info, writer, 0)
			require.NoError(t, err)

			logger.Info(test.message)
			assert.Contains(t, string(writer.data), "[info] "+test.expected)
		})
	}
}

type stdout struct {
	data []byte
}

func (s *stdout) Write(data []byte) (int, error) {
	s.data = data

	if strings.Contains(string(data), "[fatal]") {
		panic("fatal caught")
	}

	return len(data), nil
}

type expected struct {
	logLevel  string
	verbosity map[logging.Verbosity]string
}

type testcase struct {
	expected  expected
	function  func(logger logging.Logger, message any, replacements ...any)
	message   any
	verbosity logging.Verbosity
}

func testLog(t *testing.T, testcase testcase) {
	t.Helper()

	testLogLevel(t, testcase)
	testVerbosity(t, testcase)
}

func testLogLevel(t *testing.T, testcase testcase) {
	t.Helper()

	writer := &stdout{data: []byte{}}
	logger, err := logging.New(testcase.verbosity, writer, 0)
	require.NoError(t, err)

	// catch calls to fatal
	defer func() {
		_ = recover()

		assert.Contains(t, string(writer.data), testcase.expected.logLevel)
	}()

	testcase.function(logger, "test")
	assert.Contains(t, string(writer.data), testcase.expected.logLevel)
}

func testVerbosity(t *testing.T, testcase testcase) {
	t.Helper()

	writer := &stdout{data: []byte{}}
	logger, err := logging.New(testcase.verbosity, writer, 0)
	require.NoError(t, err)

	logger.Debug("test")
	assert.Contains(t, string(writer.data), testcase.expected.verbosity[logging.Debug])

	logger.Info("test")
	assert.Contains(t, string(writer.data), testcase.expected.verbosity[logging.Info])

	logger.Error("test")
	assert.Contains(t, string(writer.data), testcase.expected.verbosity[logging.Error])

	logger.Warn("test")
	assert.Contains(t, string(writer.data), testcase.expected.verbosity[logging.Warn])
}
