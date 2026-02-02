package models

import (
	"github.com/Kabanya/YAFDS/pkg/models"

	"github.com/google/uuid"
)

type OrderItemInput struct {
	RestaurantItemID uuid.UUID
	Price            float64
	Quantity         int
}

type Filter struct {
	ID         *uuid.UUID
	CustomerID *uuid.UUID
	CourierID  *uuid.UUID
	Status     string
}

type AcceptInput struct {
	OrderID    uuid.UUID
	CustomerID uuid.UUID
	CourierID  uuid.UUID
	Items      []OrderItemInput
	Status     models.OrderStatus
}

type AcceptResult struct {
	OrderID uuid.UUID `json:"order_id"`
	Status  string    `json:"status"`
}
