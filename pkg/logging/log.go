package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"

	"github.com/fatih/color"

	"github.com/sjdaws/vsphere-bridge/pkg/errors"
)

// Log implementation to Logger.
type Log struct {
	depth     int
	logger    *log.Logger
	verbosity Verbosity
	writer    io.Writer
}

const (
	// defaultDepth depth to trace back through call stack to identify calling file.
	defaultDepth = 2

	// defaultVerbosity verbosity to log messages at or above.
	defaultVerbosity = Info
)

// ErrNilWriter error when a nil writer is passed to New.
var ErrNilWriter = errors.New("nil writer specified, to suppress logs use io.Discard")

// Default create a new Logger using defaults.
func Default() Logger {
	logger, _ := New(defaultVerbosity, os.Stdout, defaultDepth)

	return logger
}

// New create a new Logger.
func New(verbosity Verbosity, writer io.Writer, depth int) (Logger, error) {
	if writer == nil {
		return nil, errors.New(ErrNilWriter)
	}

	if depth < defaultDepth {
		depth = defaultDepth
	}

	return &Log{
		depth:     depth,
		logger:    log.New(writer, "\r\n", log.LstdFlags),
		verbosity: verbosity,
		writer:    writer,
	}, nil
}

// Debug log a debug message if verbosity is Debug or above.
func (l *Log) Debug(message any, replacements ...any) {
	if l.verbosity >= Debug {
		l.logger.Printf(l.format(color.HiBlueString("[debug]"), message), replacements...)
	}
}

// Error log an error message if verbosity is Error or above.
func (l *Log) Error(message any, replacements ...any) {
	if l.verbosity >= Error {
		l.logger.Printf(l.format(color.HiRedString("[error]"), message), replacements...)
	}
}

// Fatal log a fatal message and exit application.
func (l *Log) Fatal(message any, replacements ...any) {
	l.logger.Printf(l.format(color.RedString("[fatal]"), message), replacements...)
	panic("fatal log received")
}

// Info log an informational message if verbosity is Info or above.
func (l *Log) Info(message any, replacements ...any) {
	if l.verbosity >= Info {
		l.logger.Printf(l.format(color.HiGreenString("[info]"), message), replacements...)
	}
}

// SetDepth chainable function to set depth for a single log.
func (l *Log) SetDepth(depth int) Logger {
	logger, _ := New(l.verbosity, l.writer, depth)

	return logger
}

// SetVerbosity chainable function to set verbosity for a single log.
func (l *Log) SetVerbosity(verbosity Verbosity) Logger {
	logger, _ := New(verbosity, l.writer, l.depth)

	return logger
}

// Warn log a warning message if verbosity is Warn or above.
func (l *Log) Warn(message any, replacements ...any) {
	if l.verbosity >= Warn {
		l.logger.Printf(l.format(color.HiYellowString("[warn]"), message), replacements...)
	}
}

// format message to string and prepend caller information.
func (l *Log) format(prefix string, message any) string {
	var content, format string

	// Get string version of message
	switch messageType := message.(type) {
	case errors.Error:
		content = messageType.Trace()
	case error:
		content = messageType.Error()
	default:
		content = fmt.Sprintf("%v", message)
	}

	// Get caller information
	_, filename, line, ok := runtime.Caller(l.depth)
	if ok {
		format += fmt.Sprintf("%s:%d\n", filename, line)
	}

	// Prefix loglevel if applicable
	if prefix != "" {
		format += prefix + " "
	}

	format += content

	return format
}
