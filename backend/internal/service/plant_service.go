package service

import (
	"context"

	"github.com/Project-TETI/AgriSense/backend/internal/domain"
	"github.com/Project-TETI/AgriSense/backend/internal/repository"
)

type PlantService interface {
	GetAllPlants(ctx context.Context) ([]domain.Plant, error)
}

type plantService struct {
	plantRepo repository.PlantRepository
}

func NewPlantService(plantRepo repository.PlantRepository) PlantService {
	return &plantService{plantRepo: plantRepo}
}

func (s *plantService) GetAllPlants(ctx context.Context) ([]domain.Plant, error) {
	return s.plantRepo.FindAll(ctx)
}
