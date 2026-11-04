package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Project-TETI/AgriSense/backend/internal/config"
	"github.com/Project-TETI/AgriSense/backend/internal/domain"
	"github.com/Project-TETI/AgriSense/backend/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type authServiceStub struct {
	session *domain.AuthSession
	err     error
	calls   int
}

func (s *authServiceStub) Register(context.Context, domain.RegisterRequest) (*domain.AuthSession, error) {
	s.calls++
	return s.session, s.err
}
func (s *authServiceStub) Login(context.Context, domain.LoginRequest) (*domain.AuthSession, error) {
	s.calls++
	return s.session, s.err
}
func (s *authServiceStub) RefreshToken(context.Context, string) (*domain.AuthSession, error) {
	s.calls++
	return s.session, s.err
}
func (s *authServiceStub) Logout(context.Context, string) error {
	s.calls++
	return s.err
}
func (s *authServiceStub) GetProfile(context.Context, uuid.UUID) (*domain.UserResponse, error) {
	s.calls++
	if s.err != nil {
		return nil, s.err
	}
	return &s.session.Response.User, nil
}

func TestAuthHandlerErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		operation string
		err       error
		status    int
		code      string
	}{
		{"register", domain.ErrEmailAlreadyRegistered, 409, "EMAIL_ALREADY_REGISTERED"},
		{"register", domain.ErrInvalidInput, 400, "INVALID_INPUT"},
		{"register", domain.ErrAuthUnavailable, 503, "AUTH_SERVICE_UNAVAILABLE"},
		{"login", domain.ErrInvalidCredentials, 401, "INVALID_CREDENTIALS"},
		{"login", errors.New("sensitive database details"), 500, "INTERNAL_ERROR"},
		{"refresh", domain.ErrInvalidRefreshToken, 401, "INVALID_REFRESH_TOKEN"},
		{"refresh", domain.ErrAuthUnavailable, 503, "AUTH_SERVICE_UNAVAILABLE"},
		{"refresh", errors.New("sensitive database details"), 500, "INTERNAL_ERROR"},
		{"logout", domain.ErrAuthUnavailable, 503, "AUTH_SERVICE_UNAVAILABLE"},
		{"me", domain.ErrUserNotFound, 404, "USER_NOT_FOUND"},
		{"me", domain.ErrAuthUnavailable, 503, "AUTH_SERVICE_UNAVAILABLE"},
	}
	for _, tt := range tests {
		t.Run(tt.operation+"/"+tt.code, func(t *testing.T) {
			stub := &authServiceStub{err: fmt.Errorf("service: %w", tt.err)}
			h := NewAuthHandler(stub, config.AuthHTTPConfig{CookieSecure: true, CookieSameSite: http.SameSiteLaxMode})
			r := gin.New()
			r.POST("/register", h.Register)
			r.POST("/login", h.Login)
			r.POST("/refresh", h.RefreshToken)
			r.POST("/logout", h.Logout)
			r.GET("/me", func(c *gin.Context) { c.Set("user_id", uuid.New()); h.GetMe(c) })
			method := http.MethodPost
			if tt.operation == "me" {
				method = http.MethodGet
			}
			req := httptest.NewRequest(method, "/"+tt.operation, strings.NewReader(`{"name":"Test User","email":"test@example.com","password":"test-password"}`))
			req.Header.Set("Content-Type", "application/json")
			req.AddCookie(&http.Cookie{Name: RefreshTokenCookieName, Value: strings.Repeat("a", 64)})
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			var body struct {
				Error struct{ Code, Message string }
			}
			if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			if w.Code != tt.status || body.Error.Code != tt.code || stub.calls != 1 {
				t.Fatalf("unexpected response: status=%d body=%s calls=%d", w.Code, w.Body.String(), stub.calls)
			}
			if len(w.Result().Cookies()) != 0 || w.Header().Get("Cache-Control") != "no-store" {
				t.Fatal("failed auth must not change cookies or be cacheable")
			}
			if strings.Contains(w.Body.String(), "sensitive database details") {
				t.Fatal("internal error leaked")
			}
		})
	}
}

func TestAuthSuccessKeepsRefreshOutOfJSONAndUsesActualExpiry(t *testing.T) {
	gin.SetMode(gin.TestMode)
	expiry := time.Now().UTC().Add(2 * time.Hour).Truncate(time.Second)
	session := &domain.AuthSession{
		Response:     domain.AuthResponse{AccessToken: "access-token", TokenType: "Bearer", ExpiresIn: 900},
		RefreshToken: strings.Repeat("a", 64), RefreshExpiresAt: expiry,
	}
	for _, operation := range []string{"register", "login", "refresh"} {
		t.Run(operation, func(t *testing.T) {
			h := NewAuthHandler(&authServiceStub{session: session}, config.AuthHTTPConfig{CookieSecure: true, CookieSameSite: http.SameSiteLaxMode})
			r := gin.New()
			r.POST("/register", h.Register)
			r.POST("/login", h.Login)
			r.POST("/refresh", h.RefreshToken)
			req := httptest.NewRequest("POST", "/"+operation, strings.NewReader(`{"name":"Test User","email":"test@example.com","password":"test-password"}`))
			req.Header.Set("Content-Type", "application/json")
			req.AddCookie(&http.Cookie{Name: RefreshTokenCookieName, Value: strings.Repeat("b", 64)})
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			wantStatus := 200
			if operation == "register" {
				wantStatus = 201
			}
			cookies := w.Result().Cookies()
			if w.Code != wantStatus || len(cookies) != 1 {
				t.Fatalf("unexpected response: %d %s", w.Code, w.Body.String())
			}
			cookie := cookies[0]
			if cookie.Value != session.RefreshToken || !cookie.Expires.Equal(expiry) ||
				cookie.MaxAge <= 0 || cookie.MaxAge > 7200 || !cookie.HttpOnly || !cookie.Secure || cookie.Path != CookieAuthPath ||
				cookie.Domain != "" || cookie.SameSite != http.SameSiteLaxMode {
				t.Fatal("unexpected refresh cookie attributes")
			}
			if strings.Contains(w.Body.String(), session.RefreshToken) || !strings.Contains(w.Body.String(), `"access_token":"access-token"`) {
				t.Fatal("unexpected token serialization")
			}
		})
	}
	encoded, err := json.Marshal(session)
	if err != nil || strings.Contains(string(encoded), session.RefreshToken) {
		t.Fatal("internal session serialization exposes refresh token")
	}
}

func TestAuthInputAndMiddlewareFailures(t *testing.T) {
	gin.SetMode(gin.TestMode)
	stub := &authServiceStub{}
	h := NewAuthHandler(stub, config.AuthHTTPConfig{CookieSameSite: http.SameSiteLaxMode})
	r := gin.New()
	r.POST("/register", h.Register)
	r.POST("/refresh", h.RefreshToken)
	r.GET("/me", middleware.AuthMiddleware(strings.Repeat("x", 32)), h.GetMe)
	for _, req := range []*http.Request{
		httptest.NewRequest("POST", "/register", strings.NewReader(`{}`)),
		httptest.NewRequest("POST", "/refresh", nil),
		httptest.NewRequest("GET", "/me", nil),
	} {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code < 400 || len(w.Result().Cookies()) != 0 || stub.calls != 0 {
			t.Fatalf("invalid input reached service: status=%d calls=%d", w.Code, stub.calls)
		}
		if req.URL.Path == "/me" && (w.Header().Get("WWW-Authenticate") != "Bearer" || !strings.Contains(w.Body.String(), "INVALID_ACCESS_TOKEN")) {
			t.Fatal("middleware error does not follow auth contract")
		}
	}
}
