package inspection

import (
	"bytes"
	"encoding/json"
	"image"
	"image/png"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Project-TETI/AgriSense/backend/internal/domain"
	"github.com/Project-TETI/AgriSense/backend/internal/middleware"
	"github.com/Project-TETI/AgriSense/backend/internal/repository"
	"github.com/Project-TETI/AgriSense/backend/internal/service"
	"github.com/Project-TETI/AgriSense/backend/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func isolatedDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := os.Getenv("INSPECTION_TEST_DATABASE_DSN")
	if dsn == "" {
		t.Skip("set INSPECTION_TEST_DATABASE_DSN to an isolated PostgreSQL test database")
	}
	cfg, err := pgx.ParseConfig(dsn)
	if err != nil {
		t.Fatal(err)
	}
	adminSQL := stdlib.OpenDB(*cfg)
	admin, err := gorm.Open(postgres.New(postgres.Config{Conn: adminSQL}), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		_ = adminSQL.Close()
		t.Fatal(err)
	}
	schema := "inspection_test_" + strings.ReplaceAll(uuid.NewString(), "-", "")
	if err := admin.Exec("CREATE SCHEMA " + schema).Error; err != nil {
		_ = adminSQL.Close()
		t.Fatal(err)
	}
	testCfg := cfg.Copy()
	testCfg.RuntimeParams["search_path"] = schema
	sqlDB := stdlib.OpenDB(*testCfg)
	t.Cleanup(func() {
		_ = sqlDB.Close()
		if err := admin.Exec("DROP SCHEMA " + schema + " CASCADE").Error; err != nil {
			t.Error(err)
		}
		_ = adminSQL.Close()
	})
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&domain.User{}, &domain.RefreshToken{}, &domain.Plant{}, &domain.Disease{}, &domain.Inspection{}); err != nil {
		t.Fatal(err)
	}
	return db
}

func TestInspectionFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := isolatedDB(t)
	user := domain.User{Name: "Test", Email: "test@example.com", PasswordHash: "test"}
	other := domain.User{Name: "Other", Email: "other@example.com", PasswordHash: "test"}
	plant := domain.Plant{Name: "tomato", DisplayName: "Tomat"}
	for _, item := range []any{&user, &other, &plant} {
		if err := db.Create(item).Error; err != nil {
			t.Fatal(err)
		}
	}
	disease := domain.Disease{PlantID: plant.ID, ClassName: "early_blight", DisplayName: "Early", Cause: "x", Description: "x", Recommendation: "x"}
	if err := db.Create(&disease).Error; err != nil {
		t.Fatal(err)
	}
	response := `{"plant":"tomato","class_name":"early_blight","confidence_score":0.9,"status":"SUCCESS"}`
	status := 200
	mlCalls := 0
	mlServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mlCalls++
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/ml/predict" {
			t.Errorf("unexpected ML request")
		}
		if err := r.ParseMultipartForm(6 << 20); err != nil {
			t.Error(err)
			w.WriteHeader(500)
			return
		}
		if r.FormValue("plant") != "tomato" {
			t.Errorf("wrong plant")
		}
		file, header, err := r.FormFile("file")
		if err != nil {
			t.Error(err)
			w.WriteHeader(500)
			return
		}
		_ = file.Close()
		if header.Header.Get("Content-Type") != "image/png" {
			t.Error("wrong image type")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(response))
	}))
	t.Cleanup(mlServer.Close)
	ml, err := NewMLClient(mlServer.URL)
	if err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(t.TempDir(), "photos")
	storage, err := NewStorage(dir, "")
	if err != nil {
		t.Fatal(err)
	}
	h := NewHandler(db, ml, storage, service.NewDiseaseService(repository.NewDiseaseRepository(db)))
	secret := strings.Repeat("x", 32)
	token, err := utils.GenerateAccessToken(user.ID, user.Email, user.Name, secret, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	otherToken, err := utils.GenerateAccessToken(other.ID, other.Email, other.Name, secret, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	router := gin.New()
	router.POST("/api/v1/inspections", middleware.AuthMiddleware(secret), h.Create)
	router.GET("/api/v1/inspections/:id/image", middleware.AuthMiddleware(secret), h.Image)
	var img bytes.Buffer
	if err := png.Encode(&img, image.NewRGBA(image.Rect(0, 0, 2, 2))); err != nil {
		t.Fatal(err)
	}
	request := func(token, plant string, data []byte) *httptest.ResponseRecorder {
		var body bytes.Buffer
		w := multipart.NewWriter(&body)
		_ = w.WriteField("plant", plant)
		file, err := w.CreateFormFile("file", "leaf.png")
		if err != nil {
			t.Fatal(err)
		}
		_, _ = file.Write(data)
		_ = w.Close()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/inspections", &body)
		req.Header.Set("Content-Type", w.FormDataContentType())
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)
		return res
	}
	if got := request("", "tomato", img.Bytes()); got.Code != 401 {
		t.Fatalf("unauthenticated: %d", got.Code)
	}
	for _, tc := range []struct {
		plant  string
		data   []byte
		status int
	}{{"bad-plant", img.Bytes(), 400}, {"rice", img.Bytes(), 400}, {"tomato", []byte("not an image"), 400}, {"tomato", make([]byte, MaxImageSize+1), 413}} {
		before := mlCalls
		if got := request(token, tc.plant, tc.data); got.Code != tc.status || mlCalls != before {
			t.Fatalf("invalid request: %d %s", got.Code, got.Body.String())
		}
	}
	got := request(token, "tomato", img.Bytes())
	if got.Code != 201 {
		t.Fatalf("success: %d %s", got.Code, got.Body.String())
	}
	var result struct {
		ID       uuid.UUID `json:"inspection_id"`
		ImageURL string    `json:"image_url"`
		Catalog  *bool     `json:"catalog_available"`
		Disease  *struct {
			ClassName string `json:"class_name"`
		} `json:"disease_info"`
	}
	if err := json.Unmarshal(got.Body.Bytes(), &result); err != nil || result.ID == uuid.Nil || result.Disease == nil || result.Disease.ClassName != "early_blight" || result.Catalog == nil || !*result.Catalog {
		t.Fatalf("response %s: %v", got.Body.String(), err)
	}
	var record domain.Inspection
	if err := db.First(&record, "id = ?", result.ID).Error; err != nil || record.UserID != user.ID || record.DiseaseID == nil || *record.DiseaseID != disease.ID {
		t.Fatalf("record %+v %v", record, err)
	}
	for _, tc := range []struct {
		token  string
		status int
	}{{token, 200}, {otherToken, 404}, {"", 401}} {
		req := httptest.NewRequest(http.MethodGet, result.ImageURL, nil)
		if tc.token != "" {
			req.Header.Set("Authorization", "Bearer "+tc.token)
		}
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != tc.status || (tc.status == 200 && !bytes.Equal(w.Body.Bytes(), img.Bytes())) {
			t.Fatalf("image status %d", w.Code)
		}
	}
	response = `{"plant":"tomato","class_name":"unknown","confidence_score":0.3,"status":"UNCERTAIN"}`
	got = request(token, "tomato", img.Bytes())
	if got.Code != 201 || !strings.Contains(got.Body.String(), `"status":"UNCERTAIN"`) || !strings.Contains(got.Body.String(), `"disease_info":null`) || !strings.Contains(got.Body.String(), `"message":`) {
		t.Fatalf("uncertain: %d %s", got.Code, got.Body.String())
	}
	var uncertain domain.Inspection
	if err := db.Where("status = ? AND user_id = ?", "UNCERTAIN", user.ID).First(&uncertain).Error; err != nil || uncertain.DiseaseID != nil || uncertain.PredictedClass == nil || *uncertain.PredictedClass != "unknown" {
		t.Fatalf("uncertain not persisted: %+v %v", uncertain, err)
	}
	response = `{"plant":"tomato","class_name":"healthy","confidence_score":0.9,"status":"SUCCESS"}`
	got = request(token, "tomato", img.Bytes())
	if got.Code != 201 || !strings.Contains(got.Body.String(), `"disease_info":null`) {
		t.Fatalf("healthy: %d %s", got.Code, got.Body.String())
	}
	response = `{"plant":"tomato","class_name":"missing_class","confidence_score":0.9,"status":"SUCCESS"}`
	got = request(token, "tomato", img.Bytes())
	if got.Code != 201 || !strings.Contains(got.Body.String(), `"catalog_available":false`) {
		t.Fatalf("missing catalog: %d %s", got.Code, got.Body.String())
	}
	status, response = 422, `{"detail":{"error_code":"BLURRY_IMAGE"}}`
	if got := request(token, "tomato", img.Bytes()); got.Code != 422 || !strings.Contains(got.Body.String(), "BLURRY_IMAGE") {
		t.Fatalf("blur: %d %s", got.Code, got.Body.String())
	}
	status, response = 200, `{"plant":"tomato","class_name":"healthy","confidence_score":4,"status":"SUCCESS"}`
	if got := request(token, "tomato", img.Bytes()); got.Code != 502 {
		t.Fatalf("bad ML: %d %s", got.Code, got.Body.String())
	}
	status, response = 503, `{"detail":"secret"}`
	if got := request(token, "tomato", img.Bytes()); got.Code != 503 || strings.Contains(got.Body.String(), "secret") {
		t.Fatalf("ML unavailable: %d %s", got.Code, got.Body.String())
	}
	status, response = 200, `{"plant":"tomato","class_name":"healthy","confidence_score":0.9,"status":"SUCCESS"}`
	storage.dir = filepath.Join(dir, "missing")
	if got := request(token, "tomato", img.Bytes()); got.Code != 503 {
		t.Fatalf("storage failure: %d %s", got.Code, got.Body.String())
	}
	storage.dir = dir
	if err := db.Exec("ALTER TABLE inspections ADD CONSTRAINT reject_inserts CHECK (status <> 'SUCCESS') NOT VALID").Error; err != nil {
		t.Fatal(err)
	}
	before, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		t.Fatal(err)
	}
	if got := request(token, "tomato", img.Bytes()); got.Code != 500 {
		t.Fatalf("database failure: %d %s", got.Code, got.Body.String())
	}
	after, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil || len(after) != len(before) {
		t.Fatalf("cleanup: %v %v", after, err)
	}
}
