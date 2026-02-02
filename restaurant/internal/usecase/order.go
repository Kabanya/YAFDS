package usecase

import (
	"context"

	"restaurant/internal/service"

	"github.com/Kabanya/YAFDS/pkg/models"

	pkgRepositoryModels "github.com/Kabanya/YAFDS/pkg/repository/models"
	"github.com/google/uuid"
)

type OrderUseCase interface {
	ListOrdersByRestaurantID(ctx context.Context, restaurantID uuid.UUID, status string) ([]models.Order, error)
	AcceptOrder(ctx context.Context, orderID uuid.UUID) (pkgRepositoryModels.AcceptResult, error)
}

type orderUseCase struct {
	// service service.OrdersService
	orderService service.OrdersService
}

func NewOrderUseCase(service service.OrdersService) OrderUseCase {
	return &orderUseCase{orderService: service}
}

func (u *orderUseCase) ListOrdersByRestaurantID(ctx context.Context, restaurantID uuid.UUID, status string) ([]models.Order, error) {
	panic("unimplemented")
}

// func (u *ordersUseCase) ListOrdersByRestaurantID(ctx context.Context, restaurantID uuid.UUID, status string) ([]models.Order, error) {
// 	return u.service.ListOrdersByRestaurantID(ctx, restaurantID, status)
// }

func (u *orderUseCase) AcceptOrder(ctx context.Context, orderID uuid.UUID) (pkgRepositoryModels.AcceptResult, error) {
	return u.orderService.AcceptOrder(ctx, orderID)
}
