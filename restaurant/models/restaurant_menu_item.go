package models

import "github.com/google/uuid"

type MenuItem struct {
	OrderItemID  uuid.UUID `json:"order_item_id" db:"order_item_id"`
	RestaurantID uuid.UUID `json:"restaurant_id" db:"restaurant_id"`
	Name         string    `json:"name" db:"name"`
	Price        float64   `json:"price" db:"price"`
	Quantity     int       `json:"quantity" db:"quantity"`
	Description  string    `json:"description" db:"description"`
}
