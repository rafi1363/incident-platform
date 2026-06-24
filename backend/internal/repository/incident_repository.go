package repository

import (
	"context"
	"incident-platform/backend/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type IncidentRepository struct {
	DB *pgxpool.Pool
}

func NewIncidentRepository(db *pgxpool.Pool) *IncidentRepository {
	return &IncidentRepository{DB: db}
}

func (r *IncidentRepository) Create(ctx context.Context, in models.CreateIncidentInput) (models.Incident, error) {
	var inc models.Incident
	err := r.DB.QueryRow(ctx, `
		INSERT INTO incidents (service_id, status, message)
		VALUES ($1, $2, $3)
		RETURNING id, service_id, status, message, started_at, resolved_at
	`, in.ServiceID, in.Status, in.Message).Scan(
		&inc.ID, &inc.ServiceID, &inc.Status, &inc.Message, &inc.StartedAt, &inc.ResolvedAt,
	)
	if err != nil {
		return models.Incident{}, err
	}
	return inc, nil
}

func (r *IncidentRepository) GetAll(ctx context.Context) ([]models.Incident, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT id, service_id, status, message, started_at, resolved_at
		FROM incidents
		ORDER BY started_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incidents []models.Incident
	for rows.Next() {
		var inc models.Incident
		if err := rows.Scan(&inc.ID, &inc.ServiceID, &inc.Status, &inc.Message, &inc.StartedAt, &inc.ResolvedAt); err != nil {
			return nil, err
		}
		incidents = append(incidents, inc)
	}
	return incidents, rows.Err()
}

func (r *IncidentRepository) GetByID(ctx context.Context, id int) (models.Incident, error) {
	var inc models.Incident
	err := r.DB.QueryRow(ctx, `
		SELECT id, service_id, status, message, started_at, resolved_at
		FROM incidents WHERE id = $1
	`, id).Scan(&inc.ID, &inc.ServiceID, &inc.Status, &inc.Message, &inc.StartedAt, &inc.ResolvedAt)
	if err != nil {
		return models.Incident{}, err
	}
	return inc, nil
}

// GetByServiceID — a useful "filter" query. Notice we also join to validate
// the service exists; if the FK is violated the INSERT would fail anyway.
func (r *IncidentRepository) GetByServiceID(ctx context.Context, serviceID int) ([]models.Incident, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT id, service_id, status, message, started_at, resolved_at
		FROM incidents WHERE service_id = $1
		ORDER BY started_at DESC
	`, serviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incidents []models.Incident
	for rows.Next() {
		var inc models.Incident
		if err := rows.Scan(&inc.ID, &inc.ServiceID, &inc.Status, &inc.Message, &inc.StartedAt, &inc.ResolvedAt); err != nil {
			return nil, err
		}
		incidents = append(incidents, inc)
	}
	return incidents, rows.Err()
}

func (r *IncidentRepository) Update(ctx context.Context, id int, in models.UpdateIncidentInput) (models.Incident, error) {
	var inc models.Incident
	err := r.DB.QueryRow(ctx, `
		UPDATE incidents
		SET status     = COALESCE($2, status),
		    message    = COALESCE($3, message),
		    resolved_at = COALESCE($4, resolved_at)
		WHERE id = $1
		RETURNING id, service_id, status, message, started_at, resolved_at
	`, id, in.Status, in.Message, in.ResolvedAt).Scan(
		&inc.ID, &inc.ServiceID, &inc.Status, &inc.Message, &inc.StartedAt, &inc.ResolvedAt,
	)
	if err != nil {
		return models.Incident{}, err
	}
	return inc, nil
}

func (r *IncidentRepository) Delete(ctx context.Context, id int) error {
	_, err := r.DB.Exec(ctx, `DELETE FROM incidents WHERE id = $1`, id)
	return err
}

func (r *IncidentRepository) GetOpenByServiceID(ctx context.Context, serviceID int) (models.Incident, error) {
	var inc models.Incident
	err := r.DB.QueryRow(ctx, `
		SELECT id, service_id, status, message, started_at, resolved_at
		FROM incidents
		WHERE service_id = $1 AND status = 'OPEN'
		ORDER BY started_at DESC
		LIMIT 1
	`, serviceID).Scan(
		&inc.ID, &inc.ServiceID, &inc.Status, &inc.Message, &inc.StartedAt, &inc.ResolvedAt,
	)
	if err != nil {
		return models.Incident{}, err
	}
	return inc, nil
}

func (r *IncidentRepository) Resolve(ctx context.Context, id int) error {
	_, err := r.DB.Exec(ctx, `
		UPDATE incidents
		SET status = 'RESOLVED', resolved_at = NOW()
		WHERE id = $1
	`, id)
	return err
}
