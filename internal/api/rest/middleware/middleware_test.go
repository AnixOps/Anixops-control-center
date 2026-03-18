package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestCORSMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(CORSMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("CORS header not set correctly")
	}
}

func TestCORSMiddleware_Options(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(CORSMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("OPTIONS", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, w.Code)
	}
}

func TestCORSMiddleware_Headers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(CORSMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Credentials") != "true" {
		t.Error("Allow-Credentials header not set correctly")
	}
	if w.Header().Get("Access-Control-Allow-Methods") != "POST, OPTIONS, GET, PUT, DELETE, PATCH" {
		t.Error("Allow-Methods header not set correctly")
	}
}

func TestRequestID_Existing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RequestID())
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Request-ID", "existing-id")
	router.ServeHTTP(w, req)

	if w.Header().Get("X-Request-ID") != "existing-id" {
		t.Errorf("expected existing-id, got %s", w.Header().Get("X-Request-ID"))
	}
}

func TestRequestID_New(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RequestID())
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	requestID := w.Header().Get("X-Request-ID")
	if requestID == "" {
		t.Error("request ID should be generated")
	}
}

func TestRequestID_ContextSet(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var capturedID string
	router := gin.New()
	router.Use(RequestID())
	router.GET("/test", func(c *gin.Context) {
		id, exists := c.Get("requestID")
		if exists {
			capturedID = id.(string)
		}
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if capturedID == "" {
		t.Error("request ID should be set in context")
	}
}

func TestRequestLogger(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RequestLogger())
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestRequestLogger_DifferentMethods(t *testing.T) {
	gin.SetMode(gin.TestMode)

	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}

	for _, method := range methods {
		router := gin.New()
		router.Use(RequestLogger())
		router.Handle(method, "/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(method, "/test", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("method %s: expected status %d, got %d", method, http.StatusOK, w.Code)
		}
	}
}

func TestRequestLogger_DifferentPaths(t *testing.T) {
	gin.SetMode(gin.TestMode)

	paths := []string{"/test", "/api/v1/users", "/admin/dashboard"}

	for _, path := range paths {
		router := gin.New()
		router.Use(RequestLogger())
		router.GET(path, func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", path, nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("path %s: expected status %d, got %d", path, http.StatusOK, w.Code)
		}
	}
}

func TestCORSMiddleware_AllOrigins(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(CORSMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Test with different origins
	origins := []string{"http://example.com", "http://localhost:3000", "https://app.example.com"}

	for _, origin := range origins {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Origin", origin)
		router.ServeHTTP(w, req)

		if w.Header().Get("Access-Control-Allow-Origin") != "*" {
			t.Errorf("origin %s: expected Allow-Origin to be *, got %s", origin, w.Header().Get("Access-Control-Allow-Origin"))
		}
	}
}

func TestSecurityHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(SecurityHeaders())
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// Check security headers
	tests := []struct {
		header   string
		expected string
	}{
		{"X-Frame-Options", "DENY"},
		{"X-Content-Type-Options", "nosniff"},
		{"X-XSS-Protection", "1; mode=block"},
		{"Referrer-Policy", "strict-origin-when-cross-origin"},
	}

	for _, tt := range tests {
		if got := w.Header().Get(tt.header); got != tt.expected {
			t.Errorf("header %s: expected %s, got %s", tt.header, tt.expected, got)
		}
	}
}

func TestRateLimit_Allowed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RateLimit(100, time.Minute))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestRateLimit_Exceeded(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RateLimit(2, time.Minute)) // Very low limit for testing
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Make 3 requests, third should be rate limited
	for i := 0; i < 3; i++ {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		if i < 2 && w.Code != http.StatusOK {
			t.Errorf("request %d: expected status %d, got %d", i, http.StatusOK, w.Code)
		}
		if i == 2 && w.Code != http.StatusTooManyRequests {
			t.Errorf("request %d: expected status %d, got %d", i, http.StatusTooManyRequests, w.Code)
		}
	}
}

func TestSanitizeString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"normal text", "normal text"},
		{"<script>alert('xss')</script>", "alert('xss')"},
		{"javascript:void(0)", "void(0)"},
		{"onclick=alert('xss')", "alert('xss')"},
		{"  trimmed  ", "trimmed"},
	}

	for _, tt := range tests {
		got := SanitizeString(tt.input)
		if got != tt.expected {
			t.Errorf("SanitizeString(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestValidateEmail(t *testing.T) {
	validEmails := []string{
		"test@example.com",
		"user.name@domain.org",
		"admin@company.co.uk",
	}

	invalidEmails := []string{
		"",
		"invalid",
		"@example.com",
		"test@",
		"test@example",
		"a@" + strings.Repeat("b", 250) + ".com",
	}

	for _, email := range validEmails {
		if !ValidateEmail(email) {
			t.Errorf("ValidateEmail(%q) should be true", email)
		}
	}

	for _, email := range invalidEmails {
		if ValidateEmail(email) {
			t.Errorf("ValidateEmail(%q) should be false", email)
		}
	}
}

func TestValidateURL(t *testing.T) {
	validURLs := []string{
		"http://example.com",
		"https://example.com/path",
		"https://sub.domain.com/path?query=value",
	}

	invalidURLs := []string{
		"",
		"invalid-url",
		"ftp://example.com",
		"javascript:alert('xss')",
		"JAVASCRIPT:alert('xss')",
		strings.Repeat("a", 2049),
	}

	for _, url := range validURLs {
		if !ValidateURL(url) {
			t.Errorf("ValidateURL(%q) should be true", url)
		}
	}

	for _, url := range invalidURLs {
		if ValidateURL(url) {
			t.Errorf("ValidateURL(%q) should be false", url)
		}
	}
}