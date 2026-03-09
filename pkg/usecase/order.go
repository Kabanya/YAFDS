package usecase

import (
	"context"
	"errors"

	"github.com/Kabanya/YAFDS/pkg/models"
	repositoryModels "github.com/Kabanya/YAFDS/pkg/repository/models"
	service "github.com/Kabanya/YAFDS/pkg/service"

	"github.com/google/uuid"
)

// НА УРОВНЕ USECASE ходим за продуктовыми логикой
var (
	ErrInvalidStatusTransition = errors.New("invalid order status transition")
	ErrWalletUnavailable       = errors.New("wallet service unavailable")
	ErrInsufficientFunds       = errors.New("insufficient funds")
)

type OrderUseCase interface {
	GetOrder(ctx context.Context, orderID uuid.UUID) (models.Order, error)
	ListOrders(ctx context.Context, filter repositoryModels.Filter) ([]models.Order, error)

	GetOrderStatus(ctx context.Context, orderID uuid.UUID) (models.OrderStatus, error)
	UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status models.OrderStatus) error

	AddItemIntoOrder(ctx context.Context, orderID uuid.UUID, item repositoryModels.OrderItemInput) error
	RemoveItemFromOrder(ctx context.Context, orderID uuid.UUID, restaurantItemID uuid.UUID) error
}

type orderUseCase struct {
	serviceOrder service.OrderService
}

func NewOrderUseCase(serviceOrder service.OrderService) OrderUseCase {
	return &orderUseCase{serviceOrder: serviceOrder}
}

func (u *orderUseCase) GetOrder(ctx context.Context, orderID uuid.UUID) (models.Order, error) {
	return u.serviceOrder.GetOrder(ctx, orderID)
}

func (u *orderUseCase) ListOrders(ctx context.Context, filter repositoryModels.Filter) ([]models.Order, error) {
	return u.serviceOrder.ListOrders(ctx, filter)
}

func (u *orderUseCase) GetOrderStatus(ctx context.Context, orderID uuid.UUID) (models.OrderStatus, error) {
	return u.serviceOrder.GetOrderStatus(ctx, orderID)
}

func (u *orderUseCase) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status models.OrderStatus) error {
	return u.serviceOrder.UpdateOrderStatus(ctx, orderID, status)
}

func (u *orderUseCase) AddItemIntoOrder(ctx context.Context, orderID uuid.UUID, item repositoryModels.OrderItemInput) error {
	return u.serviceOrder.AddItemIntoOrder(ctx, orderID, item)
}

func (u *orderUseCase) RemoveItemFromOrder(ctx context.Context, orderID uuid.UUID, restaurantItemID uuid.UUID) error {
	return u.serviceOrder.RemoveItemFromOrder(ctx, orderID, restaurantItemID)
}
