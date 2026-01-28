package service

// НЕ ЗНАЕТ ПРО HTTP .
// Создается, есть репу, логеры.
// есть func create order
import (
	"context"

	"github.com/Kabanya/YAFDS/pkg/models"
	pkgRepoModels "github.com/Kabanya/YAFDS/pkg/repository/models"
	"github.com/google/uuid"
)

type OrderService interface {
	CreateOrder(ctx context.Context, repo pkgRepoModels.OrderRepo, req createOrderRequest) (createOrderResponse, error)
	CreateOrderWithItems(ctx context.Context, repo pkgRepoModels.OrderRepo, req createOrderRequest, items []pkgRepoModels.OrderItemInput) (createOrderResponse, error)
	ListOrders(ctx context.Context, repo pkgRepoModels.OrderRepo, filter pkgRepoModels.Filter) ([]models.Order, error)
	GetOrder(ctx context.Context, repo pkgRepoModels.OrderRepo, orderID uuid.UUID) (models.Order, error)
	AcceptOrder(ctx context.Context, repo pkgRepoModels.OrderRepo, req acceptOrderRequest) (pkgRepoModels.AcceptResult, error)
	GetOrderStatus(ctx context.Context, repo pkgRepoModels.OrderRepo, orderID uuid.UUID) (models.OrderStatus, error)
	UpdateOrderStatus(ctx context.Context, repo pkgRepoModels.OrderRepo, orderID uuid.UUID, status models.OrderStatus) error
	CalculateOrderTotal(ctx context.Context, repo pkgRepoModels.OrderRepo, orderID uuid.UUID) (float64, error)
	GetCustomerWalletAddress(ctx context.Context, repo pkgRepoModels.OrderRepo, customerID uuid.UUID) (string, error)
	AddItemIntoOrder(ctx context.Context, repo pkgRepoModels.OrderRepo, orderID uuid.UUID, item pkgRepoModels.OrderItemInput) error
	RemoveItemFromOrder(ctx context.Context, repo pkgRepoModels.OrderRepo, orderID uuid.UUID, restaurantItemID uuid.UUID) error
}

type orderService struct {
	repo pkgRepoModels.OrderRepo
}

type createOrderRequest struct {
	CustomerID string `json:"customer_id"`
	CourierID  string `json:"courier_id"`
	Status     string `json:"status"`
}

type acceptOrderRequest struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

type createOrderResponse struct {
	OrderID string `json:"order_id"`
}

func (os *orderService) CreateOrder(ctx context.Context, req createOrderRequest) (createOrderResponse, error) {
	customerID, err := uuid.Parse(req.CustomerID)
	if err != nil {
		return createOrderResponse{}, err
	}
	courierID, err := uuid.Parse(req.CourierID)
	if err != nil {
		return createOrderResponse{}, err
	}

	order := models.Order{
		CustomerID: customerID,
		CourierID:  courierID,
		Status:     models.OrderStatusCustomerCreated,
	}

	createdOrder, err := os.repo.CreateOrder(ctx, order)
	if err != nil {
		return createOrderResponse{}, err
	}

	return createOrderResponse{
		OrderID: createdOrder.ID.String(),
	}, nil
}

func CreateOrderWithItems(ctx context.Context, repo pkgRepoModels.OrderRepo, req createOrderRequest, items []pkgRepoModels.OrderItemInput) (createOrderResponse, error) {
	customerID, err := uuid.Parse(req.CustomerID)
	if err != nil {
		return createOrderResponse{}, err
	}
	courierID, err := uuid.Parse(req.CourierID)
	if err != nil {
		return createOrderResponse{}, err
	}

	order := models.Order{
		CustomerID: customerID,
		CourierID:  courierID,
		Status:     models.OrderStatusCustomerCreated,
	}

	createdOrder, err := repo.CreateOrderWithItems(ctx, order, items)
	if err != nil {
		return createOrderResponse{}, err
	}

	return createOrderResponse{
		OrderID: createdOrder.ID.String(),
	}, nil
}

func (os *orderService) ListOrders(ctx context.Context, filter pkgRepoModels.Filter) ([]models.Order, error) {
	orders, err := os.repo.ListOrders(ctx, filter)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (os *orderService) GetOrder(ctx context.Context, orderID uuid.UUID) (models.Order, error) {
	order, err := os.repo.GetOrder(ctx, orderID)
	if err != nil {
		return models.Order{}, err
	}
	return order, nil
}

func (os *orderService) AcceptOrder(ctx context.Context, req acceptOrderRequest) (pkgRepoModels.AcceptResult, error) {
	orderID, err := uuid.Parse(req.OrderID)
	if err != nil {
		return pkgRepoModels.AcceptResult{}, err
	}

	input := pkgRepoModels.AcceptInput{
		OrderID: orderID,
		Status:  models.OrderStatus(req.Status),
	}

	result, err := os.repo.AcceptOrder(ctx, input)
	if err != nil {
		return pkgRepoModels.AcceptResult{}, err
	}

	return result, nil
}

func GetOrderStatus(ctx context.Context, repo pkgRepoModels.OrderRepo, orderID uuid.UUID) (models.OrderStatus, error) {
	status, err := repo.GetOrderStatus(ctx, orderID)
	if err != nil {
		return models.OrderStatusOrderFailed, err
	}
	return status, nil
}

func UpdateOrderStatus(ctx context.Context, repo pkgRepoModels.OrderRepo, orderID uuid.UUID, status models.OrderStatus) error {
	err := repo.UpdateOrderStatus(ctx, orderID, status)
	if err != nil {
		return err
	}
	return nil
}

func (os *orderService) CalculateOrderTotal(ctx context.Context, orderID uuid.UUID) (float64, error) {
	total, err := os.repo.CalculateOrderTotal(ctx, orderID)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (os *orderService) GetCustomerWalletAddress(ctx context.Context, customerID uuid.UUID) (string, error) {
	address, err := os.repo.GetCustomerWalletAddress(ctx, customerID)
	if err != nil {
		return "", err
	}
	return address, nil
}

func (os *orderService) AddItemIntoOrder(ctx context.Context, orderID uuid.UUID, item pkgRepoModels.OrderItemInput) error {
	err := os.repo.AddItemIntoOrder(ctx, orderID, item)
	if err != nil {
		return err
	}
	return nil
}

func (os *orderService) RemoveItemFromOrder(ctx context.Context, orderID uuid.UUID, restaurantItemID uuid.UUID) error {
	err := os.repo.RemoveItemFromOrder(ctx, orderID, restaurantItemID)
	if err != nil {
		return err
	}
	return nil
}
