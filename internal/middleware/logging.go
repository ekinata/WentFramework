package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	wentlog "went-framework/internal/logger"
)

// ResponseWriter wrapper to capture response data
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	// Write to both the original response and our buffer
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

// LoggingMiddleware logs all HTTP requests and responses
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Skip logging for health checks and static files to reduce noise
		if shouldSkipLogging(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// Read request body
		var requestBody []byte
		if r.Body != nil {
			requestBody, _ = io.ReadAll(r.Body)
			// Restore the body for the next handler
			r.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Create response writer wrapper
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // Default status code
		}

		// Call the next handler
		next.ServeHTTP(rw, r)

		// Calculate request duration
		duration := time.Since(start)

		// Log the request and response
		logHTTPRequest(r, rw, requestBody, duration)
	})
}

// logHTTPRequest logs the HTTP request and response details
func logHTTPRequest(r *http.Request, rw *responseWriter, requestBody []byte, duration time.Duration) {
	// Prepare request data
	requestData := map[string]interface{}{
		"method":       r.Method,
		"url":          r.URL.String(),
		"path":         r.URL.Path,
		"remote_addr":  getClientIP(r),
		"user_agent":   r.UserAgent(),
		"referer":      r.Referer(),
		"query_params": r.URL.Query(),
		"headers":      filterHeaders(r.Header),
		"content_type": r.Header.Get("Content-Type"),
	}

	// Add request body if present and not too large
	if len(requestBody) > 0 && len(requestBody) < 10240 { // Max 10KB
		if isJSONContent(r.Header.Get("Content-Type")) {
			var jsonBody interface{}
			if err := json.Unmarshal(requestBody, &jsonBody); err == nil {
				requestData["body"] = jsonBody
			} else {
				requestData["body"] = string(requestBody)
			}
		} else {
			requestData["body"] = string(requestBody)
		}
	} else if len(requestBody) > 0 {
		requestData["body_size"] = len(requestBody)
		requestData["body_truncated"] = true
	}

	// Prepare response data
	responseData := map[string]interface{}{
		"status_code":    rw.statusCode,
		"status_text":    http.StatusText(rw.statusCode),
		"headers":        filterHeaders(rw.Header()),
		"content_type":   rw.Header().Get("Content-Type"),
		"content_length": rw.Header().Get("Content-Length"),
	}

	// Add response body if present and not too large
	responseBody := rw.body.Bytes()
	if len(responseBody) > 0 && len(responseBody) < 10240 { // Max 10KB
		if isJSONContent(rw.Header().Get("Content-Type")) {
			var jsonBody interface{}
			if err := json.Unmarshal(responseBody, &jsonBody); err == nil {
				responseData["body"] = jsonBody
			} else {
				responseData["body"] = string(responseBody)
			}
		} else {
			responseData["body"] = string(responseBody)
		}
	} else if len(responseBody) > 0 {
		responseData["body_size"] = len(responseBody)
		responseData["body_truncated"] = true
	}

	// Prepare log data
	logData := map[string]interface{}{
		"type":        "http_request",
		"request":     requestData,
		"response":    responseData,
		"duration_ms": duration.Milliseconds(),
		"duration":    duration.String(),
		"timestamp":   time.Now().UTC(),
	}

	// Determine log level based on status code
	statusCode := rw.statusCode
	switch {
	case statusCode >= 500:
		wentlog.Error("HTTP Request - Server Error", logData)
	case statusCode >= 400:
		wentlog.Warn("HTTP Request - Client Error", logData)
	case statusCode >= 300:
		wentlog.Info("HTTP Request - Redirect", logData)
	default:
		wentlog.Info("HTTP Request", logData)
	}

	// Also log a summary line for quick scanning
	summaryMessage := fmt.Sprintf("%s %s %d %s - %v",
		r.Method,
		r.URL.Path,
		statusCode,
		http.StatusText(statusCode),
		duration,
	)

	if statusCode >= 400 {
		wentlog.Warn(summaryMessage, map[string]interface{}{
			"client_ip":  getClientIP(r),
			"user_agent": r.UserAgent(),
		})
	} else {
		wentlog.Debug(summaryMessage, map[string]interface{}{
			"client_ip": getClientIP(r),
		})
	}
}

// shouldSkipLogging determines if a request should be skipped from logging
func shouldSkipLogging(path string) bool {
	skipPaths := []string{
		"/favicon.ico",
		"/robots.txt",
		"/health",     // Skip health checks
		"/api/health", // Skip API health checks
	}

	// Skip swagger static files but log swagger API calls
	if strings.HasPrefix(path, "/swagger/") && !strings.HasSuffix(path, ".json") {
		return true
	}

	for _, skipPath := range skipPaths {
		if path == skipPath {
			return true
		}
	}

	return false
}

// getClientIP extracts the real client IP from request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (for proxies)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	return r.RemoteAddr
}

// filterHeaders removes sensitive headers from logging
func filterHeaders(headers http.Header) map[string]string {
	filtered := make(map[string]string)
	sensitiveHeaders := map[string]bool{
		"authorization": true,
		"cookie":        true,
		"set-cookie":    true,
		"x-api-key":     true,
		"x-auth-token":  true,
	}

	for key, values := range headers {
		lowerKey := strings.ToLower(key)
		if sensitiveHeaders[lowerKey] {
			filtered[key] = "[REDACTED]"
		} else if len(values) > 0 {
			filtered[key] = values[0] // Take first value to avoid clutter
		}
	}

	return filtered
}

// isJSONContent checks if content type is JSON
func isJSONContent(contentType string) bool {
	return strings.Contains(strings.ToLower(contentType), "application/json")
}

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate a simple request ID (in production, use uuid)
		requestID := generateRequestID()

		// Add to request context (you might want to use context.WithValue in production)
		r.Header.Set("X-Request-ID", requestID)
		w.Header().Set("X-Request-ID", requestID)

		next.ServeHTTP(w, r)
	})
}

// generateRequestID generates a simple request ID
func generateRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// CORSMiddleware handles CORS headers
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Configure properly in production
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// MiddlewareChain combines multiple middleware into a single handler
func MiddlewareChain(middlewares ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(final http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			final = middlewares[i](final)
		}
		return final
	}
}
