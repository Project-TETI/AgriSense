package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Project-TETI/AgriSense/backend/internal/config"
	"github.com/Project-TETI/AgriSense/backend/internal/domain"
	"github.com/Project-TETI/AgriSense/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func TestAuthRoutesEnforceSecurityAndLimits(t *testing.T) {
	gin.SetMode(gin.TestMode)
	stub := &authServiceStub{session: &domain.AuthSession{
		Response:     domain.AuthResponse{AccessToken: "access"},
		RefreshToken: strings.Repeat("a", 64), RefreshExpiresAt: time.Now().Add(time.Hour),
	}}
	r := gin.New()
	if err := r.SetTrustedProxies(nil); err != nil {
		t.Fatal(err)
	}
	r.Use(middleware.BrowserSecurity([]string{"https://app.example.com"}))
	h := NewAuthHandler(stub, config.AuthHTTPConfig{CookieSecure: true, CookieSameSite: http.SameSiteStrictMode})
	h.RegisterRoutes(r.Group("/api/v1/auth"), strings.Repeat("x", 32))
	request := func(path, body, origin, csrf string) *httptest.ResponseRecorder {
		req := httptest.NewRequest("POST", "/api/v1/auth/"+path, strings.NewReader(body))
		req.Header.Set("Origin", origin)
		req.Header.Set("X-CSRF-Protection", csrf)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w
	}
	for _, path := range []string{"register", "login", "refresh", "logout"} {
		if w := request(path, "{}", "https://evil.example", "1"); w.Code != 403 {
			t.Fatalf("foreign %s reached handler: %d", path, w.Code)
		}
		if w := request(path, "{}", "https://app.example.com", ""); w.Code != 403 {
			t.Fatalf("missing CSRF header accepted for %s", path)
		}
		if w := request(path, strings.Repeat("x", int(middleware.MaxAuthBodyBytes)+1), "https://app.example.com", "1"); w.Code != 413 {
			t.Fatalf("oversized %s reached handler: %d", path, w.Code)
		}
	}
	if stub.calls != 0 {
		t.Fatal("rejected request reached service")
	}
	// Logout must work without a bearer token, including repeated requests.
	for i := 0; i < 2; i++ {
		w := request("logout", "", "https://app.example.com", "1")
		cookies := w.Result().Cookies()
		if w.Code != 204 || w.Body.Len() != 0 || len(cookies) != 1 {
			t.Fatalf("logout failed: %d %s", w.Code, w.Body.String())
		}
		cookie := cookies[0]
		if cookie.MaxAge != -1 || cookie.Path != CookieAuthPath || !cookie.Secure || !cookie.HttpOnly ||
			cookie.Domain != "" || cookie.SameSite != http.SameSiteStrictMode {
			t.Fatal("deletion cookie does not match session cookie scope")
		}
	}
	body := `{"email":"test@example.com","password":"test-password"}`
	for i := 0; i < 9; i++ { // Oversized request above already consumed one login attempt.
		if w := request("login", body, "https://app.example.com", "1"); w.Code != 200 {
			t.Fatalf("login blocked early: %d", w.Code)
		}
	}
	before := stub.calls
	if w := request("login", body, "https://app.example.com", "1"); w.Code != 429 || w.Header().Get("Retry-After") == "" || stub.calls != before {
		t.Fatal("login limiter did not stop the service call")
	}
	if w := request("logout", "", "https://app.example.com", "1"); w.Code != 204 {
		t.Fatal("login limit must not block logout")
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/api/v1/auth/me", nil))
	if w.Code != 401 {
		t.Fatal("profile became public")
	}
}
