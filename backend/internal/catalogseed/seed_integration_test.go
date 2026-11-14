package catalogseed

import (
	"os"
	"strings"
	"testing"

	"github.com/Project-TETI/AgriSense/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func testDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := os.Getenv("CATALOG_TEST_DATABASE_DSN")
	if dsn == "" {
		t.Skip("set CATALOG_TEST_DATABASE_DSN to an isolated PostgreSQL test database")
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
	schema := "catalog_test_" + strings.ReplaceAll(uuid.NewString(), "-", "")
	if err := admin.Exec("CREATE SCHEMA " + schema).Error; err != nil {
		_ = adminSQL.Close()
		t.Fatal(err)
	}
	testCfg := cfg.Copy()
	testCfg.RuntimeParams["search_path"] = schema
	testSQL := stdlib.OpenDB(*testCfg)
	t.Cleanup(func() {
		_ = testSQL.Close()
		if err := admin.Exec("DROP SCHEMA " + schema + " CASCADE").Error; err != nil {
			t.Error(err)
		}
		_ = adminSQL.Close()
	})
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: testSQL}), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&domain.Plant{}, &domain.Disease{}); err != nil {
		t.Fatal(err)
	}
	return db
}

func TestApplyPreservesExistingAndRollsBack(t *testing.T) {
	db := testDB(t)
	input := []Plant{{Name: "tomato", DisplayName: "Tomat", Diseases: []Disease{{ClassName: "spot", DisplayName: "Bercak", Cause: "x", Description: "x", Recommendation: "x"}}}}
	counts, err := Apply(db, input)
	if err != nil || counts != (Counts{PlantsCreated: 1, DiseasesCreated: 1}) {
		t.Fatalf("first: %+v %v", counts, err)
	}
	var before domain.Disease
	if err := db.First(&before).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Model(&before).Update("display_name", "manual").Error; err != nil {
		t.Fatal(err)
	}
	input[0].DisplayName = "Changed"
	input[0].Diseases[0].DisplayName = "Changed"
	counts, err = Apply(db, input)
	if err != nil || counts != (Counts{PlantsSkipped: 1, DiseasesSkipped: 1}) {
		t.Fatalf("second: %+v %v", counts, err)
	}
	var after domain.Disease
	var plant domain.Plant
	if err := db.First(&after).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.First(&plant).Error; err != nil {
		t.Fatal(err)
	}
	if after.ID != before.ID || after.PlantID != plant.ID || after.DisplayName != "manual" || plant.DisplayName != "Tomat" {
		t.Fatalf("data overwritten: %+v %+v", plant, after)
	}
	input[0].Diseases = append(input[0].Diseases, Disease{ClassName: "new_spot", DisplayName: "New", Cause: "x", Description: "x", Recommendation: "x"})
	input = append(input, Plant{Name: "corn", DisplayName: "Corn", Diseases: []Disease{{ClassName: "spot", DisplayName: "Corn spot", Cause: "x", Description: "x", Recommendation: "x"}}})
	counts, err = Apply(db, input)
	if err != nil || counts != (Counts{PlantsCreated: 1, PlantsSkipped: 1, DiseasesCreated: 2, DiseasesSkipped: 1}) {
		t.Fatalf("additional data: %+v %v", counts, err)
	}
	if err := db.Exec("ALTER TABLE diseases ADD CONSTRAINT reject_blocked CHECK (class_name <> 'blocked')").Error; err != nil {
		t.Fatal(err)
	}
	input = []Plant{{Name: "rice", DisplayName: "Rice", Diseases: []Disease{{ClassName: "blocked", DisplayName: "Blocked", Cause: "x", Description: "x", Recommendation: "x"}}}}
	if _, err := Apply(db, input); err == nil {
		t.Fatal("expected insert failure")
	}
	var count int64
	if err := db.Model(&domain.Plant{}).Where("name = ?", "rice").Count(&count).Error; err != nil || count != 0 {
		t.Fatalf("rollback: count=%d error=%v", count, err)
	}
}
