package errors

import (
	"errors"
	"net/http"
	"testing"
)

func TestNewError(t *testing.T) {
	err := NewError(CodeNotFound, "Resource not found", LevelWarning, http.StatusNotFound)

	if err.ErrorCode != CodeNotFound {
		t.Errorf("Expected code %s, got %s", CodeNotFound, err.ErrorCode)
	}

	if err.Message != "Resource not found" {
		t.Errorf("Expected message 'Resource not found', got '%s'", err.Message)
	}

	if err.Level != LevelWarning {
		t.Errorf("Expected level %v, got %v", LevelWarning, err.Level)
	}

	if err.HTTPStatus != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, err.HTTPStatus)
	}
}

func TestErrorString(t *testing.T) {
	err := NewError(CodeNotFound, "Not found", LevelWarning, 404)

	expected := "[E0003] Not found"
	if err.Error() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, err.Error())
	}
}

func TestErrorWithDetail(t *testing.T) {
	err := NewError(CodeNotFound, "Not found", LevelWarning, 404).
		WithDetail("User with ID 123 not found")

	if err.Detail != "User with ID 123 not found" {
		t.Errorf("Expected detail, got '%s'", err.Detail)
	}
}

func TestErrorWithCause(t *testing.T) {
	cause := errors.New("underlying error")
	err := NewError(CodeInternal, "Internal error", LevelError, 500).
		WithCause(cause)

	if err.Cause != cause {
		t.Error("Cause not set correctly")
	}
}

func TestErrorWithContext(t *testing.T) {
	err := NewError(CodeNotFound, "Not found", LevelWarning, 404).
		WithContext("user_id", 123)

	if err.Context["user_id"] != 123 {
		t.Errorf("Expected context user_id=123, got %v", err.Context["user_id"])
	}
}

func TestErrorUnwrap(t *testing.T) {
	cause := errors.New("underlying error")
	err := NewError(CodeInternal, "Internal error", LevelError, 500).
		WithCause(cause)

	unwrapped := err.Unwrap()
	if unwrapped != cause {
		t.Error("Unwrap did not return the cause")
	}
}

func TestCode(t *testing.T) {
	err := NewError(CodeNotFound, "Not found", LevelWarning, 404)

	if err.Code() != CodeNotFound {
		t.Errorf("Expected code %s, got %s", CodeNotFound, err.Code())
	}
}

func TestToJSON(t *testing.T) {
	err := NewError(CodeNotFound, "Not found", LevelWarning, 404)
	json := err.ToJSON()

	if json == "" {
		t.Error("ToJSON returned empty string")
	}
}

func TestIsRetryable(t *testing.T) {
	err := NewError(CodeTimeout, "Timeout", LevelWarning, 504).WithRetryable(true)

	if !IsRetryable(err) {
		t.Error("Error should be retryable")
	}

	err2 := NewError(CodeNotFound, "Not found", LevelWarning, 404)

	if IsRetryable(err2) {
		t.Error("Error should not be retryable")
	}
}

func TestGetCode(t *testing.T) {
	err := NewError(CodeNotFound, "Not found", LevelWarning, 404)

	if GetCode(err) != CodeNotFound {
		t.Errorf("Expected code %s, got %s", CodeNotFound, GetCode(err))
	}

	// Test with standard error
	stdErr := errors.New("standard error")
	if GetCode(stdErr) != CodeUnknown {
		t.Errorf("Expected code %s for standard error, got %s", CodeUnknown, GetCode(stdErr))
	}
}

func TestGetHTTPStatus(t *testing.T) {
	err := NewError(CodeNotFound, "Not found", LevelWarning, 404)

	if GetHTTPStatus(err) != 404 {
		t.Errorf("Expected status 404, got %d", GetHTTPStatus(err))
	}

	// Test with standard error
	stdErr := errors.New("standard error")
	if GetHTTPStatus(stdErr) != http.StatusInternalServerError {
		t.Errorf("Expected status %d for standard error, got %d", http.StatusInternalServerError, GetHTTPStatus(stdErr))
	}
}

func TestWrap(t *testing.T) {
	cause := NewError(CodeNotFound, "Not found", LevelWarning, 404)
	err := Wrap(cause, CodeInternal, "Internal error")

	if err.ErrorCode != CodeInternal {
		t.Errorf("Expected code %s, got %s", CodeInternal, err.ErrorCode)
	}

	if err.Cause != cause {
		t.Error("Cause not set correctly")
	}
}

func TestIs(t *testing.T) {
	err := NewError(CodeNotFound, "Not found", LevelWarning, 404)

	if !Is(err, CodeNotFound) {
		t.Error("Is should return true for matching code")
	}

	if Is(err, CodeInternal) {
		t.Error("Is should return false for non-matching code")
	}
}

func TestPredefinedErrors(t *testing.T) {
	// Test a few predefined errors
	if ErrNotFound.ErrorCode != CodeNotFound {
		t.Error("ErrNotFound has wrong code")
	}

	if ErrTimeout.Retryable != true {
		t.Error("ErrTimeout should be retryable")
	}

	if ErrAuthFailed.HTTPStatus != http.StatusUnauthorized {
		t.Error("ErrAuthFailed should have 401 status")
	}
}

func TestErrorLevelString(t *testing.T) {
	tests := []struct {
		level    ErrorLevel
		expected string
	}{
		{LevelInfo, "info"},
		{LevelWarning, "warning"},
		{LevelError, "error"},
		{LevelCritical, "critical"},
	}

	for _, test := range tests {
		if test.level.String() != test.expected {
			t.Errorf("Expected '%s', got '%s'", test.expected, test.level.String())
		}
	}
}