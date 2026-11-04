package middleware

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const MaxAuthBodyBytes int64 = 16 << 10

// AuthBodyLimit also bounds chunked bodies before JSON parsing or password hashing.
func AuthBodyLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(http.MaxBytesReader(c.Writer, c.Request.Body, MaxAuthBodyBytes))
		if err != nil {
			if _, tooLarge := err.(*http.MaxBytesError); tooLarge {
				abortAuthRequest(c, http.StatusRequestEntityTooLarge, "BODY_TOO_LARGE", "Maximum body size exceeded")
			} else {
				abortAuthRequest(c, http.StatusBadRequest, "INVALID_INPUT", "Body could not be read")
			}
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewReader(body))
		c.Next()
	}
}

type authWindow struct {
	count int
	until time.Time
}

// shortcut: limits are per process; use a shared gateway/store before running multiple replicas.
func AuthRateLimit(limit int, window time.Duration) gin.HandlerFunc {
	var mu sync.Mutex
	windows := make(map[string]authWindow)
	nextCleanup := time.Time{}
	return func(c *gin.Context) {
		now := time.Now()
		key := c.ClientIP()
		mu.Lock()
		if !now.Before(nextCleanup) {
			for ip, entry := range windows {
				if !now.Before(entry.until) {
					delete(windows, ip)
				}
			}
			nextCleanup = now.Add(window)
		}
		entry, exists := windows[key]
		if exists && !now.Before(entry.until) {
			delete(windows, key)
			exists = false
		}
		// Fail closed at capacity instead of growing memory without bound.
		if !exists && len(windows) >= 10000 {
			mu.Unlock()
			c.Header("Retry-After", strconv.Itoa(max(1, int(window.Seconds()))))
			abortAuthRequest(c, http.StatusTooManyRequests, "RATE_LIMITED", "Too many requests")
			return
		}
		if !exists {
			entry = authWindow{until: now.Add(window)}
		}
		allowed := entry.count < limit
		if allowed {
			entry.count++
		}
		windows[key] = entry
		mu.Unlock()
		if !allowed {
			c.Header("Retry-After", strconv.Itoa(max(1, int(time.Until(entry.until).Seconds())+1)))
			abortAuthRequest(c, http.StatusTooManyRequests, "RATE_LIMITED", "Too many requests")
			return
		}
		c.Next()
	}
}
