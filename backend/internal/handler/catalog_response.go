package handler

import (
	"time"

	"github.com/Project-TETI/AgriSense/backend/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func catalogError(c *gin.Context, status int, code, message string) {
	c.JSON(status, gin.H{"error": gin.H{"code": code, "message": message}})
}

type plantResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type diseaseResponse struct {
	ID             uuid.UUID     `json:"id"`
	PlantID        uuid.UUID     `json:"plant_id"`
	ClassName      string        `json:"class_name"`
	DisplayName    string        `json:"display_name"`
	Cause          string        `json:"cause"`
	Description    string        `json:"description"`
	Symptoms       *string       `json:"symptoms"`
	Recommendation string        `json:"recommendation"`
	Plant          plantResponse `json:"plant"`
}

func catalogPlant(p domain.Plant) plantResponse {
	return plantResponse{p.ID, p.Name, p.DisplayName, p.Description, p.CreatedAt}
}

func catalogDisease(d domain.Disease) diseaseResponse {
	return diseaseResponse{d.ID, d.PlantID, d.ClassName, d.DisplayName, d.Cause, d.Description, d.Symptoms, d.Recommendation, catalogPlant(d.Plant)}
}
