package services

import (
	"context"
	"errors"

	"incident-platform/backend/internal/models"
	"incident-platform/backend/internal/repository"

	"github.com/jackc/pgx/v5"
)

var ErrNotFound = errors.New("not found")

var validStatusses = map[string]bool{
	"UP":      true,
	"DOWN":    true,
	"UNKNOWN": true,
}

type ServiceService struct {
	repo *repository.ServiceRepository
}

func NewServiceService(repo *repository.ServiceRepository) *ServiceService {
	return &ServiceService{repo: repo}
}

func (s *ServiceService) GetAll(ctx context.Context) ([]models.Service, error) {
	return s.repo.GetAll(ctx)
}

func (s *ServiceService) Create(ctx context.Context, in models.CreateServiceInput) (models.Service, error) {
	return s.repo.Create(ctx, in.Name, in.URL)
}

func (s *ServiceService) GetByID(ctx context.Context, id int) (models.Service, error) {
	svc, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Service{}, ErrNotFound
		}
		return models.Service{}, err
	}
	return svc, nil
}

func (s *ServiceService) Update(ctx context.Context, id int, in models.UpdateServiceInput) (models.Service, error) {
	if in.Status != nil && !validStatusses[*in.Status] {
		return models.Service{}, errors.New("Invalid Status")
	}

	svc, err := s.repo.Update(ctx, id, in)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Service{}, ErrNotFound
		}
		return models.Service{}, err
	}
	return svc, nil
}

func (s *ServiceService) Delete(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)
	return err
}
