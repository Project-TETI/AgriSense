package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Project-TETI/AgriSense/backend/internal/domain"
	"github.com/Project-TETI/AgriSense/backend/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			abortInvalidAccessToken(c)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			abortInvalidAccessToken(c)
			return
		}

		tokenString := parts[1]
		claims, err := utils.ValidateAccessToken(tokenString, jwtSecret)
		if err != nil {
			abortInvalidAccessToken(c)
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("name", claims.Name)

		c.Next()
	}
}

func abortInvalidAccessToken(c *gin.Context) {
	c.Header("WWW-Authenticate", "Bearer")
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"error": gin.H{"code": "INVALID_ACCESS_TOKEN", "message": domain.ErrInvalidAccessToken.Error()},
	})
}

func GetCurrentUserID(c *gin.Context) (uuid.UUID, error) {
	val, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, errors.New("user_id tidak ditemukan di context")
	}

	userID, ok := val.(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("tipe data user_id tidak valid")
	}

	return userID, nil
}
