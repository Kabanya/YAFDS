package usecase

import (
	"context"
	"errors"

	"github.com/Kabanya/YAFDS/pkg/order/domain"
	"github.com/google/uuid"
)

var (
	ErrInvalidStatusTransition = errors.New("invalid order status transition")
	ErrWalletUnavailable       = errors.New("wallet service unavailable")
	ErrInsufficientFunds       = errors.New("insufficient funds")
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
	Status     domain.OrderStatus
}

type AcceptInput struct {
	OrderID    uuid.UUID
	CustomerID uuid.UUID
	CourierID  uuid.UUID
	Items      []OrderItemInput
	Status     domain.OrderStatus
}

type AcceptResult struct {
	OrderID uuid.UUID
	Status  domain.OrderStatus
}

type OrderRepository interface {
	GetOrder(ctx context.Context, orderID uuid.UUID) (domain.Order, error)
	ListOrders(ctx context.Context, filter Filter) ([]domain.Order, error)

	GetOrderStatus(ctx context.Context, orderID uuid.UUID) (domain.OrderStatus, error)
	UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status domain.OrderStatus) error

	GetOrderItems(ctx context.Context, orderID uuid.UUID) ([]domain.MenuItem, error)
	AddItemIntoOrder(ctx context.Context, orderID uuid.UUID, item OrderItemInput) error
	RemoveItemFromOrder(ctx context.Context, orderID uuid.UUID, restaurantItemID uuid.UUID) error
}

type OrderUseCase interface {
	GetOrder(ctx context.Context, orderID uuid.UUID) (domain.Order, error)
	ListOrders(ctx context.Context, filter Filter) ([]domain.Order, error)

	GetOrderStatus(ctx context.Context, orderID uuid.UUID) (domain.OrderStatus, error)
	UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status domain.OrderStatus) error

	AddItemIntoOrder(ctx context.Context, orderID uuid.UUID, item OrderItemInput) error
	RemoveItemFromOrder(ctx context.Context, orderID uuid.UUID, restaurantItemID uuid.UUID) error
	CalculateOrderTotal(ctx context.Context, orderID uuid.UUID) (float64, error)
}

type orderUseCase struct {
	repo OrderRepository
}

func NewOrderUseCase(repo OrderRepository) OrderUseCase {
	return &orderUseCase{repo: repo}
}

func (u *orderUseCase) GetOrder(ctx context.Context, orderID uuid.UUID) (domain.Order, error) {
	return u.repo.GetOrder(ctx, orderID)
}

func (u *orderUseCase) ListOrders(ctx context.Context, filter Filter) ([]domain.Order, error) {
	return u.repo.ListOrders(ctx, filter)
}

func (u *orderUseCase) GetOrderStatus(ctx context.Context, orderID uuid.UUID) (domain.OrderStatus, error) {
	return u.repo.GetOrderStatus(ctx, orderID)
}

func (u *orderUseCase) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status domain.OrderStatus) error {
	return u.repo.UpdateOrderStatus(ctx, orderID, status)
}

func (u *orderUseCase) AddItemIntoOrder(ctx context.Context, orderID uuid.UUID, item OrderItemInput) error {
	return u.repo.AddItemIntoOrder(ctx, orderID, item)
}

func (u *orderUseCase) RemoveItemFromOrder(ctx context.Context, orderID uuid.UUID, restaurantItemID uuid.UUID) error {
	return u.repo.RemoveItemFromOrder(ctx, orderID, restaurantItemID)
}

func (u *orderUseCase) CalculateOrderTotal(ctx context.Context, orderID uuid.UUID) (float64, error) {
	items, err := u.repo.GetOrderItems(ctx, orderID)
	if err != nil {
		return 0, err
	}

	var total float64
	for _, item := range items {
		total += item.Price * float64(item.Quantity)
	}
	return total, nil
}
