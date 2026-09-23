package repository

import (
	"context"
	"errors"

	"github.com/Project-TETI/AgriSense/backend/internal/domain"
	"gorm.io/gorm"
)

type DiseaseRepository interface {
	FindAll(ctx context.Context, plantID string) ([]domain.Disease, error)
	FindByClassName(ctx context.Context, className string) (*domain.Disease, error)
}

type diseaseRepository struct {
	db *gorm.DB
}

func NewDiseaseRepository(db *gorm.DB) DiseaseRepository {
	return &diseaseRepository{db: db}
}

func (r *diseaseRepository) FindAll(ctx context.Context, plantID string) ([]domain.Disease, error) {
	var diseases []domain.Disease
	query := r.db.WithContext(ctx).Preload("Plant")
	if plantID != "" {
		query = query.Where("plant_id = ?", plantID)
	}
	err := query.Find(&diseases).Error
	return diseases, err
}

func (r *diseaseRepository) FindByClassName(ctx context.Context, className string) (*domain.Disease, error) {
	var disease domain.Disease
	err := r.db.WithContext(ctx).Where("class_name = ?", className).First(&disease).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &disease, nil
}
