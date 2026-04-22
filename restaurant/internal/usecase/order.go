package usecase

import (
	"context"

	"github.com/Kabanya/YAFDS/pkg/order/domain"
	pkgOrderUsecase "github.com/Kabanya/YAFDS/pkg/order/usecase"
	"github.com/Kabanya/YAFDS/restaurant/internal/service"

	"github.com/google/uuid"
)

type OrderUseCase interface {
	ListOrdersByRestaurantID(ctx context.Context, restaurantID uuid.UUID, status string) ([]domain.Order, error)
	AcceptOrder(ctx context.Context, orderID uuid.UUID) (pkgOrderUsecase.AcceptResult, error)
}

type orderUseCase struct {
	orderService service.OrdersService
}

func NewOrderUseCase(service service.OrdersService) OrderUseCase {
	return &orderUseCase{orderService: service}
}

func (u *orderUseCase) ListOrdersByRestaurantID(ctx context.Context, restaurantID uuid.UUID, status string) ([]domain.Order, error) {
	return u.orderService.ListOrdersByRestaurantID(ctx, restaurantID, status)
}

func (u *orderUseCase) AcceptOrder(ctx context.Context, orderID uuid.UUID) (pkgOrderUsecase.AcceptResult, error) {
	return u.orderService.AcceptOrder(ctx, orderID)
}
