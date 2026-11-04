package service_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Project-TETI/AgriSense/backend/internal/config"
	"github.com/Project-TETI/AgriSense/backend/internal/domain"
	"github.com/Project-TETI/AgriSense/backend/internal/handler"
	"github.com/Project-TETI/AgriSense/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func TestAuthHTTPBrowserSessionFlow(t *testing.T) {
	testDB := newAuthTestDB(t)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	if err := r.SetTrustedProxies(nil); err != nil {
		t.Fatal(err)
	}
	origin := "https://app.example.com"
	r.Use(middleware.BrowserSecurity([]string{origin}))
	h := handler.NewAuthHandler(authServiceFor(testDB.db), config.AuthHTTPConfig{
		CookieSecure: true, CookieSameSite: http.SameSiteNoneMode,
	})
	h.RegisterRoutes(r.Group("/api/v1/auth"), strings.Repeat("x", 32))
	server := httptest.NewTLSServer(r)
	defer server.Close()
	client := server.Client()
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	client.Jar = jar
	request := func(client *http.Client, method, path, body, bearer, raw string, want int) (*http.Response, []byte) {
		t.Helper()
		req, err := http.NewRequest(method, server.URL+"/api/v1/auth/"+path, strings.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Origin", origin)
		req.Header.Set("X-CSRF-Protection", "1")
		req.Header.Set("Content-Type", "application/json")
		if bearer != "" {
			req.Header.Set("Authorization", "Bearer "+bearer)
		}
		if raw != "" {
			req.AddCookie(&http.Cookie{Name: handler.RefreshTokenCookieName, Value: raw})
		}
		res, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		data, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != want || res.Header.Get("Cache-Control") != "no-store" {
			t.Fatalf("%s %s: status=%d body=%s", method, path, res.StatusCode, data)
		}
		return res, data
	}
	registration := `{"name":"Browser User","email":"browser@example.com","password":"test-password"}`
	res, data := request(client, "POST", "register", registration, "", "", 201)
	cookies := res.Cookies()
	if len(cookies) != 1 || !cookies[0].Secure || !cookies[0].HttpOnly || cookies[0].SameSite != http.SameSiteNoneMode || cookies[0].Domain != "" {
		t.Fatal("secure host-only refresh cookie missing")
	}
	if strings.Contains(string(data), cookies[0].Value) {
		t.Fatal("refresh token exposed in JSON")
	}
	var registered domain.AuthResponse
	if err := json.Unmarshal(data, &registered); err != nil {
		t.Fatal(err)
	}
	request(client, "GET", "me", "", "", "", 401)
	request(client, "GET", "me", "", registered.AccessToken, "", 200)
	request(client, "POST", "login", `{"email":"browser@example.com","password":"wrong-password"}`, "", "", 401)
	res, _ = request(client, "POST", "login", `{"email":"browser@example.com","password":"test-password"}`, "", "", 200)
	old := res.Cookies()[0]
	res, data = request(client, "POST", "refresh", "", "", "", 200)
	replacement := res.Cookies()[0]
	if replacement.Value == old.Value || !replacement.Expires.Equal(old.Expires) {
		t.Fatal("rotation must replace token without extending expiry")
	}
	var refreshed domain.AuthResponse
	if err := json.Unmarshal(data, &refreshed); err != nil {
		t.Fatal(err)
	}
	withoutJar := &http.Client{Transport: client.Transport}
	res, _ = request(withoutJar, "POST", "refresh", "", "", old.Value, 401)
	if len(res.Cookies()) != 0 {
		t.Fatal("losing refresh response must not erase replacement cookie")
	}
	// An expired or malformed access token must not prevent cookie-based logout.
	res, data = request(client, "POST", "logout", "", "expired-access-token", "", 204)
	if len(data) != 0 || len(res.Cookies()) != 1 || res.Cookies()[0].MaxAge != -1 || res.Cookies()[0].Path != replacement.Path || res.Cookies()[0].SameSite != replacement.SameSite {
		t.Fatal("logout must clear the same cookie and return an empty 204")
	}
	request(client, "POST", "logout", "", "", "", 204)
	request(client, "POST", "refresh", "", "", "", 401)
	request(withoutJar, "POST", "refresh", "", "", replacement.Value, 401)
	// Access JWTs remain valid until expiry; logout revokes the supplied refresh session.
	request(client, "GET", "me", "", refreshed.AccessToken, "", 200)
}
