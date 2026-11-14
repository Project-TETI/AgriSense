package service_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/Project-TETI/AgriSense/backend/internal/domain"
	"github.com/Project-TETI/AgriSense/backend/internal/repository"
	"github.com/Project-TETI/AgriSense/backend/internal/service"
)

func TestCatalogMigrationAndPlantClassLookup(t *testing.T) {
	db := newAuthTestDB(t).db
	if err := db.Exec("DROP INDEX idx_diseases_plant_class; CREATE UNIQUE INDEX idx_diseases_class_name ON diseases (class_name); ALTER TABLE diseases DROP COLUMN symptoms").Error; err != nil {
		t.Fatal(err)
	}
	plantA := domain.Plant{Name: "plant_a", DisplayName: "Plant A"}
	plantB := domain.Plant{Name: "plant_b", DisplayName: "Plant B"}
	if err := db.Create(&plantA).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&plantB).Error; err != nil {
		t.Fatal(err)
	}
	first := domain.Disease{PlantID: plantA.ID, ClassName: "shared", DisplayName: "A", Cause: "a", Description: "a", Recommendation: "a"}
	if err := db.Omit("Symptoms").Create(&first).Error; err != nil {
		t.Fatal(err)
	}
	script, err := os.ReadFile("../../migrations/002_disease_plant_class.sql")
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		if err := db.Exec(string(script)).Error; err != nil {
			t.Fatal(err)
		}
	}
	second := domain.Disease{PlantID: plantB.ID, ClassName: "shared", DisplayName: "B", Cause: "b", Description: "b", Recommendation: "b"}
	if err := db.Create(&second).Error; err != nil {
		t.Fatal(err)
	}
	duplicate := domain.Disease{PlantID: plantA.ID, ClassName: "shared", DisplayName: "duplicate", Cause: "x", Description: "x", Recommendation: "x"}
	if err := db.Create(&duplicate).Error; err == nil {
		t.Fatal("duplicate plant and class accepted")
	}
	svc := service.NewDiseaseService(repository.NewDiseaseRepository(db))
	for _, tc := range []struct{ plant, want string }{{"plant_a", "A"}, {"plant_b", "B"}} {
		got, err := svc.GetDiseaseByPlantAndClass(context.Background(), tc.plant, "shared")
		if err != nil || got.DisplayName != tc.want || got.Plant.Name != tc.plant {
			t.Fatalf("lookup %s: %+v, %v", tc.plant, got, err)
		}
	}
	if _, err := svc.GetDiseaseByPlantAndClass(context.Background(), "missing", "shared"); !errors.Is(err, service.ErrDiseaseNotFound) {
		t.Fatalf("missing: %v", err)
	}
}

func TestCatalogListsAreOrderedAndFiltered(t *testing.T) {
	db := newAuthTestDB(t).db
	a := domain.Plant{Name: "a", DisplayName: "A"}
	b := domain.Plant{Name: "b", DisplayName: "B"}
	if err := db.Create(&b).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&a).Error; err != nil {
		t.Fatal(err)
	}
	plants, err := repository.NewPlantRepository(db).FindAll(context.Background())
	if err != nil || len(plants) != 2 || plants[0].Name != "a" || plants[1].Name != "b" {
		t.Fatalf("plants=%+v err=%v", plants, err)
	}
	for _, d := range []domain.Disease{
		{PlantID: a.ID, ClassName: "z", DisplayName: "Z", Cause: "x", Description: "x", Recommendation: "x"},
		{PlantID: b.ID, ClassName: "a", DisplayName: "A", Cause: "x", Description: "x", Recommendation: "x"},
		{PlantID: a.ID, ClassName: "a", DisplayName: "A", Cause: "x", Description: "x", Recommendation: "x"},
	} {
		if err := db.Create(&d).Error; err != nil {
			t.Fatal(err)
		}
	}
	repo := repository.NewDiseaseRepository(db)
	all, err := repo.FindAll(context.Background(), "")
	if err != nil || len(all) != 3 || all[0].ClassName != "a" || all[1].ClassName != "a" || all[0].ID.String() > all[1].ID.String() || all[2].ClassName != "z" {
		t.Fatalf("all=%+v err=%v", all, err)
	}
	filtered, err := repo.FindAll(context.Background(), a.ID.String())
	if err != nil || len(filtered) != 2 || filtered[0].ClassName != "a" || filtered[1].ClassName != "z" {
		t.Fatalf("filtered=%+v err=%v", filtered, err)
	}
}
