package catalogseed

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Project-TETI/AgriSense/backend/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Disease struct {
	ClassName      string  `json:"class_name"`
	DisplayName    string  `json:"display_name"`
	Cause          string  `json:"cause"`
	Description    string  `json:"description"`
	Symptoms       *string `json:"symptoms"`
	Recommendation string  `json:"recommendation"`
}

type Plant struct {
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Description *string   `json:"description"`
	Diseases    []Disease `json:"diseases"`
}

var identifier = regexp.MustCompile(`^[a-z0-9]+(_[a-z0-9]+)*$`)

func Load(dir string) ([]Plant, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.json"))
	if err != nil || len(files) == 0 {
		return nil, fmt.Errorf("no catalog JSON files found in %s", dir)
	}
	plants := make([]Plant, 0, len(files))
	seen := make(map[string]bool)
	for _, name := range files {
		file, err := os.Open(name)
		if err != nil {
			return nil, fmt.Errorf("open %s: %w", name, err)
		}
		var plant Plant
		decoder := json.NewDecoder(file)
		decoder.DisallowUnknownFields()
		err = decoder.Decode(&plant)
		if err == nil {
			var extra any
			if err = decoder.Decode(&extra); err == io.EOF {
				err = nil
			} else if err == nil {
				err = fmt.Errorf("multiple JSON values")
			}
		}
		closeErr := file.Close()
		if err != nil {
			return nil, fmt.Errorf("decode %s: %w", name, err)
		}
		if closeErr != nil {
			return nil, fmt.Errorf("close %s: %w", name, closeErr)
		}
		if len(plant.Name) > 100 || !identifier.MatchString(plant.Name) || !required(plant.DisplayName, 100) || seen[plant.Name] {
			return nil, fmt.Errorf("invalid or duplicate plant in %s", name)
		}
		seen[plant.Name] = true
		classes := make(map[string]bool)
		for _, disease := range plant.Diseases {
			if len(disease.ClassName) > 150 || !identifier.MatchString(disease.ClassName) || classes[disease.ClassName] ||
				!required(disease.DisplayName, 150) || !required(disease.Cause, 0) || !required(disease.Description, 0) || !required(disease.Recommendation, 0) {
				return nil, fmt.Errorf("invalid or duplicate disease in %s", name)
			}
			classes[disease.ClassName] = true
		}
		plants = append(plants, plant)
	}
	return plants, nil
}

func required(value string, max int) bool {
	return strings.TrimSpace(value) != "" && (max == 0 || len([]rune(value)) <= max)
}

type Counts struct{ PlantsCreated, PlantsSkipped, DiseasesCreated, DiseasesSkipped int64 }

func Apply(db *gorm.DB, plants []Plant) (Counts, error) {
	var counts Counts
	err := db.Transaction(func(tx *gorm.DB) error {
		for _, input := range plants {
			plant := domain.Plant{Name: input.Name, DisplayName: input.DisplayName, Description: input.Description}
			result := tx.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "name"}}, DoNothing: true}).Omit(clause.Associations).Create(&plant)
			if result.Error != nil {
				return result.Error
			}
			counts.PlantsCreated += result.RowsAffected
			counts.PlantsSkipped += 1 - result.RowsAffected
			if result.RowsAffected == 0 {
				plant = domain.Plant{}
				if err := tx.Where("name = ?", input.Name).First(&plant).Error; err != nil {
					return err
				}
			}
			for _, item := range input.Diseases {
				disease := domain.Disease{PlantID: plant.ID, ClassName: item.ClassName, DisplayName: item.DisplayName, Cause: item.Cause, Description: item.Description, Symptoms: item.Symptoms, Recommendation: item.Recommendation}
				result := tx.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "plant_id"}, {Name: "class_name"}}, DoNothing: true}).Omit(clause.Associations).Create(&disease)
				if result.Error != nil {
					return result.Error
				}
				counts.DiseasesCreated += result.RowsAffected
				counts.DiseasesSkipped += 1 - result.RowsAffected
			}
		}
		return nil
	})
	if err != nil {
		return Counts{}, err
	}
	return counts, nil
}
