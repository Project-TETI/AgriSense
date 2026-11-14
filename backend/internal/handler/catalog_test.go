package handler

import (
	"context"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Project-TETI/AgriSense/backend/internal/domain"
	"github.com/Project-TETI/AgriSense/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type catalogStub struct {
	plant, class string
	err          error
	listErr      error
	filter       string
	calls        int
}

func (s *catalogStub) GetAllDiseases(_ context.Context, filter string) ([]domain.Disease, error) {
	s.calls++
	s.filter = filter
	return nil, s.listErr
}
func (s *catalogStub) GetDiseaseByPlantAndClass(_ context.Context, plant, class string) (*domain.Disease, error) {
	s.calls++
	s.plant, s.class = plant, class
	if s.err != nil {
		return nil, s.err
	}
	return &domain.Disease{ClassName: class, Plant: domain.Plant{Name: plant}}, nil
}

func TestCatalogDetail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	stub := &catalogStub{}
	r := gin.New()
	r.GET("/plants/:plant_name/diseases/:class_name", NewDiseaseHandler(stub).GetDiseaseByPlantAndClass)
	for _, tc := range []struct {
		err    error
		status int
		body   string
	}{
		{nil, 200, `"name":"tomato"`},
		{service.ErrDiseaseNotFound, 404, "DISEASE_NOT_FOUND"},
		{errors.New("private database detail"), 500, "INTERNAL_ERROR"},
	} {
		stub.err = tc.err
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest("GET", "/plants/tomato/diseases/spot", nil))
		if w.Code != tc.status || !strings.Contains(w.Body.String(), tc.body) || strings.Contains(w.Body.String(), "private database detail") || stub.plant != "tomato" || stub.class != "spot" {
			t.Fatalf("status %d body %s plant %s class %s", w.Code, w.Body.String(), stub.plant, stub.class)
		}
	}
}

type plantStub struct{ err error }

func (s plantStub) GetAllPlants(context.Context) ([]domain.Plant, error) { return nil, s.err }

func TestCatalogValidationAndLists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	stub := &catalogStub{}
	r := gin.New()
	r.GET("/diseases", NewDiseaseHandler(stub).GetDiseases)
	r.GET("/plants/:plant_name/diseases/:class_name", NewDiseaseHandler(stub).GetDiseaseByPlantAndClass)
	r.GET("/plants", NewPlantHandler(plantStub{}).GetPlants)
	r.GET("/plant-error", NewPlantHandler(plantStub{err: errors.New("private database detail")}).GetPlants)
	validID := "5b0548c2-8ad7-4301-94ee-b476609a154e"
	for _, tc := range []struct {
		path   string
		status int
		body   string
		called bool
	}{
		{"/plants", 200, `"plants":[]`, false},
		{"/plant-error", 500, `"code":"INTERNAL_ERROR"`, false},
		{"/diseases", 200, `"diseases":[]`, true},
		{"/diseases?plant_id=" + validID, 200, `"diseases":[]`, true},
		{"/diseases?plant_id=", 400, `"code":"INVALID_INPUT"`, false},
		{"/diseases?plant_id=bad", 400, `"code":"INVALID_INPUT"`, false},
		{"/diseases?plant_id=00000000-0000-0000-0000-000000000000", 400, `"code":"INVALID_INPUT"`, false},
		{"/diseases?plant_id=" + validID + "&plant_id=" + validID, 400, `"code":"INVALID_INPUT"`, false},
		{"/plants/Tomato/diseases/spot", 400, `"code":"INVALID_INPUT"`, false},
		{"/plants/tomato/diseases/early__blight", 400, `"code":"INVALID_INPUT"`, false},
		{"/plants/" + strings.Repeat("a", 101) + "/diseases/spot", 400, `"code":"INVALID_INPUT"`, false},
	} {
		t.Run(tc.path, func(t *testing.T) {
			before := stub.calls
			w := httptest.NewRecorder()
			r.ServeHTTP(w, httptest.NewRequest("GET", tc.path, nil))
			if w.Code != tc.status || !strings.Contains(w.Body.String(), tc.body) || strings.Contains(w.Body.String(), "private database detail") || (stub.calls > before) != tc.called {
				t.Fatalf("status=%d body=%s calls=%d", w.Code, w.Body.String(), stub.calls-before)
			}
			if strings.Contains(tc.path, validID) && tc.status == 200 && stub.filter != validID {
				t.Fatalf("filter=%s", stub.filter)
			}
		})
	}
	stub.listErr = errors.New("private database detail")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/diseases", nil))
	if w.Code != 500 || !strings.Contains(w.Body.String(), "INTERNAL_ERROR") || strings.Contains(w.Body.String(), "private database detail") {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
}
