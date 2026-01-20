package usecase

import (
	"context"

	pkg_repository "customer/pkg/repository"
	"restaurant/internal/service"

	"github.com/google/uuid"
)

type OrdersUseCase interface {
	ListOrdersByRestaurantID(ctx context.Context, restaurantID uuid.UUID, status string) ([]pkg_repository.Order, error)
}

type ordersUseCase struct {
	service service.OrdersService
}

func NewOrdersUseCase(service service.OrdersService) OrdersUseCase {
	return &ordersUseCase{service: service}
}

func (u *ordersUseCase) ListOrdersByRestaurantID(ctx context.Context, restaurantID uuid.UUID, status string) ([]pkg_repository.Order, error) {
	return u.service.ListOrdersByRestaurantID(ctx, restaurantID, status)
}
