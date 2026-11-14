package handler

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Project-TETI/AgriSense/backend/internal/service"
)

var catalogIdentifier = regexp.MustCompile(`^[a-z0-9]+(_[a-z0-9]+)*$`)

type DiseaseHandler struct {
	diseaseService service.DiseaseService
}

func NewDiseaseHandler(diseaseService service.DiseaseService) *DiseaseHandler {
	return &DiseaseHandler{diseaseService: diseaseService}
}

func (h *DiseaseHandler) GetDiseases(c *gin.Context) {
	var plantID string
	if values, present := c.Request.URL.Query()["plant_id"]; present {
		if len(values) != 1 {
			catalogError(c, http.StatusBadRequest, "INVALID_INPUT", "Filter tanaman tidak valid")
			return
		}
		id, err := uuid.Parse(values[0])
		if err != nil || id == uuid.Nil {
			catalogError(c, http.StatusBadRequest, "INVALID_INPUT", "Filter tanaman tidak valid")
			return
		}
		plantID = id.String()
	}
	diseases, err := h.diseaseService.GetAllDiseases(c.Request.Context(), plantID)
	if err != nil {
		catalogError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Terjadi kesalahan internal")
		return
	}

	response := make([]diseaseResponse, 0, len(diseases))
	for _, disease := range diseases {
		response = append(response, catalogDisease(disease))
	}
	c.JSON(http.StatusOK, gin.H{"diseases": response})
}

func (h *DiseaseHandler) GetDiseaseByPlantAndClass(c *gin.Context) {
	plantName := c.Param("plant_name")
	className := c.Param("class_name")
	if len(plantName) > 100 || len(className) > 150 || !catalogIdentifier.MatchString(plantName) || !catalogIdentifier.MatchString(className) {
		catalogError(c, http.StatusBadRequest, "INVALID_INPUT", "Kode tanaman atau kelas tidak valid")
		return
	}

	disease, err := h.diseaseService.GetDiseaseByPlantAndClass(c.Request.Context(), plantName, className)
	if err != nil {
		if errors.Is(err, service.ErrDiseaseNotFound) {
			catalogError(c, http.StatusNotFound, "DISEASE_NOT_FOUND", "Data penyakit tidak ditemukan")
		} else {
			catalogError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Terjadi kesalahan internal")
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"disease": catalogDisease(*disease),
	})
}
