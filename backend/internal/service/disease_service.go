package service

import (
	"context"
	"errors"

	"github.com/Project-TETI/AgriSense/backend/internal/domain"
	"github.com/Project-TETI/AgriSense/backend/internal/repository"
)

type DiseaseService interface {
	GetAllDiseases(ctx context.Context, plantID string) ([]domain.Disease, error)
	GetDiseaseByClassName(ctx context.Context, className string) (*domain.Disease, error)
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

func (s *diseaseService) GetDiseaseByClassName(ctx context.Context, className string) (*domain.Disease, error) {
	disease, err := s.diseaseRepo.FindByClassName(ctx, className)
	if err != nil {
		return nil, err
	}
	if disease == nil {
		return nil, errors.New("Disease not found")
	}
	return disease, nil
}
