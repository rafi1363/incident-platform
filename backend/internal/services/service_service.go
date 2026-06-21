package services

import (
	"context"

	"incident-platform/backend/internal/models"
	"incident-platform/backend/internal/repository"
)

type ServiceService struct {
	repo *repository.ServiceRepository
}

func NewServiceService(repo *repository.ServiceRepository) *ServiceService {
	return &ServiceService{repo: repo}
}

func (s *ServiceService) GetAll(ctx context.Context) ([]models.Service, error) {
	services, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	for i := range services {
		if services[i].Status == "UP" {
			services[i].Status = "HEALTHY"
		}
	}
	return services, nil
}
