package models

import "time"

type Incident struct {
	ID         int        `json:"id"`
	ServiceID  int        `json:"service_id"`
	Status     string     `json:"status"` // e.g. "OPEN", "RESOLVED"
	Message    string     `json:"message"`
	StartedAt  time.Time  `json:"started_at"`
	ResolvedAt *time.Time `json:"resolved_at"` // pointer = nullable in DB
}

type CreateIncidentInput struct {
	ServiceID int    `json:"service_id" binding:"required"`
	Status    string `json:"status"      binding:"required"`
	Message   string `json:"message"`
}

type UpdateIncidentInput struct {
	Status     *string    `json:"status"`
	Message    *string    `json:"message"`
	ResolvedAt *time.Time `json:"resolved_at"`
}
