package repository

import (
	"context"

	"github.com/Project-TETI/AgriSense/backend/internal/domain"
	"gorm.io/gorm"
)

type PlantRepository interface {
	FindAll(ctx context.Context) ([]domain.Plant, error)
}

type plantRepository struct {
	db *gorm.DB
}

func NewPlantRepository(db *gorm.DB) PlantRepository {
	return &plantRepository{db: db}
}

func (r *plantRepository) FindAll(ctx context.Context) ([]domain.Plant, error) {
	var plants []domain.Plant
	err := r.db.WithContext(ctx).Find(&plants).Error
	if err != nil {
		return nil, err
	}
	return plants, nil
}
