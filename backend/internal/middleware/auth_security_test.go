package middleware

import (
	"io"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestBrowserSecurity(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, tt := range []struct {
		name, method, path, origin, csrf, preflightMethod, preflightHeaders string
		status                                                              int
	}{
		{name: "allowed login", method: "POST", path: "/api/v1/auth/login", origin: "https://app.example.com", csrf: "1", status: 200},
		{name: "foreign login", method: "POST", path: "/api/v1/auth/login", origin: "https://evil.example", csrf: "1", status: 403},
		{name: "suffix trick", method: "POST", path: "/api/v1/auth/refresh", origin: "https://app.example.com.evil.example", csrf: "1", status: 403},
		{name: "missing origin", method: "POST", path: "/api/v1/auth/refresh", csrf: "1", status: 403},
		{name: "null origin", method: "POST", path: "/api/v1/auth/logout", origin: "null", csrf: "1", status: 403},
		{name: "missing CSRF", method: "POST", path: "/api/v1/auth/register", origin: "https://app.example.com", status: 403},
		{name: "wrong CSRF", method: "POST", path: "/api/v1/auth/logout", origin: "https://app.example.com", csrf: "0", status: 403},
		{name: "native read", method: "GET", path: "/api/v1/auth/me", status: 200},
		{name: "foreign read", method: "GET", path: "/api/v1/auth/me", origin: "https://evil.example", status: 403},
		{name: "allowed preflight", method: "OPTIONS", path: "/api/v1/auth/login", origin: "https://app.example.com", preflightMethod: "POST", preflightHeaders: "content-type, X-CSRF-Protection", status: 204},
		{name: "bearer preflight", method: "OPTIONS", path: "/api/v1/auth/me", origin: "https://app.example.com", preflightMethod: "GET", preflightHeaders: "Authorization", status: 204},
		{name: "foreign preflight", method: "OPTIONS", path: "/api/v1/auth/login", origin: "https://evil.example", preflightMethod: "POST", status: 403},
		{name: "bad method", method: "OPTIONS", path: "/api/v1/auth/login", origin: "https://app.example.com", preflightMethod: "DELETE", status: 403},
		{name: "bad header", method: "OPTIONS", path: "/api/v1/auth/login", origin: "https://app.example.com", preflightMethod: "POST", preflightHeaders: "X-Unknown", status: 403},
	} {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.Use(BrowserSecurity([]string{"https://app.example.com"}))
			calls := 0
			r.Handle(tt.method, tt.path, func(c *gin.Context) { calls++; c.Status(200) })
			req := httptest.NewRequest(tt.method, tt.path, nil)
			req.Header.Set("Origin", tt.origin)
			req.Header.Set("X-CSRF-Protection", tt.csrf)
			req.Header.Set("Access-Control-Request-Method", tt.preflightMethod)
			req.Header.Set("Access-Control-Request-Headers", tt.preflightHeaders)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			if w.Code != tt.status || (tt.status != 200 && calls != 0) {
				t.Fatalf("status=%d calls=%d body=%s", w.Code, calls, w.Body.String())
			}
			if tt.origin == "https://app.example.com" && (w.Header().Get("Access-Control-Allow-Origin") != tt.origin || w.Header().Get("Access-Control-Allow-Credentials") != "true") {
				t.Fatal("missing credentialed CORS for allowed origin")
			}
			if tt.origin != "https://app.example.com" && w.Header().Get("Access-Control-Allow-Origin") != "" {
				t.Fatal("untrusted origin received CORS permission")
			}
		})
	}
}

func TestAuthBodyLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, size := range []int{int(MaxAuthBodyBytes), int(MaxAuthBodyBytes) + 1} {
		r := gin.New()
		r.POST("/", AuthBodyLimit(), func(c *gin.Context) {
			body, err := io.ReadAll(c.Request.Body)
			if err != nil || len(body) != size {
				t.Error("body was not preserved")
			}
			c.Status(204)
		})
		req := httptest.NewRequest("POST", "/", strings.NewReader(strings.Repeat("x", size)))
		req.ContentLength = -1 // Unknown/chunked length must still be bounded.
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		want := 204
		if int64(size) > MaxAuthBodyBytes {
			want = 413
		}
		if w.Code != want {
			t.Fatalf("size=%d status=%d body=%s", size, w.Code, w.Body.String())
		}
	}
}

func TestAuthRateLimitConcurrentAndUntrustedForwarding(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	if err := r.SetTrustedProxies(nil); err != nil {
		t.Fatal(err)
	}
	var accepted atomic.Int32
	r.POST("/", AuthRateLimit(10, time.Minute), func(c *gin.Context) { accepted.Add(1); c.Status(204) })
	var wg sync.WaitGroup
	for i := 0; i < 30; i++ {
		wg.Go(func() {
			req := httptest.NewRequest("POST", "/", nil)
			req.RemoteAddr = "192.0.2.1:1234"
			req.Header.Set("X-Forwarded-For", "198.51.100.1")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			if w.Code != 204 && (w.Code != 429 || w.Header().Get("Retry-After") == "") {
				t.Errorf("unexpected rate limit response: %d", w.Code)
			}
		})
	}
	wg.Wait()
	if accepted.Load() != 10 {
		t.Fatalf("accepted=%d", accepted.Load())
	}
	for _, ip := range []string{"192.0.2.1:1234", "192.0.2.2:1234"} {
		req := httptest.NewRequest("POST", "/", nil)
		req.RemoteAddr = ip
		req.Header.Set("X-Forwarded-For", "203.0.113.55")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		want := 204
		if ip == "192.0.2.1:1234" {
			want = 429
		}
		if w.Code != want {
			t.Fatalf("forwarded IP bypass or independent IP blocked: %d", w.Code)
		}
	}
}

func TestAuthRateLimitWindowExpires(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/", AuthRateLimit(1, 50*time.Millisecond), func(c *gin.Context) { c.Status(204) })
	request := func() int {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest("POST", "/", nil))
		return w.Code
	}
	if request() != 204 {
		t.Fatal("first request blocked")
	}
	time.Sleep(60 * time.Millisecond)
	if request() != 204 {
		t.Fatal("window did not reset")
	}
}
