package service

import (
	"context"
	"errors"

	"github.com/Kabanya/YAFDS/pkg/order/domain"
	"github.com/Kabanya/YAFDS/pkg/order/usecase"
	"github.com/google/uuid"
)

type OrderService interface {
	GetOrder(ctx context.Context, orderID uuid.UUID) (domain.Order, error)
	ListOrders(ctx context.Context, filter usecase.Filter) ([]domain.Order, error)

	GetOrderStatus(ctx context.Context, orderID uuid.UUID) (domain.OrderStatus, error)
	UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status domain.OrderStatus) error

	AddItemIntoOrder(ctx context.Context, orderID uuid.UUID, item usecase.OrderItemInput) error
	RemoveItemFromOrder(ctx context.Context, orderID uuid.UUID, restaurantItemID uuid.UUID) error
	CalculateOrderTotal(ctx context.Context, orderID uuid.UUID) (float64, error)
}

type orderService struct {
	repo usecase.OrderRepository
}

func NewOrderService(repo usecase.OrderRepository) OrderService {
	return &orderService{repo: repo}
}

type CreateOrderResponse struct {
	OrderID string `json:"order_id"`
}

var (
	ErrNotPayedOrderState = errors.New("order not in paid state : OrderStatusCustomerPaid")
)

func (os *orderService) ListOrders(ctx context.Context, filter usecase.Filter) ([]domain.Order, error) {
	return os.repo.ListOrders(ctx, filter)
}

func (os *orderService) GetOrder(ctx context.Context, orderID uuid.UUID) (domain.Order, error) {
	return os.repo.GetOrder(ctx, orderID)
}

func (os *orderService) GetOrderStatus(ctx context.Context, orderID uuid.UUID) (domain.OrderStatus, error) {
	return os.repo.GetOrderStatus(ctx, orderID)
}

func (os *orderService) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status domain.OrderStatus) error {
	return os.repo.UpdateOrderStatus(ctx, orderID, status)
}

func (os *orderService) AddItemIntoOrder(ctx context.Context, orderID uuid.UUID, item usecase.OrderItemInput) error {
	return os.repo.AddItemIntoOrder(ctx, orderID, item)
}

func (os *orderService) RemoveItemFromOrder(ctx context.Context, orderID uuid.UUID, restaurantItemID uuid.UUID) error {
	return os.repo.RemoveItemFromOrder(ctx, orderID, restaurantItemID)
}

func (os *orderService) CalculateOrderTotal(ctx context.Context, orderID uuid.UUID) (float64, error) {
	items, err := os.repo.GetOrderItems(ctx, orderID)
	if err != nil {
		return 0, err
	}

	var total float64
	for _, item := range items {
		total += item.Price * float64(item.Quantity)
	}
	return total, nil
}
