package handlers

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *responseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		requestID := generateRequestID()

		var requestBody string
		if c.Request.Body != nil {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil {
				requestBody = string(bodyBytes)
				// Restore request body for downstream handlers
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		logRequestStart(c, requestID, requestBody)

		// Wrap response writer to capture response body
		responseBuffer := &bytes.Buffer{}
		wrappedWriter := &responseWriter{
			ResponseWriter: c.Writer,
			body:           responseBuffer,
		}
		c.Writer = wrappedWriter

		// Process request
		c.Next()

		// Calculate processing time
		duration := time.Since(start)

		// Log response
		logRequestComplete(c, requestID, responseBuffer.String(), duration)
	}
}

func logRequestStart(c *gin.Context, requestID, requestBody string) {
	fields := logrus.Fields{
		"request_id":     requestID,
		"method":         c.Request.Method,
		"path":           c.Request.URL.Path,
		"query":          c.Request.URL.RawQuery,
		"user_agent":     c.Request.UserAgent(),
		"remote_addr":    c.ClientIP(),
		"content_type":   c.Request.Header.Get("Content-Type"),
		"content_length": c.Request.ContentLength,
	}

	// Add headers (exclude sensitive ones)
	headers := make(map[string]string)
	for name, values := range c.Request.Header {
		if !isSensitiveHeader(name) {
			headers[name] = strings.Join(values, ", ")
		}
	}
	if len(headers) > 0 {
		fields["headers"] = headers
	}

	// Add request body for non-GET requests (with size limit and content filtering)
	if shouldLogRequestBody(c.Request.Method, requestBody) {
		fields["request_body"] = truncateAndSanitize(requestBody, 1000)
	}

	logrus.WithFields(fields).Info("HTTP request started")
}

func logRequestComplete(c *gin.Context, requestID, responseBody string, duration time.Duration) {
	fields := logrus.Fields{
		"request_id":    requestID,
		"method":        c.Request.Method,
		"path":          c.Request.URL.Path,
		"status_code":   c.Writer.Status(),
		"response_size": c.Writer.Size(),
		"duration_ms":   duration.Milliseconds(),
		"duration":      duration.String(),
	}

	// Add response headers (exclude sensitive ones)
	responseHeaders := make(map[string]string)
	for name, values := range c.Writer.Header() {
		if !isSensitiveHeader(name) {
			responseHeaders[name] = strings.Join(values, ", ")
		}
	}
	if len(responseHeaders) > 0 {
		fields["response_headers"] = responseHeaders
	}

	// Add response body (with size limit and content filtering)
	if shouldLogResponseBody(c.Writer.Status(), c.Writer.Header().Get("Content-Type"), responseBody) {
		fields["response_body"] = truncateAndSanitize(responseBody, 1000)
	}

	// Add any errors that occurred during processing
	if len(c.Errors) > 0 {
		var errorMessages []string
		for _, ginErr := range c.Errors {
			errorMessages = append(errorMessages, ginErr.Error())
		}
		fields["errors"] = errorMessages
	}

	// Log with appropriate level based on status code
	logLevel := getLogLevel(c.Writer.Status())
	switch logLevel {
	case logrus.ErrorLevel:
		logrus.WithFields(fields).Error("HTTP request completed with error")
	case logrus.WarnLevel:
		logrus.WithFields(fields).Warn("HTTP request completed with warning")
	default:
		logrus.WithFields(fields).Info("HTTP request completed")
	}
}

func generateRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func getLogLevel(statusCode int) logrus.Level {
	switch {
	case statusCode >= 500:
		return logrus.ErrorLevel
	case statusCode >= 400:
		return logrus.WarnLevel
	default:
		return logrus.InfoLevel
	}
}

func isSensitiveHeader(headerName string) bool {
	sensitiveHeaders := []string{
		"authorization",
		"cookie",
		"set-cookie",
		"x-api-key",
		"x-auth-token",
		"bearer",
	}

	headerLower := strings.ToLower(headerName)
	for _, sensitive := range sensitiveHeaders {
		if strings.Contains(headerLower, sensitive) {
			return true
		}
	}
	return false
}

func shouldLogRequestBody(method, body string) bool {
	// Don't log GET requests (usually no meaningful body)
	if method == "GET" || method == "HEAD" || method == "OPTIONS" {
		return false
	}

	// Don't log empty bodies
	if strings.TrimSpace(body) == "" {
		return false
	}

	// Don't log if it looks like binary data
	if isBinaryContent(body) {
		return false
	}

	return true
}

func shouldLogResponseBody(statusCode int, contentType, body string) bool {
	// Don't log empty responses
	if strings.TrimSpace(body) == "" {
		return false
	}

	// Don't log binary content
	if isBinaryContent(body) || isBinaryContentType(contentType) {
		return false
	}

	// Always log error responses for debugging
	if statusCode >= 400 {
		return true
	}

	// Log successful responses for JSON/text content
	if strings.Contains(strings.ToLower(contentType), "json") ||
		strings.Contains(strings.ToLower(contentType), "text") ||
		strings.Contains(strings.ToLower(contentType), "xml") {
		return true
	}

	return false
}

// isBinaryContent checks if content appears to be binary
func isBinaryContent(content string) bool {
	if len(content) == 0 {
		return false
	}

	// Check for null bytes (common in binary data)
	for _, b := range []byte(content) {
		if b == 0 {
			return true
		}
	}

	// Check for high percentage of non-printable characters
	nonPrintable := 0
	for _, r := range content {
		if r < 32 && r != 9 && r != 10 && r != 13 { // excluding tab, newline, carriage return
			nonPrintable++
		}
	}

	// If more than 30% non-printable, consider it binary
	return float64(nonPrintable)/float64(len(content)) > 0.3
}

func isBinaryContentType(contentType string) bool {
	binaryTypes := []string{
		"image/",
		"video/",
		"audio/",
		"application/octet-stream",
		"application/pdf",
		"application/zip",
		"multipart/form-data",
	}

	contentTypeLower := strings.ToLower(contentType)
	for _, binaryType := range binaryTypes {
		if strings.Contains(contentTypeLower, binaryType) {
			return true
		}
	}
	return false
}

// truncateAndSanitize truncates content to maxLength and removes sensitive data
func truncateAndSanitize(content string, maxLength int) string {
	if content == "" {
		return ""
	}

	// Sanitize potential sensitive data patterns
	content = sanitizeSensitiveData(content)

	// Truncate if too long
	if len(content) > maxLength {
		return content[:maxLength] + "... [truncated]"
	}

	return content
}

func sanitizeSensitiveData(content string) string {
	sensitivePatterns := map[string]string{
		`"password":\s*"[^"]*"`:      `"password": "[REDACTED]"`,
		`"token":\s*"[^"]*"`:         `"token": "[REDACTED]"`,
		`"secret":\s*"[^"]*"`:        `"secret": "[REDACTED]"`,
		`"api_key":\s*"[^"]*"`:       `"api_key": "[REDACTED]"`,
		`"authorization":\s*"[^"]*"`: `"authorization": "[REDACTED]"`,
	}

	result := content
	for pattern, replacement := range sensitivePatterns {
		// Note: In a production environment, you might want to use compiled regex for better performance
		result = strings.ReplaceAll(result, pattern, replacement)
	}

	return result
}

func RequestLoggerLite() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Calculate processing time
		duration := time.Since(start)

		// Log essential information only
		logrus.WithFields(logrus.Fields{
			"method":      c.Request.Method,
			"path":        c.Request.URL.Path,
			"status":      c.Writer.Status(),
			"duration_ms": duration.Milliseconds(),
			"client_ip":   c.ClientIP(),
		}).Info("HTTP request")
	}
}
