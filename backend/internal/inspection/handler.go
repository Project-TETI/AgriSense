package inspection

import (
	"bytes"
	"context"
	"errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Project-TETI/AgriSense/backend/internal/domain"
	"github.com/Project-TETI/AgriSense/backend/internal/middleware"
	"github.com/Project-TETI/AgriSense/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type Handler struct {
	db       *gorm.DB
	ml       *MLClient
	storage  *Storage
	diseases service.DiseaseService
}

func NewHandler(db *gorm.DB, ml *MLClient, storage *Storage, diseases service.DiseaseService) *Handler {
	return &Handler{db: db, ml: ml, storage: storage, diseases: diseases}
}

func fail(c *gin.Context, status int, code, message string) {
	c.JSON(status, gin.H{"error": gin.H{"code": code, "message": message}})
}

func (h *Handler) Create(c *gin.Context) {
	userID, err := middleware.GetCurrentUserID(c)
	if err != nil {
		fail(c, 401, "INVALID_ACCESS_TOKEN", "Token akses tidak valid")
		return
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxImageSize+(64<<10))
	err = c.Request.ParseMultipartForm(MaxImageSize + (64 << 10))
	if c.Request.MultipartForm != nil {
		defer c.Request.MultipartForm.RemoveAll()
	}
	if err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			fail(c, 413, "FILE_TOO_LARGE", "Ukuran foto terlalu besar")
		} else {
			fail(c, 400, "INVALID_INPUT", "Formulir tidak valid")
		}
		return
	}
	form := c.Request.MultipartForm
	if len(form.Value) != 1 || len(form.Value["plant"]) != 1 || len(form.File) != 1 || len(form.File["file"]) != 1 {
		fail(c, 400, "INVALID_INPUT", "Satu foto dan satu kode tanaman diperlukan")
		return
	}
	plantName := form.Value["plant"][0]
	if len(plantName) > 100 || !codePattern.MatchString(plantName) {
		fail(c, 400, "INVALID_INPUT", "Kode tanaman tidak valid")
		return
	}
	file, err := form.File["file"][0].Open()
	if err != nil {
		fail(c, 400, "INVALID_IMAGE", "Foto tidak dapat dibaca")
		return
	}
	data, err := io.ReadAll(io.LimitReader(file, MaxImageSize+1))
	_ = file.Close()
	if err != nil {
		fail(c, 400, "INVALID_IMAGE", "Foto tidak dapat dibaca")
		return
	}
	if len(data) > MaxImageSize {
		fail(c, 413, "FILE_TOO_LARGE", "Ukuran foto maksimal 5 MB")
		return
	}
	config, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil || (format != "jpeg" && format != "png") || config.Width <= 0 || config.Height <= 0 || int64(config.Width)*int64(config.Height) > 20000000 {
		fail(c, 400, "INVALID_IMAGE", "Foto harus berupa JPEG atau PNG yang valid")
		return
	}
	if _, _, err := image.Decode(bytes.NewReader(data)); err != nil {
		fail(c, 400, "INVALID_IMAGE", "Foto tidak valid")
		return
	}
	var plant domain.Plant
	err = h.db.WithContext(c.Request.Context()).Where("name = ?", plantName).First(&plant).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fail(c, 400, "INVALID_PLANT", "Tanaman tidak ditemukan")
		return
	}
	if err != nil {
		fail(c, 500, "INTERNAL_ERROR", "Terjadi kesalahan internal")
		return
	}
	prediction, err := h.ml.Predict(c.Request.Context(), plantName, "image/"+format, data)
	if err != nil {
		var mlErr *MLError
		if errors.As(err, &mlErr) {
			message := "Layanan pemeriksaan tidak dapat memproses permintaan"
			if mlErr.Code == "BLURRY_IMAGE" {
				message = "Foto terlalu buram. Silakan ambil foto ulang."
			}
			fail(c, mlErr.Status, mlErr.Code, message)
			return
		}
		fail(c, 500, "INTERNAL_ERROR", "Terjadi kesalahan internal")
		return
	}
	var disease *domain.Disease
	var catalogAvailable *bool
	if prediction.Status == "SUCCESS" && prediction.ClassName != "healthy" {
		disease, err = h.diseases.GetDiseaseByPlantAndClass(c.Request.Context(), plantName, prediction.ClassName)
		if err != nil && !errors.Is(err, service.ErrDiseaseNotFound) {
			fail(c, 500, "INTERNAL_ERROR", "Terjadi kesalahan internal")
			return
		}
		found := disease != nil
		catalogAvailable = &found
	}
	key := uuid.NewString() + "." + format
	if err := h.storage.Put(c.Request.Context(), key, "image/"+format, data); err != nil {
		// shortcut: a lost Azure upload response may leave an orphan; add reconciliation before production.
		if h.storage.container != nil {
			log.Printf("inspection upload may need reconciliation for key=%s", key)
		}
		fail(c, 503, "STORAGE_UNAVAILABLE", "Penyimpanan foto tidak tersedia")
		return
	}
	inspection := domain.Inspection{UserID: userID, PlantID: plant.ID, ImageURL: key, PredictedClass: &prediction.ClassName, ConfidenceScore: prediction.Confidence, Status: prediction.Status}
	if disease != nil {
		inspection.DiseaseID = &disease.ID
	}
	if err := h.db.WithContext(c.Request.Context()).Omit("User", "Plant", "Disease").Create(&inspection).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && (strings.HasPrefix(pgErr.Code, "22") || strings.HasPrefix(pgErr.Code, "23")) {
			h.cleanup(key)
		} else {
			// shortcut: an ambiguous database failure may leave an orphan; add reconciliation before production.
			log.Printf("inspection image may need reconciliation for key=%s", key)
		}
		fail(c, 500, "INTERNAL_ERROR", "Hasil pemeriksaan tidak dapat disimpan")
		return
	}
	var diseaseInfo any
	if disease != nil {
		diseaseInfo = gin.H{"class_name": disease.ClassName, "display_name": disease.DisplayName, "cause": disease.Cause, "description": disease.Description, "symptoms": disease.Symptoms, "recommendation": disease.Recommendation}
	}
	response := gin.H{"inspection_id": inspection.ID, "plant": plantName, "predicted_class": prediction.ClassName, "confidence_score": *prediction.Confidence, "status": prediction.Status, "disease_info": diseaseInfo, "catalog_available": catalogAvailable, "image_url": "/api/v1/inspections/" + inspection.ID.String() + "/image", "created_at": inspection.CreatedAt}
	if prediction.Status == "UNCERTAIN" {
		response["message"] = "Hasil kurang meyakinkan. Ambil foto ulang dengan pencahayaan dan fokus yang lebih baik."
	}
	c.JSON(http.StatusCreated, response)
}

func (h *Handler) cleanup(key string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := h.storage.Delete(ctx, key); err != nil {
		log.Printf("inspection image cleanup failed for key=%s", key)
	}
}

func (h *Handler) Image(c *gin.Context) {
	userID, err := middleware.GetCurrentUserID(c)
	if err != nil {
		fail(c, 401, "INVALID_ACCESS_TOKEN", "Token akses tidak valid")
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil || id == uuid.Nil {
		fail(c, 400, "INVALID_INPUT", "ID pemeriksaan tidak valid")
		return
	}
	var inspection domain.Inspection
	err = h.db.WithContext(c.Request.Context()).Where("id = ? AND user_id = ?", id, userID).First(&inspection).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fail(c, 404, "INSPECTION_NOT_FOUND", "Pemeriksaan tidak ditemukan")
		return
	}
	if err != nil {
		fail(c, 500, "INTERNAL_ERROR", "Terjadi kesalahan internal")
		return
	}
	parts := strings.Split(inspection.ImageURL, ".")
	if len(parts) != 2 || (parts[1] != "jpeg" && parts[1] != "png") {
		fail(c, 500, "INTERNAL_ERROR", "Foto tidak tersedia")
		return
	}
	if keyID, err := uuid.Parse(parts[0]); err != nil || keyID.String() != parts[0] {
		fail(c, 500, "INTERNAL_ERROR", "Foto tidak tersedia")
		return
	}
	data, err := h.storage.Get(c.Request.Context(), inspection.ImageURL)
	if err != nil {
		fail(c, 503, "STORAGE_UNAVAILABLE", "Foto tidak tersedia")
		return
	}
	c.Header("X-Content-Type-Options", "nosniff")
	c.Data(200, "image/"+parts[1], data)
}
