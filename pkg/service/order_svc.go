package service

import (
	"context"

	"github.com/Kabanya/YAFDS/pkg/models"
	pkgRepoModels "github.com/Kabanya/YAFDS/pkg/repository/models"
	"github.com/google/uuid"
)

type OrderService interface {
	CreateOrder(ctx context.Context, customerID string, courierID string, status models.OrderStatus) (CreateOrderResponse, error)
	CreateOrderWithItems(ctx context.Context, customerID string, courierID string, status models.OrderStatus, items []pkgRepoModels.OrderItemInput) (CreateOrderResponse, error)
	ListOrders(ctx context.Context, filter pkgRepoModels.Filter) ([]models.Order, error)
	GetOrder(ctx context.Context, orderID uuid.UUID) (models.Order, error)
	AcceptOrder(ctx context.Context, OrderID string, CustomerID string, CourierID string, items []pkgRepoModels.OrderItemInput, status models.OrderStatus) (pkgRepoModels.AcceptResult, error)
	GetOrderStatus(ctx context.Context, orderID uuid.UUID) (models.OrderStatus, error)
	UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status models.OrderStatus) error
	CalculateOrderTotal(ctx context.Context, orderID uuid.UUID) (float64, error)
	GetCustomerWalletAddress(ctx context.Context, customerID uuid.UUID) (string, error)
	AddItemIntoOrder(ctx context.Context, orderID uuid.UUID, item pkgRepoModels.OrderItemInput) error
	RemoveItemFromOrder(ctx context.Context, orderID uuid.UUID, restaurantItemID uuid.UUID) error
}

type orderService struct {
	repo pkgRepoModels.OrderRepo
}

func NewOrderService(repo pkgRepoModels.OrderRepo) OrderService {
	return &orderService{repo: repo}
}

type CreateOrderResponse struct {
	OrderID string `json:"order_id"`
}

func (os *orderService) CreateOrder(ctx context.Context, customerID string, courierID string, status models.OrderStatus) (CreateOrderResponse, error) {
	custID, err := uuid.Parse(customerID)
	if err != nil {
		return CreateOrderResponse{}, err
	}
	courID, err := uuid.Parse(courierID)
	if err != nil {
		return CreateOrderResponse{}, err
	}

	order := models.Order{
		CustomerID: custID,
		CourierID:  courID,
		Status:     status,
	}

	createdOrder, err := os.repo.CreateOrder(ctx, order)
	if err != nil {
		return CreateOrderResponse{}, err
	}

	return CreateOrderResponse{
		OrderID: createdOrder.ID.String(),
	}, nil
}

func (os *orderService) CreateOrderWithItems(ctx context.Context, customerID string, courierID string, status models.OrderStatus, items []pkgRepoModels.OrderItemInput) (CreateOrderResponse, error) {
	custID, err := uuid.Parse(customerID)
	if err != nil {
		return CreateOrderResponse{}, err
	}
	courID, err := uuid.Parse(courierID)
	if err != nil {
		return CreateOrderResponse{}, err
	}

	order := models.Order{
		CustomerID: custID,
		CourierID:  courID,
		Status:     models.OrderStatus(status),
	}

	createdOrder, err := os.repo.CreateOrderWithItems(ctx, order, items)
	if err != nil {
		return CreateOrderResponse{}, err
	}

	return CreateOrderResponse{
		OrderID: createdOrder.ID.String(),
	}, nil
}

func (os *orderService) ListOrders(ctx context.Context, filter pkgRepoModels.Filter) ([]models.Order, error) {
	return os.repo.ListOrders(ctx, filter)
}

// нужны ли в простых методах сервиса какие либо провеки если проверки были на уровне repository
// func (os *orderService) ListOrders(ctx context.Context, filter pkgRepoModels.Filter) ([]models.Order, error) {
// 	if filter.CustomerID != nil {
// 		_, err := uuid.Parse(filter.CustomerID.String())
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
// 	if filter.CourierID != nil {
// 		_, err := uuid.Parse(filter.CourierID.String())
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
// 	return os.repo.ListOrders(ctx, filter)
// }

func (os *orderService) GetOrder(ctx context.Context, orderID uuid.UUID) (models.Order, error) {
	return os.repo.GetOrder(ctx, orderID)
}

func (os *orderService) AcceptOrder(ctx context.Context, OrderID string, CustomerID string, CourierID string, Items []pkgRepoModels.OrderItemInput, Status models.OrderStatus) (pkgRepoModels.AcceptResult, error) { // Updated signature
	orderID, err := uuid.Parse(OrderID)
	if err != nil {
		return pkgRepoModels.AcceptResult{}, err
	}
	customerID, err := uuid.Parse(CustomerID)
	if err != nil {
		return pkgRepoModels.AcceptResult{}, err
	}
	courierID, err := uuid.Parse(CourierID)
	if err != nil {
		return pkgRepoModels.AcceptResult{}, err
	}

	input := pkgRepoModels.AcceptInput{
		OrderID:    orderID,
		CustomerID: customerID,
		CourierID:  courierID,
		Items:      Items,
		Status:     Status,
	}

	return os.repo.AcceptOrder(ctx, input)
}

func (os *orderService) GetOrderStatus(ctx context.Context, orderID uuid.UUID) (models.OrderStatus, error) {
	return os.repo.GetOrderStatus(ctx, orderID)
}

func (os *orderService) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status models.OrderStatus) error {
	return os.repo.UpdateOrderStatus(ctx, orderID, status)
}

func (os *orderService) CalculateOrderTotal(ctx context.Context, orderID uuid.UUID) (float64, error) {
	return os.repo.CalculateOrderTotal(ctx, orderID)
}

func (os *orderService) GetCustomerWalletAddress(ctx context.Context, customerID uuid.UUID) (string, error) {
	return os.repo.GetCustomerWalletAddress(ctx, customerID)
}

func (os *orderService) AddItemIntoOrder(ctx context.Context, orderID uuid.UUID, item pkgRepoModels.OrderItemInput) error {
	return os.repo.AddItemIntoOrder(ctx, orderID, item)
}

func (os *orderService) RemoveItemFromOrder(ctx context.Context, orderID uuid.UUID, restaurantItemID uuid.UUID) error {
	return os.repo.RemoveItemFromOrder(ctx, orderID, restaurantItemID)
}
