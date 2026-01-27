package models

import (
	"context"

	"github.com/Kabanya/YAFDS/pkg/models"

	"github.com/google/uuid"
)

type OrderRepo interface { //dto - data transfer object. dto похож на model но другое. Если совпадают, то приоритет модели
	//                   по возможности не плодим dto
	CreateOrder(ctx context.Context, order models.Order) (models.Order, error)
	CreateOrderWithItems(ctx context.Context, order models.Order, items []OrderItemInput) (models.Order, error)
	ListOrders(ctx context.Context, filter Filter) ([]models.Order, error)

	GetOrder(ctx context.Context, orderID uuid.UUID) (models.Order, error)
	AcceptOrder(ctx context.Context, input AcceptInput) (AcceptResult, error)

	GetOrderStatus(ctx context.Context, orderID uuid.UUID) (models.OrderStatus, error)
	UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status models.OrderStatus) error

	CalculateOrderTotal(ctx context.Context, orderID uuid.UUID) (float64, error)
	GetCustomerWalletAddress(ctx context.Context, customerID uuid.UUID) (string, error)

	AddItemIntoOrder(ctx context.Context, orderID uuid.UUID, item OrderItemInput) error
	RemoveItemFromOrder(ctx context.Context, orderID uuid.UUID, restaurantItemID uuid.UUID) error
}

type OrderItemInput struct {
	RestaurantItemID uuid.UUID
	Price            float64
	Quantity         int
}

type Filter struct {
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
