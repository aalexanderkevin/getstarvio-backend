package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	maxHTTPBodyLogBytes   = 2048
	maxHTTPHeaderLogBytes = 4096
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	_, _ = w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *bodyLogWriter) WriteString(s string) (int, error) {
	_, _ = w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func HTTPLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		rid := c.GetString("request_id")
		if rid == "" {
			rid = c.GetHeader("X-Request-ID")
		}

		routePath := c.FullPath()
		if routePath == "" && c.Request != nil && c.Request.URL != nil {
			routePath = c.Request.URL.Path
		}

		requestBody := readRequestBody(c)
		requestHeaders := headersToLog(c.Request.Header)

		writer := &bodyLogWriter{ResponseWriter: c.Writer}
		c.Writer = writer

		log.Printf("http request rid=%s method=%s path=%s query=%s ip=%s ua=%q headers=%s payload=%q",
			rid,
			c.Request.Method,
			routePath,
			c.Request.URL.RawQuery,
			c.ClientIP(),
			c.Request.UserAgent(),
			requestHeaders,
			truncateForLog(requestBody, maxHTTPBodyLogBytes),
		)

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		respBody := writer.body.String()
		responseHeaders := headersToLog(c.Writer.Header())

		errMsg := ""
		if len(c.Errors) > 0 {
			errMsg = c.Errors.String()
		}

		log.Printf("http response rid=%s method=%s path=%s status=%d latency=%s size=%d error=%q headers=%s payload=%q",
			rid,
			c.Request.Method,
			routePath,
			status,
			latency,
			c.Writer.Size(),
			errMsg,
			responseHeaders,
			truncateForLog(respBody, maxHTTPBodyLogBytes),
		)
	}
}

func readRequestBody(c *gin.Context) string {
	if c.Request == nil || c.Request.Body == nil {
		return ""
	}
	if !shouldLogBody(c.ContentType()) {
		return ""
	}

	raw, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return fmt.Sprintf("<read-error: %v>", err)
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(raw))
	return string(raw)
}

func shouldLogBody(contentType string) bool {
	if strings.HasPrefix(contentType, "multipart/form-data") {
		return false
	}
	return true
}

func headersToLog(headers http.Header) string {
	if headers == nil {
		return "{}"
	}

	sanitized := map[string][]string{}
	for k, vals := range headers {
		if isSensitiveHeader(k) {
			sanitized[k] = []string{"<redacted>"}
			continue
		}

		cloned := make([]string, len(vals))
		copy(cloned, vals)
		sanitized[k] = cloned
	}

	b, err := json.Marshal(sanitized)
	if err != nil {
		return `{"error":"marshal headers"}`
	}
	return truncateForLog(string(b), maxHTTPHeaderLogBytes)
}

func isSensitiveHeader(key string) bool {
	switch strings.ToLower(strings.TrimSpace(key)) {
	case "authorization", "cookie", "set-cookie", "x-api-key", "x-auth-token":
		return true
	default:
		return false
	}
}

func truncateForLog(v string, max int) string {
	if max <= 0 || len(v) <= max {
		return v
	}
	return v[:max] + "...(truncated)"
}
