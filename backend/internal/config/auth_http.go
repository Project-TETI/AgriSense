package config

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// AuthHTTPConfig controls browser credentials and the proxy trust boundary.
type AuthHTTPConfig struct {
	CookieSecure   bool
	CookieSameSite http.SameSite
	AllowedOrigins []string
	TrustedProxies []string
}

func LoadAuthHTTPConfig() (AuthHTTPConfig, error) {
	cfg := AuthHTTPConfig{CookieSecure: true, CookieSameSite: http.SameSiteLaxMode}
	if raw := strings.TrimSpace(os.Getenv("AUTH_COOKIE_SECURE")); raw != "" {
		secure, err := strconv.ParseBool(raw)
		if err != nil {
			return cfg, fmt.Errorf("AUTH_COOKIE_SECURE must be true or false")
		}
		cfg.CookieSecure = secure
	}
	switch strings.ToLower(strings.TrimSpace(os.Getenv("AUTH_COOKIE_SAME_SITE"))) {
	case "", "lax":
	case "strict":
		cfg.CookieSameSite = http.SameSiteStrictMode
	case "none":
		cfg.CookieSameSite = http.SameSiteNoneMode
	default:
		return cfg, fmt.Errorf("AUTH_COOKIE_SAME_SITE must be lax, strict, or none")
	}
	if !cfg.CookieSecure && (cfg.CookieSameSite == http.SameSiteNoneMode || os.Getenv("GIN_MODE") == "release") {
		return cfg, fmt.Errorf("secure cookies are required for SameSite=None and release mode")
	}
	for _, raw := range strings.Split(os.Getenv("AUTH_ALLOWED_ORIGINS"), ",") {
		origin := strings.TrimSpace(raw)
		u, err := url.Parse(origin)
		if err != nil || u == nil || (u.Scheme != "http" && u.Scheme != "https") ||
			u.Hostname() == "" || u.User != nil || u.Opaque != "" || u.Path != "" ||
			u.RawQuery != "" || u.ForceQuery || u.Fragment != "" || strings.ContainsAny(origin, "*#\\") {
			return cfg, fmt.Errorf("AUTH_ALLOWED_ORIGINS must contain explicit http(s) origins without paths or wildcards")
		}
		if os.Getenv("GIN_MODE") == "release" && u.Scheme != "https" {
			return cfg, fmt.Errorf("AUTH_ALLOWED_ORIGINS must use HTTPS in release mode")
		}
		cfg.AllowedOrigins = append(cfg.AllowedOrigins, origin)
	}
	if raw := strings.TrimSpace(os.Getenv("TRUSTED_PROXIES")); raw != "" {
		for _, entry := range strings.Split(raw, ",") {
			entry = strings.TrimSpace(entry)
			if net.ParseIP(entry) == nil {
				_, network, err := net.ParseCIDR(entry)
				if err != nil {
					return cfg, fmt.Errorf("TRUSTED_PROXIES must contain explicit proxy IPs or CIDRs, not all addresses")
				}
				if bits, _ := network.Mask.Size(); bits == 0 {
					return cfg, fmt.Errorf("TRUSTED_PROXIES must not trust all addresses")
				}
			}
			cfg.TrustedProxies = append(cfg.TrustedProxies, entry)
		}
	}
	return cfg, nil
}
