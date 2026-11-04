package domain

import "errors"

var (
	ErrInvalidInput           = errors.New("invalid authentication input")
	ErrEmailAlreadyRegistered = errors.New("email already registered")
	ErrInvalidCredentials     = errors.New("invalid credentials")
	ErrInvalidAccessToken     = errors.New("invalid access token")
	ErrInvalidRefreshToken    = errors.New("invalid refresh token")
	ErrUserNotFound           = errors.New("user not found")
	ErrAuthUnavailable        = errors.New("authentication storage unavailable")
)
