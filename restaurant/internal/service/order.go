package service

import (
	"context"

	"restaurant/internal/repository"

	"github.com/Kabanya/YAFDS/pkg/models"

	"github.com/google/uuid"
)

type OrdersService interface {
	ListOrdersByRestaurantID(ctx context.Context, restaurantID uuid.UUID, status string) ([]models.Order, error)
}

type ordersService struct {
	repo repository.OrdersRepo
}

func NewOrdersService(repo repository.OrdersRepo) OrdersService {
	return &ordersService{repo: repo}
}

func (s *ordersService) ListOrdersByRestaurantID(ctx context.Context, restaurantID uuid.UUID, status string) ([]models.Order, error) {
	return s.repo.ListOrdersByRestaurantID(ctx, restaurantID, status)
}
