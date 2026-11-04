package utils

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestGenerateAccessToken(t *testing.T) {
	secret := strings.Repeat("x", 32)
	userID := uuid.New()
	raw, err := GenerateAccessToken(userID, "contoh@example.com", "Pengguna", secret, 15*time.Minute)
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}
	claims, err := ValidateAccessToken(raw, secret)
	if err != nil {
		t.Fatalf("validate generated token: %v", err)
	}
	if claims.UserID != userID || claims.Subject != userID.String() ||
		claims.Email != "contoh@example.com" || claims.Name != "Pengguna" {
		t.Fatal("unexpected user claims")
	}
	if claims.Issuer != "agrisense" || len(claims.Audience) != 1 || claims.Audience[0] != "agrisense-api" {
		t.Fatal("unexpected issuer or audience")
	}
	if claims.ExpiresAt.Sub(claims.IssuedAt.Time) != 15*time.Minute {
		t.Fatal("unexpected token lifetime")
	}
}

func TestGenerateAccessTokenRejectsInvalidInput(t *testing.T) {
	secret := strings.Repeat("x", 32)
	userID := uuid.New()
	tests := []struct {
		name     string
		userID   uuid.UUID
		duration time.Duration
	}{
		{name: "nil user ID", duration: time.Minute},
		{name: "zero duration", userID: userID},
		{name: "negative duration", userID: userID, duration: -time.Minute},
		{name: "subsecond duration", userID: userID, duration: time.Millisecond},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := GenerateAccessToken(tt.userID, "", "", secret, tt.duration); err == nil {
				t.Fatal("expected invalid input to fail")
			}
		})
	}
}

func TestValidateAccessToken(t *testing.T) {
	secret := strings.Repeat("x", 32)
	now := time.Now().UTC()
	userID := uuid.New()
	tests := []struct {
		name       string
		change     func(*JWTClaims)
		method     jwt.SigningMethod
		signingKey interface{}
		wantErr    bool
	}{
		{name: "valid"},
		{
			name: "expired",
			change: func(c *JWTClaims) {
				c.IssuedAt = jwt.NewNumericDate(now.Add(-2 * time.Hour))
				c.ExpiresAt = jwt.NewNumericDate(now.Add(-time.Hour))
			},
			wantErr: true,
		},
		{
			name:    "missing expiration",
			change:  func(c *JWTClaims) { c.ExpiresAt = nil },
			wantErr: true,
		},
		{
			name:    "missing issued at",
			change:  func(c *JWTClaims) { c.IssuedAt = nil },
			wantErr: true,
		},
		{
			name: "future issued at",
			change: func(c *JWTClaims) {
				c.IssuedAt = jwt.NewNumericDate(now.Add(time.Hour))
				c.ExpiresAt = jwt.NewNumericDate(now.Add(2 * time.Hour))
			},
			wantErr: true,
		},
		{
			name:    "wrong issuer",
			change:  func(c *JWTClaims) { c.Issuer = "other-app" },
			wantErr: true,
		},
		{
			name:    "missing issuer",
			change:  func(c *JWTClaims) { c.Issuer = "" },
			wantErr: true,
		},
		{
			name:    "wrong audience",
			change:  func(c *JWTClaims) { c.Audience = jwt.ClaimStrings{"other-api"} },
			wantErr: true,
		},
		{
			name:    "missing audience",
			change:  func(c *JWTClaims) { c.Audience = nil },
			wantErr: true,
		},
		{
			name:    "missing subject",
			change:  func(c *JWTClaims) { c.Subject = "" },
			wantErr: true,
		},
		{
			name:    "mismatched subject",
			change:  func(c *JWTClaims) { c.Subject = uuid.NewString() },
			wantErr: true,
		},
		{
			name:    "nil user ID",
			change:  func(c *JWTClaims) { c.UserID = uuid.Nil },
			wantErr: true,
		},
		{
			name:    "wrong algorithm",
			method:  jwt.SigningMethodHS384,
			wantErr: true,
		},
		{
			name:       "unsigned token",
			method:     jwt.SigningMethodNone,
			signingKey: jwt.UnsafeAllowNoneSignatureType,
			wantErr:    true,
		},
		{
			name:       "wrong signing key",
			signingKey: []byte(strings.Repeat("y", 32)),
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims := JWTClaims{
				UserID: userID,
				Email:  "contoh@example.com",
				Name:   "Pengguna",
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "agrisense",
					Audience:  jwt.ClaimStrings{"agrisense-api"},
					Subject:   userID.String(),
					IssuedAt:  jwt.NewNumericDate(now.Add(-time.Minute)),
					ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
				},
			}
			if tt.change != nil {
				tt.change(&claims)
			}
			method := tt.method
			if method == nil {
				method = jwt.SigningMethodHS256
			}
			key := tt.signingKey
			if key == nil {
				key = []byte(secret)
			}
			raw, err := jwt.NewWithClaims(method, claims).SignedString(key)
			if err != nil {
				t.Fatalf("build test token: %v", err)
			}
			got, err := ValidateAccessToken(raw, secret)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error=%v, wantErr=%v", err, tt.wantErr)
			}
			if tt.wantErr && got != nil {
				t.Fatal("invalid token must not return claims")
			}
		})
	}
}

func TestValidateAccessTokenRejectsMalformedToken(t *testing.T) {
	for _, raw := range []string{"", "not-a-jwt"} {
		if claims, err := ValidateAccessToken(raw, strings.Repeat("x", 32)); err == nil || claims != nil {
			t.Fatal("expected malformed token to fail without returning claims")
		}
	}
}
