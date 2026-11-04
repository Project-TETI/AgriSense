package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	accessTokenIssuer   = "agrisense"
	accessTokenAudience = "agrisense-api"
)

type JWTClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Name   string    `json:"name"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userID uuid.UUID, email, name, secret string, duration time.Duration) (string, error) {
	if userID == uuid.Nil {
		return "", errors.New("user ID not valid")
	}
	if duration < time.Second {
		return "", errors.New("access token duration must be at least 1 second")
	}

	now := time.Now().UTC()
	return GenerateAccessTokenWithExpiry(userID, email, name, secret, now, now.Add(duration))
}

// GenerateAccessTokenWithExpiry lets session rotation preserve its absolute expiry.
func GenerateAccessTokenWithExpiry(userID uuid.UUID, email, name, secret string, issuedAt, expiresAt time.Time) (string, error) {
	if userID == uuid.Nil {
		return "", errors.New("user ID not valid")
	}
	iat := jwt.NewNumericDate(issuedAt.UTC())
	exp := jwt.NewNumericDate(expiresAt.UTC())
	if !exp.Time.After(iat.Time) {
		return "", errors.New("expire time must be after issued time")
	}
	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		Name:   name,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    accessTokenIssuer,
			Audience:  jwt.ClaimStrings{accessTokenAudience},
			Subject:   userID.String(),
			IssuedAt:  iat,
			ExpiresAt: exp,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateAccessToken(tokenString, secret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaims{},
		func(_ *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
		jwt.WithIssuer(accessTokenIssuer),
		jwt.WithAudience(accessTokenAudience),
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("token not valid")
	}

	if claims.UserID == uuid.Nil || claims.Subject != claims.UserID.String() {
		return nil, errors.New("token identity not valid")
	}
	if claims.IssuedAt == nil || claims.ExpiresAt == nil {
		return nil, errors.New("token missing required time claims")
	}
	if !claims.ExpiresAt.Time.After(claims.IssuedAt.Time) {
		return nil, errors.New("expire time must be after issued time")
	}

	return claims, nil
}

func GenerateRefreshToken() (rawToken string, tokenHash string, err error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", err
	}

	rawToken = hex.EncodeToString(bytes)
	hash := sha256.Sum256([]byte(rawToken))
	tokenHash = hex.EncodeToString(hash[:])

	return rawToken, tokenHash, nil
}

func HashToken(rawToken string) string {
	hash := sha256.Sum256([]byte(rawToken))
	return hex.EncodeToString(hash[:])
}
