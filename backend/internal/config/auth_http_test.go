package config

import (
	"net/http"
	"testing"
)

func TestLoadAuthHTTPConfig(t *testing.T) {
	for _, tt := range []struct {
		name, origins, secure, sameSite, proxies, mode string
		wantErr                                        bool
	}{
		{name: "secure defaults", origins: "https://app.example.com"},
		{name: "local HTTP", origins: "http://localhost:3000", secure: "false", sameSite: "lax"},
		{name: "cross site HTTPS", origins: "https://app.example.com,https://admin.example.com", sameSite: "none", proxies: "127.0.0.1,10.0.0.0/24", mode: "release"},
		{name: "strict", origins: "https://app.example.com", sameSite: "strict"},
		{name: "missing origins", wantErr: true},
		{name: "wildcard", origins: "*", wantErr: true},
		{name: "wildcard host", origins: "https://*.example.com", wantErr: true},
		{name: "null", origins: "null", wantErr: true},
		{name: "path", origins: "https://app.example.com/", wantErr: true},
		{name: "userinfo", origins: "https://user@app.example.com", wantErr: true},
		{name: "query", origins: "https://app.example.com?", wantErr: true},
		{name: "fragment", origins: "https://app.example.com#", wantErr: true},
		{name: "insecure none", origins: "https://app.example.com", secure: "false", sameSite: "none", wantErr: true},
		{name: "insecure release", origins: "https://app.example.com", secure: "false", mode: "release", wantErr: true},
		{name: "HTTP release", origins: "http://app.example.com", mode: "release", wantErr: true},
		{name: "invalid bool", origins: "https://app.example.com", secure: "maybe", wantErr: true},
		{name: "invalid same site", origins: "https://app.example.com", sameSite: "disabled", wantErr: true},
		{name: "trust every IPv4", origins: "https://app.example.com", proxies: "0.0.0.0/0", wantErr: true},
		{name: "trust every IPv6", origins: "https://app.example.com", proxies: "::/0", wantErr: true},
		{name: "invalid proxy", origins: "https://app.example.com", proxies: "proxy.local", wantErr: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("AUTH_ALLOWED_ORIGINS", tt.origins)
			t.Setenv("AUTH_COOKIE_SECURE", tt.secure)
			t.Setenv("AUTH_COOKIE_SAME_SITE", tt.sameSite)
			t.Setenv("TRUSTED_PROXIES", tt.proxies)
			t.Setenv("GIN_MODE", tt.mode)
			cfg, err := LoadAuthHTTPConfig()
			if (err != nil) != tt.wantErr {
				t.Fatalf("unexpected config error: %v", err)
			}
			if !tt.wantErr && tt.secure == "" && !cfg.CookieSecure {
				t.Fatal("cookies must default to secure")
			}
			if !tt.wantErr && tt.sameSite == "" && cfg.CookieSameSite != http.SameSiteLaxMode {
				t.Fatal("cookies must default to SameSite=Lax")
			}
		})
	}
}
