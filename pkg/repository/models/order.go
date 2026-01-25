package models

import (
	"context"

	"github.com/Kabanya/YAFDS/pkg/models"

	"github.com/google/uuid"
)

type Order interface { //dto - data transfer object. dto похож на model но другое. Если совпадают, то приоритет модели
	//                   по возможности не плодим dto
	Create(ctx context.Context, order models.Order) (models.Order, error)
	CreateWithItems(ctx context.Context, order models.Order, items []OrderItemInput) (models.Order, error)
	List(ctx context.Context, filter Filter) ([]models.Order, error)
	Get(ctx context.Context, orderID uuid.UUID) (models.Order, error)
	GetOrderStatus(ctx context.Context, orderID uuid.UUID) (models.OrderStatus, error)
	UpdateStatus(ctx context.Context, orderID uuid.UUID, status models.OrderStatus) error
	GetOrderTotal(ctx context.Context, orderID uuid.UUID) (float64, error)
	GetCustomerWalletAddress(ctx context.Context, customerID uuid.UUID) (string, error)
	Accept(ctx context.Context, input AcceptInput) (AcceptResult, error)
	AddItem(ctx context.Context, orderID uuid.UUID, item OrderItemInput) error
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
