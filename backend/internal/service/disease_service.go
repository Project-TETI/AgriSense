package service

import (
	"context"
	"errors"

	"github.com/Project-TETI/AgriSense/backend/internal/domain"
	"github.com/Project-TETI/AgriSense/backend/internal/repository"
)

type DiseaseService interface {
	GetAllDiseases(ctx context.Context, plantID string) ([]domain.Disease, error)
	GetDiseaseByPlantAndClass(ctx context.Context, plantName, className string) (*domain.Disease, error)
}

type diseaseService struct {
	diseaseRepo repository.DiseaseRepository
}

func NewDiseaseService(diseaseRepo repository.DiseaseRepository) DiseaseService {
	return &diseaseService{diseaseRepo: diseaseRepo}
}

func (s *diseaseService) GetAllDiseases(ctx context.Context, plantID string) ([]domain.Disease, error) {
	return s.diseaseRepo.FindAll(ctx, plantID)
}

var ErrDiseaseNotFound = errors.New("disease not found")

func (s *diseaseService) GetDiseaseByPlantAndClass(ctx context.Context, plantName, className string) (*domain.Disease, error) {
	disease, err := s.diseaseRepo.FindByPlantAndClass(ctx, plantName, className)
	if err != nil {
		return nil, err
	}
	if disease == nil {
		return nil, ErrDiseaseNotFound
	}
	return disease, nil
}
