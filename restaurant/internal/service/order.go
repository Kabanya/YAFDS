package service

import (
	"context"

	pkg_repository "customer/pkg/repository"
	"restaurant/internal/repository"

	"github.com/google/uuid"
)

type OrdersService interface {
	ListOrdersByRestaurantID(ctx context.Context, restaurantID uuid.UUID, status string) ([]pkg_repository.Order, error)
}

type ordersService struct {
	repo repository.OrdersRepo
}

func NewOrdersService(repo repository.OrdersRepo) OrdersService {
	return &ordersService{repo: repo}
}

func (s *ordersService) ListOrdersByRestaurantID(ctx context.Context, restaurantID uuid.UUID, status string) ([]pkg_repository.Order, error) {
	return s.repo.ListOrdersByRestaurantID(ctx, restaurantID, status)
}
