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
		catalogError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Terjadi kesalahan internal")
		return
	}

	response := make([]plantResponse, 0, len(plants))
	for _, plant := range plants {
		response = append(response, catalogPlant(plant))
	}
	c.JSON(http.StatusOK, gin.H{"plants": response})
}
