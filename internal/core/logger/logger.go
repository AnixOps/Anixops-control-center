package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

// Level represents log level
type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Field represents a log field
type Field struct {
	Key   string
	Value interface{}
}

// F creates a log field
func F(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

// Fields is a collection of log fields
type Fields []Field

// Logger is a structured logger
type Logger struct {
	mu       sync.Mutex
	writer   io.Writer
	level    Level
	fields   Fields
	prefix   string
	json     bool
	callers  int
}

// Options holds logger options
type Options struct {
	Writer  io.Writer
	Level   Level
	Fields  Fields
	Prefix  string
	JSON    bool
	Callers int // number of caller frames to include
}

// New creates a new logger
func New(opts Options) *Logger {
	if opts.Writer == nil {
		opts.Writer = os.Stdout
	}
	return &Logger{
		writer:  opts.Writer,
		level:   opts.Level,
		fields:  opts.Fields,
		prefix:  opts.Prefix,
		json:    opts.JSON,
		callers: opts.Callers,
	}
}

// DefaultLogger is the default logger instance
var DefaultLogger = New(Options{
	Level: InfoLevel,
})

// SetLevel sets the log level
func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// SetWriter sets the output writer
func (l *Logger) SetWriter(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.writer = w
}

// SetJSON enables or disables JSON output
func (l *Logger) SetJSON(enabled bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.json = enabled
}

// WithPrefix returns a new logger with a prefix
func (l *Logger) WithPrefix(prefix string) *Logger {
	return &Logger{
		writer:  l.writer,
		level:   l.level,
		fields:  l.fields,
		prefix:  l.prefix + prefix,
		json:    l.json,
		callers: l.callers,
	}
}

// WithFields returns a new logger with additional fields
func (l *Logger) WithFields(fields ...Field) *Logger {
	newFields := make(Fields, 0, len(l.fields)+len(fields))
	newFields = append(newFields, l.fields...)
	newFields = append(newFields, fields...)
	return &Logger{
		writer:  l.writer,
		level:   l.level,
		fields:  newFields,
		prefix:  l.prefix,
		json:    l.json,
		callers: l.callers,
	}
}

// WithError returns a logger with an error field
func (l *Logger) WithError(err error) *Logger {
	if err == nil {
		return l
	}
	return l.WithFields(F("error", err.Error()))
}

// log writes a log entry
func (l *Logger) log(level Level, msg string, fields ...Field) {
	if level < l.level {
		return
	}

	entry := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339Nano),
		"level":     level.String(),
		"message":   msg,
	}

	if l.prefix != "" {
		entry["prefix"] = l.prefix
	}

	// Add base fields
	for _, f := range l.fields {
		entry[f.Key] = f.Value
	}

	// Add log-specific fields
	for _, f := range fields {
		entry[f.Key] = f.Value
	}

	// Add caller info
	if l.callers > 0 {
		_, file, line, ok := runtime.Caller(2)
		if ok {
			entry["file"] = file
			entry["line"] = line
		}
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.json {
		b, _ := json.Marshal(entry)
		fmt.Fprintln(l.writer, string(b))
	} else {
		// Text format
		output := fmt.Sprintf("%s [%s] %s",
			entry["timestamp"],
			entry["level"],
			msg,
		)
		if l.prefix != "" {
			output = fmt.Sprintf("%s [%s]", output, l.prefix)
		}
		for _, f := range fields {
			output = fmt.Sprintf("%s %s=%v", output, f.Key, f.Value)
		}
		fmt.Fprintln(l.writer, output)
	}

	// Handle fatal
	if level == FatalLevel {
		os.Exit(1)
	}
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, fields ...Field) {
	l.log(DebugLevel, msg, fields...)
}

// Debugf logs a formatted debug message
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.log(DebugLevel, fmt.Sprintf(format, args...))
}

// Info logs an info message
func (l *Logger) Info(msg string, fields ...Field) {
	l.log(InfoLevel, msg, fields...)
}

// Infof logs a formatted info message
func (l *Logger) Infof(format string, args ...interface{}) {
	l.log(InfoLevel, fmt.Sprintf(format, args...))
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, fields ...Field) {
	l.log(WarnLevel, msg, fields...)
}

// Warnf logs a formatted warning message
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.log(WarnLevel, fmt.Sprintf(format, args...))
}

// Error logs an error message
func (l *Logger) Error(msg string, fields ...Field) {
	l.log(ErrorLevel, msg, fields...)
}

// Errorf logs a formatted error message
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.log(ErrorLevel, fmt.Sprintf(format, args...))
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string, fields ...Field) {
	l.log(FatalLevel, msg, fields...)
}

// Fatalf logs a formatted fatal message and exits
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.log(FatalLevel, fmt.Sprintf(format, args...))
}

// Context key for logger
type ctxKey struct{}

// WithContext returns a context with the logger
func WithContext(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, logger)
}

// FromContext retrieves the logger from context
func FromContext(ctx context.Context) *Logger {
	if logger, ok := ctx.Value(ctxKey{}).(*Logger); ok {
		return logger
	}
	return DefaultLogger
}

// Package-level functions using default logger
func Debug(msg string, fields ...Field) { DefaultLogger.Debug(msg, fields...) }
func Debugf(format string, args ...interface{}) {
	DefaultLogger.Debugf(format, args...)
}
func Info(msg string, fields ...Field)  { DefaultLogger.Info(msg, fields...) }
func Infof(format string, args ...interface{}) {
	DefaultLogger.Infof(format, args...)
}
func Warn(msg string, fields ...Field)  { DefaultLogger.Warn(msg, fields...) }
func Warnf(format string, args ...interface{}) {
	DefaultLogger.Warnf(format, args...)
}
func Error(msg string, fields ...Field) { DefaultLogger.Error(msg, fields...) }
func Errorf(format string, args ...interface{}) {
	DefaultLogger.Errorf(format, args...)
}
func Fatal(msg string, fields ...Field) { DefaultLogger.Fatal(msg, fields...) }
func Fatalf(format string, args ...interface{}) {
	DefaultLogger.Fatalf(format, args...)
}

// WithFields returns a new logger with additional fields
func WithFields(fields ...Field) *Logger { return DefaultLogger.WithFields(fields...) }

// WithError returns a logger with an error field
func WithError(err error) *Logger { return DefaultLogger.WithError(err) }

// SetLevel sets the default logger level
func SetLevel(level Level) { DefaultLogger.SetLevel(level) }

// SetJSON enables JSON output for default logger
func SetJSON(enabled bool) { DefaultLogger.SetJSON(enabled) }