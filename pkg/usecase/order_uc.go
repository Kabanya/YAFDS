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
	CreateOrder(ctx context.Context, customerID uuid.UUID, courierID uuid.UUID) (models.Order, error)
	CreateOrderWithItems(ctx context.Context, customerID uuid.UUID, courierID uuid.UUID, items []repositoryModels.OrderItemInput) (models.Order, error)
	GetOrder(ctx context.Context, orderID uuid.UUID) (models.Order, error)
	AcceptOrder(ctx context.Context, orderID uuid.UUID, customerID uuid.UUID, courierID uuid.UUID, items []repositoryModels.OrderItemInput, status models.OrderStatus) (repositoryModels.AcceptResult, error)
	GetOrderStatus(ctx context.Context, orderID uuid.UUID) (models.OrderStatus, error)
	UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status models.OrderStatus) error
	CalculateOrderTotal(ctx context.Context, orderID uuid.UUID) (float64, error)
	GetCustomerWalletAddress(ctx context.Context, customerID uuid.UUID) (string, error)
	AddItemIntoOrder(ctx context.Context, orderID uuid.UUID, item repositoryModels.OrderItemInput) error
	RemoveItemFromOrder(ctx context.Context, orderID uuid.UUID, restaurantItemID uuid.UUID) error
}

type orderUseCase struct {
	serviceOrder service.OrderService
}

func NewOrderUseCase(serviceOrder service.OrderService) OrderUseCase {
	return &orderUseCase{serviceOrder: serviceOrder}
}

func (u *orderUseCase) CreateOrder(ctx context.Context, customerID uuid.UUID, courierID uuid.UUID) (models.Order, error) {
	resp, err := u.serviceOrder.CreateOrder(ctx, customerID.String(), courierID.String(), models.OrderStatusCustomerCreated)
	if err != nil {
		return models.Order{}, err
	}
	orderID, err := uuid.Parse(resp.OrderID)
	if err != nil {
		return models.Order{}, err
	}
	return u.serviceOrder.GetOrder(ctx, orderID)
}

func (u *orderUseCase) CreateOrderWithItems(ctx context.Context, customerID uuid.UUID, courierID uuid.UUID, items []repositoryModels.OrderItemInput) (models.Order, error) {
	resp, err := u.serviceOrder.CreateOrderWithItems(ctx, customerID.String(), courierID.String(), models.OrderStatusCustomerCreated, items)
	if err != nil {
		return models.Order{}, err
	}
	orderID, err := uuid.Parse(resp.OrderID)
	if err != nil {
		return models.Order{}, err
	}
	return u.serviceOrder.GetOrder(ctx, orderID)
}

func (u *orderUseCase) GetOrder(ctx context.Context, orderID uuid.UUID) (models.Order, error) {
	return u.serviceOrder.GetOrder(ctx, orderID)
}

func (u *orderUseCase) AcceptOrder(ctx context.Context, orderID uuid.UUID, customerID uuid.UUID, courierID uuid.UUID, items []repositoryModels.OrderItemInput, status models.OrderStatus) (repositoryModels.AcceptResult, error) {
	return u.serviceOrder.AcceptOrder(ctx, orderID.String(), customerID.String(), courierID.String(), items, status)
}

func (u *orderUseCase) GetOrderStatus(ctx context.Context, orderID uuid.UUID) (models.OrderStatus, error) {
	return u.serviceOrder.GetOrderStatus(ctx, orderID)
}

func (u *orderUseCase) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status models.OrderStatus) error {
	return u.serviceOrder.UpdateOrderStatus(ctx, orderID, status)
}

func (u *orderUseCase) CalculateOrderTotal(ctx context.Context, orderID uuid.UUID) (float64, error) {
	return u.serviceOrder.CalculateOrderTotal(ctx, orderID)
}

func (u *orderUseCase) GetCustomerWalletAddress(ctx context.Context, customerID uuid.UUID) (string, error) {
	return u.serviceOrder.GetCustomerWalletAddress(ctx, customerID)
}

func (u *orderUseCase) AddItemIntoOrder(ctx context.Context, orderID uuid.UUID, item repositoryModels.OrderItemInput) error {
	return u.serviceOrder.AddItemIntoOrder(ctx, orderID, item)
}

func (u *orderUseCase) RemoveItemFromOrder(ctx context.Context, orderID uuid.UUID, restaurantItemID uuid.UUID) error {
	return u.serviceOrder.RemoveItemFromOrder(ctx, orderID, restaurantItemID)
}
