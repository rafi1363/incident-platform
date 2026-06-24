package services

import (
	"context"
	"errors"

	"incident-platform/backend/internal/models"
	"incident-platform/backend/internal/repository"

	"github.com/jackc/pgx/v5"
)

var validIncidentStatuses = map[string]bool{
	"OPEN":       true,
	"RESOLVED":   true,
	"MONITORING": true,
}

type IncidentService struct {
	repo *repository.IncidentRepository
}

func NewIncidentService(repo *repository.IncidentRepository) *IncidentService {
	return &IncidentService{repo: repo}
}

func (s *IncidentService) Create(ctx context.Context, in models.CreateIncidentInput) (models.Incident, error) {
	if !validIncidentStatuses[in.Status] {
		return models.Incident{}, errors.New("invalid incident status")
	}
	return s.repo.Create(ctx, in)
}

func (s *IncidentService) GetAll(ctx context.Context) ([]models.Incident, error) {
	return s.repo.GetAll(ctx)
}

func (s *IncidentService) GetByID(ctx context.Context, id int) (models.Incident, error) {
	inc, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Incident{}, ErrNotFound
		}
		return models.Incident{}, err
	}
	return inc, nil
}

func (s *IncidentService) GetByServiceID(ctx context.Context, serviceID int) ([]models.Incident, error) {
	return s.repo.GetByServiceID(ctx, serviceID)
}

func (s *IncidentService) Update(ctx context.Context, id int, in models.UpdateIncidentInput) (models.Incident, error) {
	if in.Status != nil && !validIncidentStatuses[*in.Status] {
		return models.Incident{}, errors.New("invalid incident status")
	}
	inc, err := s.repo.Update(ctx, id, in)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Incident{}, ErrNotFound
		}
		return models.Incident{}, err
	}
	return inc, nil
}

func (s *IncidentService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *IncidentService) GetOpenByServiceID(ctx context.Context, serviceID int) (models.Incident, error) {
	inc, err := s.repo.GetOpenByServiceID(ctx, serviceID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Incident{}, ErrNotFound
		}
		return models.Incident{}, err
	}
	return inc, nil
}

func (s *IncidentService) Resolve(ctx context.Context, id int) error {
	return s.repo.Resolve(ctx, id)
}
