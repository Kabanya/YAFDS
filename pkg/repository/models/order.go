package models

import (
	"context"

	"github.com/Kabanya/YAFDS/pkg/models"

	"github.com/google/uuid"
)

type Order interface { //dto - data transfer object. dto похож на model но другое. Если совпадают, то приоритет модели
	//                   по возможности не плодим dto
	Create(ctx context.Context, order Order) (models.Order /*Order*/, error)
	CreateWithItems(ctx context.Context, order Order, items []OrderItemInput) (Order, error)
	List(ctx context.Context, filter Filter) ([]Order, error)
	Get(ctx context.Context, orderID uuid.UUID) (Order, error)
	GetOrderStatus(ctx context.Context, orderID uuid.UUID) (models.OrderStatus, error)
	UpdateStatus(ctx context.Context, orderID uuid.UUID, status models.OrderStatus) error
	GetOrderTotal(ctx context.Context, orderID uuid.UUID) (float64, error)
	GetCustomerWalletAddress(ctx context.Context, customerID uuid.UUID) (string, error)
	Accept(ctx context.Context, input AcceptInput) (AcceptResult, error)
	AddItem(ctx context.Context, orderID uuid.UUID, item OrderItemInput) error
}
