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

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return services, nil
}

func (r *ServiceRepository) Create(ctx context.Context, name, url string) (models.Service, error) {
	var s models.Service

	err := r.DB.QueryRow(ctx, `
		INSERT INTO services (name, url)
		VALUES ($1, $2)
		RETURNING id, name, url, status, created_at, updated_at
		`, name, url).Scan(
		&s.ID, &s.Name, &s.URL, &s.Status, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return models.Service{}, err
	}
	return s, nil
}

func (r *ServiceRepository) GetByID(ctx context.Context, id int) (models.Service, error) {
	var s models.Service

	err := r.DB.QueryRow(ctx, `
		SELECT id, name, url, status, created_at, updated_at
		FROM services
		WHERE id = $1
		`, id).Scan(
		&s.ID, &s.Name, &s.URL, &s.Status, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return models.Service{}, err
	}
	return s, nil
}
func (r *ServiceRepository) Update(ctx context.Context, id int, in models.UpdateServiceInput) (models.Service, error) {
	var s models.Service

	err := r.DB.QueryRow(ctx, `
		UPDATE services
		SET name    = COALESCE($2, name),
		    url     = COALESCE($3, url),
		    status  = COALESCE($4, status),
		    updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, url, status, created_at, updated_at
		`, id, in.Name, in.URL, in.Status).Scan(
		&s.ID, &s.Name, &s.URL, &s.Status, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return models.Service{}, err
	}
	return s, nil
}

func (r *ServiceRepository) Delete(ctx context.Context, id int) error {
	_, err := r.DB.Exec(ctx, `DELETE FROM services WHERE id = $1`, id)
	return err
}
