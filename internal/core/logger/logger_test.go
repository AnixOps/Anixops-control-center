package logger

import (
	"bytes"
	"context"
	"testing"
)

func TestNewLogger(t *testing.T) {
	buf := &bytes.Buffer{}
	log := New(Options{
		Writer: buf,
		Level:  InfoLevel,
	})

	if log == nil {
		t.Fatal("Logger is nil")
	}
}

func TestLogLevel(t *testing.T) {
	buf := &bytes.Buffer{}
	log := New(Options{
		Writer: buf,
		Level:  WarnLevel,
	})

	log.Debug("debug message")
	if buf.Len() > 0 {
		t.Error("Debug message should not be logged at WarnLevel")
	}

	buf.Reset()
	log.Warn("warn message")
	if buf.Len() == 0 {
		t.Error("Warn message should be logged")
	}
}

func TestWithFields(t *testing.T) {
	buf := &bytes.Buffer{}
	log := New(Options{
		Writer: buf,
		Level:  InfoLevel,
		JSON:   false,
	})

	log.WithFields(F("key", "value")).Info("test message")

	output := buf.String()
	if output == "" {
		t.Error("Expected output")
	}
}

func TestWithError(t *testing.T) {
	buf := &bytes.Buffer{}
	log := New(Options{
		Writer: buf,
		Level:  InfoLevel,
	})

	// Test with nil error - should not panic
	log.WithError(nil).Info("test")

	output := buf.String()
	if output == "" {
		t.Error("Expected output")
	}

	buf.Reset()

	// Test with actual error
	log.WithError(nil).Error("test error")
	output = buf.String()
	if output == "" {
		t.Error("Expected output for error")
	}
}

func TestWithPrefix(t *testing.T) {
	buf := &bytes.Buffer{}
	log := New(Options{
		Writer: buf,
		Level:  InfoLevel,
		Prefix: "app: ",
	})

	prefixedLog := log.WithPrefix("sub: ")
	prefixedLog.Info("test")

	output := buf.String()
	if output == "" {
		t.Error("Expected output")
	}
}

func TestJSONOutput(t *testing.T) {
	buf := &bytes.Buffer{}
	log := New(Options{
		Writer: buf,
		Level:  InfoLevel,
		JSON:   true,
	})

	log.Info("test message", F("key", "value"))

	output := buf.String()
	if output == "" {
		t.Error("Expected output")
	}

	// Should contain JSON fields
	if !bytes.Contains(buf.Bytes(), []byte(`"level":"INFO"`)) {
		t.Error("Expected JSON level field")
	}
}

func TestLevelString(t *testing.T) {
	tests := []struct {
		level    Level
		expected string
	}{
		{DebugLevel, "DEBUG"},
		{InfoLevel, "INFO"},
		{WarnLevel, "WARN"},
		{ErrorLevel, "ERROR"},
		{FatalLevel, "FATAL"},
	}

	for _, test := range tests {
		if test.level.String() != test.expected {
			t.Errorf("Expected '%s', got '%s'", test.expected, test.level.String())
		}
	}
}

func TestDefaultLogger(t *testing.T) {
	if DefaultLogger == nil {
		t.Fatal("DefaultLogger is nil")
	}
}

func TestPackageFunctions(t *testing.T) {
	// These should not panic
	Info("test")
	Debug("test")
	Warn("test")
	Error("test")
	Infof("test %s", "message")
	Debugf("test %s", "message")
	Warnf("test %s", "message")
	Errorf("test %s", "message")
}

func TestSetLevel(t *testing.T) {
	buf := &bytes.Buffer{}
	log := New(Options{
		Writer: buf,
		Level:  InfoLevel,
	})

	log.SetLevel(DebugLevel)

	log.Debug("debug message")
	if buf.Len() == 0 {
		t.Error("Debug should be logged after level change")
	}
}

func TestSetWriter(t *testing.T) {
	buf1 := &bytes.Buffer{}
	buf2 := &bytes.Buffer{}
	log := New(Options{
		Writer: buf1,
		Level:  InfoLevel,
	})

	log.Info("first")
	if buf1.Len() == 0 {
		t.Error("First buffer should have content")
	}

	log.SetWriter(buf2)
	log.Info("second")

	if buf2.Len() == 0 {
		t.Error("Second buffer should have content")
	}
}

func TestSetJSON(t *testing.T) {
	buf := &bytes.Buffer{}
	log := New(Options{
		Writer: buf,
		Level:  InfoLevel,
		JSON:   false,
	})

	log.Info("text message")
	buf.Reset()

	log.SetJSON(true)
	log.Info("json message")

	output := buf.String()
	if !bytes.Contains([]byte(output), []byte(`"level":"INFO"`)) {
		t.Error("Expected JSON output")
	}
}

func TestFatal(t *testing.T) {
	// Can't easily test Fatal as it calls os.Exit
	// Just verify the function exists
	_ = Fatal
	_ = Fatalf
}

func TestFromContext(t *testing.T) {
	ctx := context.Background()

	// Without logger in context
	log := FromContext(ctx)
	if log != DefaultLogger {
		t.Error("Should return DefaultLogger when not in context")
	}

	// With logger in context
	customLog := New(Options{Level: DebugLevel})
	ctx = WithContext(ctx, customLog)
	log = FromContext(ctx)
	if log != customLog {
		t.Error("Should return custom logger from context")
	}
}

func TestLoggerFatal(t *testing.T) {
	// Note: Can't test Fatal as it calls os.Exit(1)
	// Just verify method exists
	buf := &bytes.Buffer{}
	log := New(Options{
		Writer: buf,
		Level:  InfoLevel,
	})

	// Verify Fatal method exists
	_ = log.Fatal
	_ = log.Fatalf
}