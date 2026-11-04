package config

import (
	"strings"
	"testing"
	"time"
)

func TestLoadAuthConfig(t *testing.T) {
	validSecret := strings.Repeat("x", 32)

	tests := []struct {
		name        string
		secret      string
		access      string
		refresh     string
		wantAccess  time.Duration
		wantRefresh time.Duration
		wantErr     bool
	}{
		{
			name:        "default durations",
			secret:      validSecret,
			wantAccess:  15 * time.Minute,
			wantRefresh: 7 * 24 * time.Hour,
		},
		{
			name:        "custom durations",
			secret:      validSecret,
			access:      "10m",
			refresh:     "24h",
			wantAccess:  10 * time.Minute,
			wantRefresh: 24 * time.Hour,
		},
		{
			name:    "missing secret",
			wantErr: true,
		},
		{
			name:    "blank secret",
			secret:  "   ",
			wantErr: true,
		},
		{
			name:    "short secret",
			secret:  "short",
			wantErr: true,
		},
		{
			name:    "invalid access",
			secret:  validSecret,
			access:  "invalid",
			wantErr: true,
		},
		{
			name:    "invalid refresh",
			secret:  validSecret,
			refresh: "invalid",
			wantErr: true,
		},
		{
			name:    "zero access",
			secret:  validSecret,
			access:  "0s",
			wantErr: true,
		},
		{
			name:    "subsecond access",
			secret:  validSecret,
			access:  "500ms",
			wantErr: true,
		},
		{
			name:    "negative refresh",
			secret:  validSecret,
			refresh: "-1h",
			wantErr: true,
		},
		{
			name:    "equal durations",
			secret:  validSecret,
			access:  "1h",
			refresh: "1h",
			wantErr: true,
		},
		{
			name:    "shorter refresh",
			secret:  validSecret,
			access:  "2h",
			refresh: "1h",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("JWT_SECRET", tt.secret)
			t.Setenv("JWT_ACCESS_EXPIRATION", tt.access)
			t.Setenv("JWT_REFRESH_EXPIRATION", tt.refresh)

			got, err := LoadAuthConfig()

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected an error")
				}

				if tt.secret != "" &&
					strings.TrimSpace(tt.secret) != "" &&
					strings.Contains(err.Error(), tt.secret) {
					t.Fatal("error exposes secret")
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got.JWTSecret != tt.secret {
				t.Error("secret was not loaded correctly")
			}

			if got.AccessDuration != tt.wantAccess ||
				got.RefreshDuration != tt.wantRefresh {
				t.Errorf(
					"unexpected durations: access=%s refresh=%s",
					got.AccessDuration,
					got.RefreshDuration,
				)
			}
		})
	}
}
