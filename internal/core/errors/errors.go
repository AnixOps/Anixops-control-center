package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ErrorCoder represents an error with a code
type ErrorCoder interface {
	Code() string
}

// ErrorLevel represents the severity level of an error
type ErrorLevel int

const (
	LevelInfo ErrorLevel = iota
	LevelWarning
	LevelError
	LevelCritical
)

func (l ErrorLevel) String() string {
	switch l {
	case LevelInfo:
		return "info"
	case LevelWarning:
		return "warning"
	case LevelError:
		return "error"
	case LevelCritical:
		return "critical"
	default:
		return "unknown"
	}
}

// AppError represents a structured application error
type AppError struct {
	ErrorCode  string                 `json:"code"`
	Message    string                 `json:"message"`
	Detail     string                 `json:"detail,omitempty"`
	Level      ErrorLevel             `json:"level"`
	HTTPStatus int                    `json:"http_status"`
	Context    map[string]interface{} `json:"context,omitempty"`
	Cause      error                  `json:"-"`
	Retryable  bool                   `json:"retryable"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %s (caused by: %v)", e.ErrorCode, e.Message, e.Detail, e.Cause)
	}
	if e.Detail != "" {
		return fmt.Sprintf("[%s] %s: %s", e.ErrorCode, e.Message, e.Detail)
	}
	return fmt.Sprintf("[%s] %s", e.ErrorCode, e.Message)
}

// Code implements ErrorCoder interface
func (e *AppError) Code() string {
	return e.ErrorCode
}

// Unwrap returns the underlying cause
func (e *AppError) Unwrap() error {
	return e.Cause
}

// ToJSON returns the error as JSON
func (e *AppError) ToJSON() string {
	b, _ := json.Marshal(e)
	return string(b)
}

// WithDetail adds detail to the error
func (e *AppError) WithDetail(detail string) *AppError {
	e.Detail = detail
	return e
}

// WithContext adds context to the error
func (e *AppError) WithContext(key string, value interface{}) *AppError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// WithCause adds a cause to the error
func (e *AppError) WithCause(cause error) *AppError {
	e.Cause = cause
	return e
}

// NewError creates a new application error
func NewError(code, message string, level ErrorLevel, httpStatus int) *AppError {
	return &AppError{
		ErrorCode:  code,
		Message:    message,
		Level:      level,
		HTTPStatus: httpStatus,
	}
}

// Common error codes
const (
	// General errors (E0000-E0099)
	CodeSuccess        = "E0000"
	CodeUnknown        = "E0001"
	CodeInvalidRequest = "E0002"
	CodeNotFound       = "E0003"
	CodeAlreadyExists  = "E0004"
	CodePermissionDenied = "E0005"
	CodeTimeout        = "E0006"
	CodeUnavailable    = "E0007"
	CodeRateLimited    = "E0008"
	CodeInternal       = "E0009"

	// Plugin errors (E0100-E0199)
	CodePluginNotFound      = "E0100"
	CodePluginInitFailed    = "E0101"
	CodePluginStartFailed   = "E0102"
	CodePluginStopFailed    = "E0103"
	CodePluginNotExecutable = "E0104"
	CodePluginTimeout       = "E0105"
	CodePluginDependency    = "E0106"
	CodePluginConfig        = "E0107"
	CodePluginHealth        = "E0108"

	// Auth errors (E0200-E0299)
	CodeAuthFailed      = "E0200"
	CodeTokenExpired    = "E0201"
	CodeTokenInvalid    = "E0202"
	CodeUnauthorized    = "E0203"
	CodeForbidden       = "E0204"
	CodeOAuthFailed     = "E0205"
	CodeLDAPFailed      = "E0206"
	CodeSAMLFailed      = "E0207"

	// Database errors (E0300-E0399)
	CodeDatabaseError   = "E0300"
	CodeRecordNotFound  = "E0301"
	CodeDuplicateEntry  = "E0302"
	CodeTransactionFail = "E0303"
	CodeMigrationError  = "E0304"

	// Node errors (E0400-E0499)
	CodeNodeNotFound    = "E0400"
	CodeNodeOffline     = "E0401"
	CodeNodeError       = "E0402"
	CodeNodeTimeout     = "E0403"

	// Task errors (E0500-E0599)
	CodeTaskNotFound    = "E0500"
	CodeTaskFailed      = "E0501"
	CodeTaskTimeout     = "E0502"
	CodeTaskCancelled   = "E0503"

	// Config errors (E0600-E0699)
	CodeConfigInvalid   = "E0600"
	CodeConfigNotFound  = "E0601"
	CodeConfigParse     = "E0602"
)

// Predefined errors
var (
	// General
	ErrUnknown = NewError(CodeUnknown, "Unknown error", LevelError, http.StatusInternalServerError).
			WithDetail("An unexpected error occurred")

	ErrInvalidRequest = NewError(CodeInvalidRequest, "Invalid request", LevelWarning, http.StatusBadRequest)

	ErrNotFound = NewError(CodeNotFound, "Resource not found", LevelWarning, http.StatusNotFound)

	ErrAlreadyExists = NewError(CodeAlreadyExists, "Resource already exists", LevelWarning, http.StatusConflict)

	ErrPermissionDenied = NewError(CodePermissionDenied, "Permission denied", LevelWarning, http.StatusForbidden)

	ErrTimeout = NewError(CodeTimeout, "Operation timed out", LevelWarning, http.StatusGatewayTimeout).
			WithRetryable(true)

	ErrUnavailable = NewError(CodeUnavailable, "Service unavailable", LevelError, http.StatusServiceUnavailable).
			WithRetryable(true)

	ErrRateLimited = NewError(CodeRateLimited, "Rate limit exceeded", LevelWarning, http.StatusTooManyRequests).
			WithRetryable(true)

	// Plugin
	ErrPluginNotFound = NewError(CodePluginNotFound, "Plugin not found", LevelWarning, http.StatusNotFound)

	ErrPluginInitFailed = NewError(CodePluginInitFailed, "Plugin initialization failed", LevelError, http.StatusInternalServerError)

	ErrPluginStartFailed = NewError(CodePluginStartFailed, "Plugin start failed", LevelError, http.StatusInternalServerError)

	ErrPluginStopFailed = NewError(CodePluginStopFailed, "Plugin stop failed", LevelError, http.StatusInternalServerError)

	ErrPluginNotExecutable = NewError(CodePluginNotExecutable, "Plugin not executable", LevelWarning, http.StatusBadRequest)

	ErrPluginDependency = NewError(CodePluginDependency, "Plugin dependency error", LevelError, http.StatusInternalServerError)

	ErrPluginConfig = NewError(CodePluginConfig, "Plugin configuration error", LevelError, http.StatusBadRequest)

	// Auth
	ErrAuthFailed = NewError(CodeAuthFailed, "Authentication failed", LevelWarning, http.StatusUnauthorized)

	ErrTokenExpired = NewError(CodeTokenExpired, "Token expired", LevelInfo, http.StatusUnauthorized)

	ErrTokenInvalid = NewError(CodeTokenInvalid, "Invalid token", LevelWarning, http.StatusUnauthorized)

	ErrUnauthorized = NewError(CodeUnauthorized, "Unauthorized", LevelWarning, http.StatusUnauthorized)

	ErrForbidden = NewError(CodeForbidden, "Forbidden", LevelWarning, http.StatusForbidden)
)

// WithRetryable marks an error as retryable
func (e *AppError) WithRetryable(retryable bool) *AppError {
	e.Retryable = retryable
	return e
}

// IsRetryable checks if an error is retryable
func IsRetryable(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Retryable
	}
	return false
}

// GetCode extracts the error code from an error
func GetCode(err error) string {
	if coder, ok := err.(ErrorCoder); ok {
		return coder.Code()
	}
	return CodeUnknown
}

// GetHTTPStatus extracts HTTP status from an error
func GetHTTPStatus(err error) int {
	if appErr, ok := err.(*AppError); ok {
		return appErr.HTTPStatus
	}
	return http.StatusInternalServerError
}

// Wrap wraps an error with context
func Wrap(err error, code, message string) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return &AppError{
			ErrorCode:  code,
			Message:    message,
			Detail:     appErr.Detail,
			Level:      appErr.Level,
			HTTPStatus: appErr.HTTPStatus,
			Context:    appErr.Context,
			Cause:      appErr,
			Retryable:  appErr.Retryable,
		}
	}
	return NewError(code, message, LevelError, http.StatusInternalServerError).
		WithCause(err)
}

// Is checks if an error has a specific code
func Is(err error, code string) bool {
	return GetCode(err) == code
}