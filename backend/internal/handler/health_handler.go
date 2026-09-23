package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) CheckHealth(c *gin.Context) {
	dbStatus := "connected"
	isHealthy := true
	var errMsg string

	if h.db == nil {
		dbStatus = "disconnected"
		isHealthy = false
		errMsg = "Database connection is not initialized"
	} else {
		sqlDB, err := h.db.DB()
		if err != nil || sqlDB.Ping() != nil {
			dbStatus = "disconnected"
			isHealthy = false
			if err != nil {
				errMsg = err.Error()
			} else {
				errMsg = "Database ping failed"
			}
		}
	}

	res := gin.H{
		"timestamp": time.Now().Format(time.RFC3339),
		"services": gin.H{
			"database": dbStatus,
		},
	}

	if !isHealthy {
		res["status"] = "DEGRADED"
		res["error"] = errMsg
		c.JSON(http.StatusServiceUnavailable, res)
		return
	}

	res["status"] = "UP"
	c.JSON(http.StatusOK, res)
}
