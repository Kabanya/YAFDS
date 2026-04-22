package usecase

import (
	"context"
	"errors"

	customerservice "github.com/Kabanya/YAFDS/customer/internal/service"
	domain "github.com/Kabanya/YAFDS/pkg/order/domain"
	orderservice "github.com/Kabanya/YAFDS/pkg/order/service"
	models "github.com/Kabanya/YAFDS/pkg/order/usecase"
	"github.com/Kabanya/YAFDS/pkg/wallet"
	"github.com/google/uuid"
)

type OrderUseCase interface {
	CreateOrder(ctx context.Context, customerID uuid.UUID, courierID uuid.UUID) (domain.Order, error)
	CreateOrderWithItems(ctx context.Context, customerID uuid.UUID, courierID uuid.UUID, items []models.OrderItemInput) (domain.Order, error)
	GetOrder(ctx context.Context, orderID uuid.UUID) (domain.Order, error)
	ListOrders(ctx context.Context, filter models.Filter) ([]domain.Order, error)

	GetOrderStatus(ctx context.Context, orderID uuid.UUID) (domain.OrderStatus, error)
	UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status domain.OrderStatus) error

	AddItemIntoOrder(ctx context.Context, orderID uuid.UUID, item models.OrderItemInput) error
	RemoveItemFromOrder(ctx context.Context, orderID uuid.UUID, restaurantItemID uuid.UUID) error
	CalculateOrderTotal(ctx context.Context, orderID uuid.UUID) (float64, error)
	GetCustomerWalletAddress(ctx context.Context, customerID uuid.UUID) (string, error)
	PayOrder(ctx context.Context, orderID uuid.UUID) error
}

type orderUseCase struct {
	orderService         orderservice.OrderService
	customerOrderService customerservice.OrderService
	walletClient         wallet.WalletClient
}

func NewOrderUseCase(orderService orderservice.OrderService, customerOrderService customerservice.OrderService, walletClient wallet.WalletClient) OrderUseCase {
	return &orderUseCase{
		orderService:         orderService,
		customerOrderService: customerOrderService,
		walletClient:         walletClient,
	}
}

func (o *orderUseCase) CreateOrder(ctx context.Context, customerID uuid.UUID, courierID uuid.UUID) (domain.Order, error) {
	return o.customerOrderService.CreateOrder(ctx, customerID, courierID)
}

func (o *orderUseCase) CreateOrderWithItems(ctx context.Context, customerID uuid.UUID, courierID uuid.UUID, items []models.OrderItemInput) (domain.Order, error) {
	return o.customerOrderService.CreateOrderWithItems(ctx, customerID, courierID, items)
}

func (o *orderUseCase) GetOrder(ctx context.Context, orderID uuid.UUID) (domain.Order, error) {
	return o.orderService.GetOrder(ctx, orderID)
}

func (o *orderUseCase) ListOrders(ctx context.Context, filter models.Filter) ([]domain.Order, error) {
	return o.orderService.ListOrders(ctx, filter)
}

func (o *orderUseCase) GetOrderStatus(ctx context.Context, orderID uuid.UUID) (domain.OrderStatus, error) {
	return o.orderService.GetOrderStatus(ctx, orderID)
}

func (o *orderUseCase) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status domain.OrderStatus) error {
	return o.orderService.UpdateOrderStatus(ctx, orderID, status)
}

func (o *orderUseCase) AddItemIntoOrder(ctx context.Context, orderID uuid.UUID, item models.OrderItemInput) error {
	return o.orderService.AddItemIntoOrder(ctx, orderID, item)
}

func (o *orderUseCase) RemoveItemFromOrder(ctx context.Context, orderID uuid.UUID, restaurantItemID uuid.UUID) error {
	return o.orderService.RemoveItemFromOrder(ctx, orderID, restaurantItemID)
}

func (o *orderUseCase) CalculateOrderTotal(ctx context.Context, orderID uuid.UUID) (float64, error) {
	return o.orderService.CalculateOrderTotal(ctx, orderID)
}

func (o *orderUseCase) GetCustomerWalletAddress(ctx context.Context, customerID uuid.UUID) (string, error) {
	return o.customerOrderService.GetCustomerWalletAddress(ctx, customerID)
}

func (o *orderUseCase) PayOrder(ctx context.Context, orderID uuid.UUID) error {
	order, err := o.orderService.GetOrder(ctx, orderID)
	if err != nil {
		return err
	}
	if order.Status != domain.OrderStatusCustomerCreated {
		return errors.New("order can only be paid when in CUSTOMER_CREATED status")
	}

	walletAddress, err := o.customerOrderService.GetCustomerWalletAddress(ctx, order.CustomerID)
	if err != nil {
		return err
	}

	total, err := o.orderService.CalculateOrderTotal(ctx, orderID)
	if err != nil {
		return err
	}

	if o.walletClient == nil {
		return models.ErrWalletUnavailable
	}

	balance, err := o.walletClient.GetBalanceWallet(ctx, walletAddress)
	if err != nil {
		return models.ErrWalletUnavailable
	}
	if balance < total {
		return models.ErrInsufficientFunds
	}
	if err := o.walletClient.PayToWallet(ctx, walletAddress, total); err != nil {
		return errors.Join(models.ErrWalletUnavailable, err)
	}

	return o.customerOrderService.PayOrder(ctx, orderID)
}
