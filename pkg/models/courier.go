package models

import "github.com/google/uuid"

type Courier struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	TransportType string    `json:"transport_type"`
	IsActive      bool      `json:"is_active"`
}
