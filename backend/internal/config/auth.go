package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type AuthConfig struct {
	JWTSecret       string
	AccessDuration  time.Duration
	RefreshDuration time.Duration
}

func LoadAuthConfig() (AuthConfig, error) {
	secret := strings.TrimSpace(os.Getenv("JWT_SECRET"))

	if len(secret) < 32 {
		return AuthConfig{}, fmt.Errorf("Must provide a JWT secret with at least 32 characters. Current length: %d", len(secret))
	}

	accessDuration, err := readDurationEnv(
		"JWT_ACCESS_EXPIRATION",
		15*time.Minute,
	)
	if err != nil {
		return AuthConfig{}, err
	}

	refreshDuration, err := readDurationEnv(
		"JWT_REFRESH_EXPIRATION",
		7*24*time.Hour,
	)
	if err != nil {
		return AuthConfig{}, err
	}

	if refreshDuration <= accessDuration {
		return AuthConfig{}, fmt.Errorf(
			"JWT_REFRESH_EXPIRATION harus lebih panjang dari JWT_ACCESS_EXPIRATION",
		)
	}

	return AuthConfig{
		JWTSecret:       secret,
		AccessDuration:  accessDuration,
		RefreshDuration: refreshDuration,
	}, nil
}

func readDurationEnv(
	key string,
	defaultValue time.Duration,
) (time.Duration, error) {
	raw := strings.TrimSpace(os.Getenv(key))

	if raw == "" {
		return defaultValue, nil
	}

	duration, err := time.ParseDuration(raw)
	if err != nil || duration <= 0 {
		return 0, fmt.Errorf("%s must be a valid positive duration string, e.g., '15m', '1h', '24h'. Got: '%s'", key, raw)
	}

	return duration, nil
}
