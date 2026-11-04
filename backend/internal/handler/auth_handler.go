package handler

import (
	"errors"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/Project-TETI/AgriSense/backend/internal/config"
	"github.com/Project-TETI/AgriSense/backend/internal/domain"
	"github.com/Project-TETI/AgriSense/backend/internal/middleware"
	"github.com/Project-TETI/AgriSense/backend/internal/service"
	"github.com/gin-gonic/gin"
)

const (
	RefreshTokenCookieName = "refresh_token"
	CookieAuthPath         = "/api/v1/auth"
)

type AuthHandler struct {
	authService service.AuthService
	httpConfig  config.AuthHTTPConfig
}

func NewAuthHandler(authService service.AuthService, httpConfig config.AuthHTTPConfig) *AuthHandler {
	return &AuthHandler{authService: authService, httpConfig: httpConfig}
}

// RegisterRoutes is shared by the server and HTTP integration tests.
func (h *AuthHandler) RegisterRoutes(group *gin.RouterGroup, jwtSecret string) {
	posts := group.Group("")
	posts.POST("/register", middleware.AuthRateLimit(5, time.Minute), middleware.AuthBodyLimit(), h.Register)
	posts.POST("/login", middleware.AuthRateLimit(10, time.Minute), middleware.AuthBodyLimit(), h.Login)
	posts.POST("/refresh", middleware.AuthRateLimit(60, time.Minute), middleware.AuthBodyLimit(), h.RefreshToken)
	posts.POST("/logout", middleware.AuthRateLimit(60, time.Minute), middleware.AuthBodyLimit(), h.Logout)
	group.GET("/me", middleware.AuthMiddleware(jwtSecret), h.GetMe)
}

func writeAuthError(c *gin.Context, err error) {
	c.Header("Cache-Control", "no-store")
	status, code, message := http.StatusInternalServerError, "INTERNAL_ERROR", "Terjadi kesalahan internal"
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		status, code, message = http.StatusBadRequest, "INVALID_INPUT", "Input tidak valid"
	case errors.Is(err, domain.ErrEmailAlreadyRegistered):
		status, code, message = http.StatusConflict, "EMAIL_ALREADY_REGISTERED", "Email sudah terdaftar"
	case errors.Is(err, domain.ErrInvalidCredentials):
		status, code, message = http.StatusUnauthorized, "INVALID_CREDENTIALS", "Email atau kata sandi tidak valid"
	case errors.Is(err, domain.ErrInvalidAccessToken):
		status, code, message = http.StatusUnauthorized, "INVALID_ACCESS_TOKEN", "Token akses tidak valid atau kedaluwarsa"
		c.Header("WWW-Authenticate", "Bearer")
	case errors.Is(err, domain.ErrInvalidRefreshToken):
		status, code, message = http.StatusUnauthorized, "INVALID_REFRESH_TOKEN", "Refresh token tidak valid atau kedaluwarsa"
	case errors.Is(err, domain.ErrUserNotFound):
		status, code, message = http.StatusNotFound, "USER_NOT_FOUND", "Pengguna tidak ditemukan"
	case errors.Is(err, domain.ErrAuthUnavailable):
		status, code, message = http.StatusServiceUnavailable, "AUTH_SERVICE_UNAVAILABLE", "Layanan autentikasi sementara tidak tersedia"
	}
	if status >= 500 {
		// Do not log SQL errors here: driver details can contain credential hashes.
		log.Printf("auth request failed method=%s path=%s error_type=%T", c.Request.Method, c.FullPath(), err)
	}
	c.JSON(status, gin.H{"error": gin.H{"code": code, "message": message}})
}

func (h *AuthHandler) Register(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	var req domain.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeAuthError(c, domain.ErrInvalidInput)
		return
	}
	session, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		writeAuthError(c, err)
		return
	}
	h.setRefreshTokenCookie(c, session.RefreshToken, session.RefreshExpiresAt)
	c.JSON(http.StatusCreated, session.Response)
}

func (h *AuthHandler) Login(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeAuthError(c, domain.ErrInvalidInput)
		return
	}
	session, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		writeAuthError(c, err)
		return
	}
	h.setRefreshTokenCookie(c, session.RefreshToken, session.RefreshExpiresAt)
	c.JSON(http.StatusOK, session.Response)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	raw, err := c.Cookie(RefreshTokenCookieName)
	if err != nil {
		writeAuthError(c, domain.ErrInvalidRefreshToken)
		return
	}
	session, err := h.authService.RefreshToken(c.Request.Context(), raw)
	if err != nil {
		// A failed parallel refresh must not erase a newer cookie set by its peer.
		writeAuthError(c, err)
		return
	}
	h.setRefreshTokenCookie(c, session.RefreshToken, session.RefreshExpiresAt)
	c.JSON(http.StatusOK, session.Response)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	raw, _ := c.Cookie(RefreshTokenCookieName)
	if err := h.authService.Logout(c.Request.Context(), raw); err != nil {
		writeAuthError(c, err)
		return
	}
	h.clearRefreshTokenCookie(c)
	c.Status(http.StatusNoContent)
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	userID, err := middleware.GetCurrentUserID(c)
	if err != nil {
		writeAuthError(c, domain.ErrInvalidAccessToken)
		return
	}
	user, err := h.authService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		writeAuthError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *AuthHandler) setRefreshTokenCookie(c *gin.Context, token string, expiresAt time.Time) {
	maxAge := int(math.Ceil(time.Until(expiresAt).Seconds()))
	if maxAge <= 0 {
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name: RefreshTokenCookieName, Value: token, Path: CookieAuthPath,
		Expires: expiresAt.UTC(), MaxAge: maxAge, HttpOnly: true,
		Secure: h.httpConfig.CookieSecure, SameSite: h.httpConfig.CookieSameSite,
	})
}

func (h *AuthHandler) clearRefreshTokenCookie(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name: RefreshTokenCookieName, Path: CookieAuthPath, MaxAge: -1,
		Expires: time.Unix(1, 0).UTC(), HttpOnly: true,
		Secure: h.httpConfig.CookieSecure, SameSite: h.httpConfig.CookieSameSite,
	})
}
