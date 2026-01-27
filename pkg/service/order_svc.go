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

type createOrderRequest struct {
	CustomerID string `json:"customer_id"`
	CourierID  string `json:"courier_id"`
	Status     string `json:"status"`
}

type acceptOrderRequest struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

type createOrderResponce struct {
	OrderID string `json:"order_id"`
}

func CreateOrder(ctx context.Context, repo pkgRepoModels.OrderRepo, req createOrderRequest) (createOrderResponce, error) {
	customerID, err := uuid.Parse(req.CustomerID)
	if err != nil {
		return createOrderResponce{}, err
	}
	courierID, err := uuid.Parse(req.CourierID)
	if err != nil {
		return createOrderResponce{}, err
	}

	order := models.Order{
		CustomerID: customerID,
		CourierID:  courierID,
		Status:     req.Status,
	}

	createdOrder, err := repo.CreateOrder(ctx, order)
	if err != nil {
		return createOrderResponce{}, err
	}

	return createOrderResponce{
		OrderID: createdOrder.ID.String(),
	}, nil
}

func CreateOrderWithItems(ctx context.Context, repo pkgRepoModels.OrderRepo, req createOrderRequest, items []pkgRepoModels.OrderItemInput) (createOrderResponce, error) {
	customerID, err := uuid.Parse(req.CustomerID)
	if err != nil {
		return createOrderResponce{}, err
	}
	courierID, err := uuid.Parse(req.CourierID)
	if err != nil {
		return createOrderResponce{}, err
	}

	order := models.Order{
		CustomerID: customerID,
		CourierID:  courierID,
		Status:     req.Status,
	}

	createdOrder, err := repo.CreateOrderWithItems(ctx, order, items)
	if err != nil {
		return createOrderResponce{}, err
	}

	return createOrderResponce{
		OrderID: createdOrder.ID.String(),
	}, nil
}

func ListOrders(ctx context.Context, repo pkgRepoModels.OrderRepo, filter pkgRepoModels.Filter) ([]models.Order, error) {
	orders, err := repo.ListOrders(ctx, filter)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func GetOrder(ctx context.Context, repo pkgRepoModels.OrderRepo, orderID uuid.UUID) (models.Order, error) {
	order, err := repo.GetOrder(ctx, orderID)
	if err != nil {
		return models.Order{}, err
	}
	return order, nil
}

func AcceptOrder(ctx context.Context, repo pkgRepoModels.OrderRepo, req acceptOrderRequest) (pkgRepoModels.AcceptResult, error) {
	orderID, err := uuid.Parse(req.OrderID)
	if err != nil {
		return pkgRepoModels.AcceptResult{}, err
	}

	input := pkgRepoModels.AcceptInput{
		OrderID: orderID,
		Status:  models.OrderStatus(req.Status),
	}

	result, err := repo.AcceptOrder(ctx, input)
	if err != nil {
		return pkgRepoModels.AcceptResult{}, err
	}

	return result, nil
}

func GetOrderStatus(ctx context.Context, repo pkgRepoModels.OrderRepo, orderID uuid.UUID) (models.OrderStatus, error) {
	status, err := repo.GetOrderStatus(ctx, orderID)
	if err != nil {
		return models.OrderStatus("UNDEFINED"), err
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

func CalculateOrderTotal(ctx context.Context, repo pkgRepoModels.OrderRepo, orderID uuid.UUID) (float64, error) {
	total, err := repo.CalculateOrderTotal(ctx, orderID)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func GetCustomerWalletAddress(ctx context.Context, repo pkgRepoModels.OrderRepo, customerID uuid.UUID) (string, error) {
	address, err := repo.GetCustomerWalletAddress(ctx, customerID)
	if err != nil {
		return "", err
	}
	return address, nil
}

func AddItemIntoOrder(ctx context.Context, repo pkgRepoModels.OrderRepo, orderID uuid.UUID, item pkgRepoModels.OrderItemInput) error {
	err := repo.AddItemIntoOrder(ctx, orderID, item)
	if err != nil {
		return err
	}
	return nil
}

func RemoveItemFromOrder(ctx context.Context, repo pkgRepoModels.OrderRepo, orderID uuid.UUID, restaurantItemID uuid.UUID) error {
	err := repo.RemoveItemFromOrder(ctx, orderID, restaurantItemID)
	if err != nil {
		return err
	}
	return nil
}
