package service

import (
	"context"

	"restaurant/internal/repository"

	"github.com/Kabanya/YAFDS/pkg/models"

	pkgRepoModels "github.com/Kabanya/YAFDS/pkg/repository/models"
	pkgRepo "github.com/Kabanya/YAFDS/pkg/repository/postgres"
	pkgOrderService "github.com/Kabanya/YAFDS/pkg/service"
	"github.com/google/uuid"
)

type OrdersService interface {
	ListOrdersByRestaurantID(ctx context.Context, restaurantID uuid.UUID, status string) ([]models.Order, error)
	AcceptOrder(ctx context.Context, orderID uuid.UUID) (pkgRepoModels.AcceptResult, error)
	CalculateOrderTotal(ctx context.Context, orderID uuid.UUID) (float64, error)
}

type orderService struct {
	repo            repository.OrdersRepo
	pkgOrderRepo    pkgRepo.OrderRepo
	pkgOrderService pkgOrderService.OrderService
}

func NewOrderService(repo repository.OrdersRepo, pkgOrderRepo pkgRepo.OrderRepo, pkgOrderService pkgOrderService.OrderService) OrdersService {
	return &orderService{repo: repo, pkgOrderRepo: pkgOrderRepo, pkgOrderService: pkgOrderService}
}

func (s *orderService) ListOrdersByRestaurantID(ctx context.Context, restaurantID uuid.UUID, status string) ([]models.Order, error) {
	return s.repo.ListOrdersByRestaurantID(ctx, restaurantID, status)
}

func (s *orderService) CalculateOrderTotal(ctx context.Context, orderID uuid.UUID) (float64, error) {
	items, err := s.pkgOrderService.GetOrderItems(ctx, orderID)
	if err != nil {
		return 0, err
	}

	var total float64
	for _, item := range items {
		total += item.Price * float64(item.Quantity)
	}
	return total, nil
}

func (os *orderService) AcceptOrder(ctx context.Context, orderID uuid.UUID) (pkgRepoModels.AcceptResult, error) {
	status, err := os.pkgOrderRepo.GetOrderStatus(ctx, orderID)
	if err != nil {
		return pkgRepoModels.AcceptResult{}, err
	}
	if status != models.OrderStatusCustomerPaid {
		return pkgRepoModels.AcceptResult{}, pkgOrderService.ErrNotPayedOrderState
	}

	total, err := os.CalculateOrderTotal(ctx, orderID)
	if err != nil {
		return pkgRepoModels.AcceptResult{}, err
	}

	newStatus := models.OrderStatusKitchenAccepted
	if total <= 0 {
		newStatus = models.OrderStatusKitchenDenied
	}

	if err := os.pkgOrderService.UpdateOrderStatus(ctx, orderID, newStatus); err != nil {
		return pkgRepoModels.AcceptResult{}, err
	}

	return pkgRepoModels.AcceptResult{
		OrderID: orderID,
		Status:  string(newStatus),
	}, nil
}
