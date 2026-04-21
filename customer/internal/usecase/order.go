package usecase

import (
	"context"
	"github.com/Kabanya/YAFDS/customer/internal/models"
	orderservice "github.com/Kabanya/YAFDS/customer/internal/service"

	"github.com/google/uuid"
)

type OrderUseCase interface {
	GetOrder(ctx context.Context, orderID uuid.UUID) (models.Order, error)
	ListOrders(ctx context.Context, filter models.Filter) ([]models.Order, error)

	GetOrderStatus(ctx context.Context, orderID uuid.UUID) (models.OrderStatus, error)
	UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status models.OrderStatus) error

	AddItemIntoOrder(ctx context.Context, orderID uuid.UUID, item models.OrderItemInput) error
	RemoveItemFromOrder(ctx context.Context, orderID uuid.UUID, restaurantItemID uuid.UUID) error
}

type orderUseCase struct {
	service orderservice.OrderService
}

func NewOrderUseCase(service orderservice.OrderService) OrderUseCase {
	return &orderUseCase{service: service}
}

func (o *orderUseCase) GetOrder(ctx context.Context, orderID uuid.UUID) (models.Order, error) {
	return o.service.GetOrder(ctx, orderID)
}

func (o *orderUseCase) ListOrders(ctx context.Context, filter models.Filter) ([]models.Order, error) {
	return o.service.ListOrders(ctx, filter)
}

func (o *orderUseCase) GetOrderStatus(ctx context.Context, orderID uuid.UUID) (models.OrderStatus, error) {
	return o.service.GetOrderStatus(ctx, orderID)
}

func (o *orderUseCase) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status models.OrderStatus) error {
	return o.service.UpdateOrderStatus(ctx, orderID, status)
}

func (o *orderUseCase) AddItemIntoOrder(ctx context.Context, orderID uuid.UUID, item models.OrderItemInput) error {
	return o.service.AddItemIntoOrder(ctx, orderID, item)
}

func (o *orderUseCase) RemoveItemFromOrder(ctx context.Context, orderID uuid.UUID, restaurantItemID uuid.UUID) error {
	return o.service.RemoveItemFromOrder(ctx, orderID, restaurantItemID)
}
