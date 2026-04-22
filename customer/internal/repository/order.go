package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type OrderDTO struct {
	ID         uuid.UUID
	CustomerID uuid.UUID
	CourierID  uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Status     string
}

type OrderItemDTO struct {
	RestaurantItemID uuid.UUID
	Price            float64
	Quantity         int
}

type CustomerOrderRepository interface {
	CreateOrder(ctx context.Context, order OrderDTO) (OrderDTO, error)
	CreateOrderWithItems(ctx context.Context, order OrderDTO, items []OrderItemDTO) (OrderDTO, error)
	GetCustomerWalletAddress(ctx context.Context, customerID uuid.UUID) (string, error)
	UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status string) error
}
