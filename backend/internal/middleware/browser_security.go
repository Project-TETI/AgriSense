package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func abortAuthRequest(c *gin.Context, status int, code, message string) {
	c.Header("Cache-Control", "no-store")
	c.AbortWithStatusJSON(status, gin.H{"error": gin.H{"code": code, "message": message}})
}

// BrowserSecurity checks origins before any handler that can change a session.
func BrowserSecurity(origins []string) gin.HandlerFunc {
	allowed := make(map[string]bool, len(origins))
	for _, origin := range origins {
		allowed[origin] = true
	}
	return func(c *gin.Context) {
		c.Writer.Header().Add("Vary", "Origin")
		origin := c.GetHeader("Origin")
		if len(c.Request.Header.Values("Origin")) > 1 || (origin != "" && !allowed[origin]) {
			abortAuthRequest(c, http.StatusForbidden, "ORIGIN_NOT_ALLOWED", "Origin tidak diizinkan")
			return
		}
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Expose-Headers", "Retry-After, WWW-Authenticate")
		}
		if c.Request.Method == http.MethodOptions {
			c.Writer.Header().Add("Vary", "Access-Control-Request-Method")
			c.Writer.Header().Add("Vary", "Access-Control-Request-Headers")
			method := c.GetHeader("Access-Control-Request-Method")
			if origin == "" || (method != "POST" && method != "GET") {
				abortAuthRequest(c, http.StatusForbidden, "CORS_NOT_ALLOWED", "Preflight tidak diizinkan")
				return
			}
			for _, header := range strings.Split(c.GetHeader("Access-Control-Request-Headers"), ",") {
				switch strings.ToLower(strings.TrimSpace(header)) {
				case "", "content-type", "authorization", "x-csrf-protection":
				default:
					abortAuthRequest(c, http.StatusForbidden, "CORS_NOT_ALLOWED", "Header preflight tidak diizinkan")
					return
				}
			}
			c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-CSRF-Protection")
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		if c.Request.Method == http.MethodPost && strings.HasPrefix(c.Request.URL.Path, "/api/v1/auth/") {
			if origin == "" || c.GetHeader("X-CSRF-Protection") != "1" {
				abortAuthRequest(c, http.StatusForbidden, "CSRF_REJECTED", "Origin dan header X-CSRF-Protection: 1 diperlukan")
				return
			}
		}
		c.Next()
	}
}
