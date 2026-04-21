package domain

import (
	"context"

	"github.com/google/uuid"
)

type Courier struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	TransportType string    `json:"transport_type"`
	IsActive      bool      `json:"is_active"`
}

type CourierRepo interface {
	ListCouriers(ctx context.Context) ([]Courier, error)
}
