package models

import "github.com/google/uuid"

type Restaurant struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Address string    `json:"address"`
	Status  bool      `json:"status"`
}

type RestaurantMenuItem struct {
	ID           uuid.UUID `json:"id"`
	RestaurantID uuid.UUID `json:"restaurant_id"`
	Name         string    `json:"name"`
	Price        float64   `json:"price"`
}
