package handler

import (
	"net/http"

	"github.com/Project-TETI/AgriSense/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type PlantHandler struct {
	plantService service.PlantService
}

func NewPlantHandler(plantService service.PlantService) *PlantHandler {
	return &PlantHandler{plantService: plantService}
}

func (h *PlantHandler) GetPlants(c *gin.Context) {
	plants, err := h.plantService.GetAllPlants(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"plants": plants})
}
