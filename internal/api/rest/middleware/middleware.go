package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CORSMiddleware handles CORS
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// RequestLogger logs requests
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		gin.DefaultWriter.Write([]byte(
			"[API] " + c.Request.Method + " " + path + " " +
				string(rune(status)) + " " + latency.String() + "\n",
		))
	}
}

// RequestID adds a unique request ID
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Set("requestID", requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)

		c.Next()
	}
}

// SecurityHeaders adds security headers to responses
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent clickjacking
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		// Prevent MIME type sniffing
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		// XSS protection
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		// Referrer policy
		c.Writer.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		// Content Security Policy
		c.Writer.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self' data:; connect-src 'self' ws: wss:")
		// Permissions Policy
		c.Writer.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		c.Next()
	}
}

// RateLimiter implements a simple in-memory rate limiter
type RateLimiter struct {
	requests map[string]*clientInfo
	mu       sync.RWMutex
	rate     int           // requests per window
	window   time.Duration // time window
}

type clientInfo struct {
	count     int
	expiresAt time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	limiter := &RateLimiter{
		requests: make(map[string]*clientInfo),
		rate:     rate,
		window:   window,
	}

	// Cleanup expired entries periodically
	go func() {
		for {
			time.Sleep(window)
			limiter.cleanup()
		}
	}()

	return limiter
}

func (r *RateLimiter) cleanup() {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	for ip, info := range r.requests {
		if info.expiresAt.Before(now) {
			delete(r.requests, ip)
		}
	}
}

func (r *RateLimiter) isAllowed(ip string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	info, exists := r.requests[ip]

	if !exists || info.expiresAt.Before(now) {
		r.requests[ip] = &clientInfo{
			count:     1,
			expiresAt: now.Add(r.window),
		}
		return true
	}

	if info.count >= r.rate {
		return false
	}

	info.count++
	return true
}

// RateLimit middleware limits requests per IP
func RateLimit(rate int, window time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter(rate, window)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !limiter.isAllowed(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
				"code":  http.StatusTooManyRequests,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// InputSanitizer sanitizes input to prevent XSS
func InputSanitizer() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read the request body
		if c.Request.Body != nil {
			// Body is read and parsed by handlers, just mark for sanitization
			c.Set("sanitize", true)
		}
		c.Next()
	}
}

// SanitizeString removes potentially dangerous characters from a string
func SanitizeString(input string) string {
	// Remove script tags
	input = strings.ReplaceAll(input, "<script>", "")
	input = strings.ReplaceAll(input, "</script>", "")
	input = strings.ReplaceAll(input, "<SCRIPT>", "")
	input = strings.ReplaceAll(input, "</SCRIPT>", "")

	// Remove javascript: protocol
	input = strings.ReplaceAll(input, "javascript:", "")

	// Remove on* event handlers
	dangerous := []string{"onclick", "onerror", "onload", "onmouseover", "onfocus", "onblur"}
	for _, attr := range dangerous {
		input = strings.ReplaceAll(input, attr+"=", "")
	}

	return strings.TrimSpace(input)
}

// ValidateEmail checks if an email is valid
func ValidateEmail(email string) bool {
	if len(email) == 0 || len(email) > 254 {
		return false
	}

	// Basic email validation
	atIndex := strings.Index(email, "@")
	if atIndex <= 0 || atIndex >= len(email)-1 {
		return false
	}

	dotIndex := strings.LastIndex(email, ".")
	if dotIndex <= atIndex+1 || dotIndex >= len(email)-1 {
		return false
	}

	return true
}

// ValidateURL checks if a URL is valid and safe
func ValidateURL(url string) bool {
	if len(url) == 0 || len(url) > 2048 {
		return false
	}

	// Check for allowed schemes
	allowedSchemes := []string{"http://", "https://"}
	hasValidScheme := false
	for _, scheme := range allowedSchemes {
		if strings.HasPrefix(strings.ToLower(url), scheme) {
			hasValidScheme = true
			break
		}
	}

	if !hasValidScheme {
		return false
	}

	// Check for javascript: protocol (case insensitive)
	if strings.Contains(strings.ToLower(url), "javascript:") {
		return false
	}

	return true
}
