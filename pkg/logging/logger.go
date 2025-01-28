package logging

// Logger interface.
type Logger interface {
	Debug(message any, replacements ...any)
	Error(message any, replacements ...any)
	Fatal(message any, replacements ...any)
	Info(message any, replacements ...any)
	SetDepth(depth int) Logger
	SetVerbosity(verbosity Verbosity) Logger
	Warn(message any, replacements ...any)
}

// Verbosity log verbosity.
type Verbosity int

const (
	// Error only publish Error logs.
	Error Verbosity = iota + 1

	// Warn only publish Warn and Error logs.
	Warn

	// Info only publish Info, Warn and Error logs.
	Info

	// Debug publish all logs.
	Debug
)
