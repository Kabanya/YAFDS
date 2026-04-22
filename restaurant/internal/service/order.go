package service

import (
	"context"

	"github.com/Kabanya/YAFDS/pkg/order/domain"
	pkgRepo "github.com/Kabanya/YAFDS/pkg/order/repository/postgres"
	pkgOrderService "github.com/Kabanya/YAFDS/pkg/order/service"
	pkgOrderUsecase "github.com/Kabanya/YAFDS/pkg/order/usecase"
	repository "github.com/Kabanya/YAFDS/restaurant/internal/repository/postgres"
	"github.com/google/uuid"
)

type OrdersService interface {
	ListOrdersByRestaurantID(ctx context.Context, restaurantID uuid.UUID, status string) ([]domain.Order, error)
	AcceptOrder(ctx context.Context, orderID uuid.UUID) (pkgOrderUsecase.AcceptResult, error)
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

func (s *orderService) ListOrdersByRestaurantID(ctx context.Context, restaurantID uuid.UUID, status string) ([]domain.Order, error) {
	return s.repo.ListOrdersByRestaurantID(ctx, restaurantID, status)
}

func (s *orderService) CalculateOrderTotal(ctx context.Context, orderID uuid.UUID) (float64, error) {
	return s.pkgOrderService.CalculateOrderTotal(ctx, orderID)
}

func (os *orderService) AcceptOrder(ctx context.Context, orderID uuid.UUID) (pkgOrderUsecase.AcceptResult, error) {
	status, err := os.pkgOrderRepo.GetOrderStatus(ctx, orderID)
	if err != nil {
		return pkgOrderUsecase.AcceptResult{}, err
	}
	if status != domain.OrderStatusCustomerPaid {
		return pkgOrderUsecase.AcceptResult{}, pkgOrderService.ErrNotPayedOrderState
	}

	total, err := os.CalculateOrderTotal(ctx, orderID)
	if err != nil {
		return pkgOrderUsecase.AcceptResult{}, err
	}

	newStatus := domain.OrderStatusKitchenAccepted
	if total <= 0 {
		newStatus = domain.OrderStatusKitchenDenied
	}

	if err := os.pkgOrderService.UpdateOrderStatus(ctx, orderID, newStatus); err != nil {
		return pkgOrderUsecase.AcceptResult{}, err
	}

	return pkgOrderUsecase.AcceptResult{
		OrderID: orderID,
		Status:  newStatus,
	}, nil
}
