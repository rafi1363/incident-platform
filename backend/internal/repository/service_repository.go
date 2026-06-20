package repository

import (
	"context"
	"incident-platform/backend/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ServiceRepository struct {
	DB *pgxpool.Pool
}

func NewServiceRepository(db *pgxpool.Pool) *ServiceRepository {
	return &ServiceRepository{
		DB: db,
	}
}

func (r *ServiceRepository) GetAll(ctx context.Context) ([]models.Service, error) {
	rows, err := r.DB.Query(
		ctx,
		`
		SELECT
			id,
			name,
			url,
			status,
			created_at,
			updated_at
		FROM services
		ORDER BY id
		`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []models.Service

	for rows.Next() {
		var s models.Service

		err := rows.Scan(
			&s.ID,
			&s.Name,
			&s.URL,
			&s.Status,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		services = append(services, s)
	}
	return services, nil
}
