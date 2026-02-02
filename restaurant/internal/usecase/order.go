package usecase

import (
	"context"

	"restaurant/internal/service"

	"github.com/Kabanya/YAFDS/pkg/models"

	"github.com/google/uuid"
)

type OrdersUseCase interface {
	ListOrdersByRestaurantID(ctx context.Context, restaurantID uuid.UUID, status string) ([]models.Order, error)
}

type ordersUseCase struct {
	// service service.OrdersService
	orderService service.OrdersService
}

func (u *ordersUseCase) ListOrdersByRestaurantID(ctx context.Context, restaurantID uuid.UUID, status string) ([]models.Order, error) {
	panic("unimplemented")
}

func NewOrdersUseCase(service service.OrdersService) OrdersUseCase {
	return &ordersUseCase{orderService: service}
}

// func (u *ordersUseCase) ListOrdersByRestaurantID(ctx context.Context, restaurantID uuid.UUID, status string) ([]models.Order, error) {
// 	return u.service.ListOrdersByRestaurantID(ctx, restaurantID, status)
// }

func (u *ordersUseCase) AcceptOrder(ctx context.Context, orderID uuid.UUID) error {
	return nil
}
