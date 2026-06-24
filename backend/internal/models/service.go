package models

import "time"

type Service struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateServiceInput struct {
	Name string `json:"name" binding:"required"`
	URL  string `json:"url" binding:"required"`
}

type UpdateServiceInput struct {
	Name   *string `json:"name"`
	URL    *string `jsong:"url"`
	Status *string `json:"status"`
}
