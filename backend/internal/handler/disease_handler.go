package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Project-TETI/AgriSense/backend/internal/service"
)

type DiseaseHandler struct {
	diseaseService service.DiseaseService
}

func NewDiseaseHandler(diseaseService service.DiseaseService) *DiseaseHandler {
	return &DiseaseHandler{diseaseService: diseaseService}
}

func (h *DiseaseHandler) GetDiseases(c *gin.Context) {
	plantID := c.Query("plant_id")
	diseases, err := h.diseaseService.GetAllDiseases(c.Request.Context(), plantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot retrieve diseases: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"diseases": diseases,
	})
}

// /api/v1/diseases/:class_name
func (h *DiseaseHandler) GetDiseaseByClassName(c *gin.Context) {
	className := c.Param("class_name")

	disease, err := h.diseaseService.GetDiseaseByClassName(c.Request.Context(), className)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cannot retrieve disease: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"disease": disease,
	})
}
