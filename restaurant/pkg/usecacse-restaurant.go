package usecase

import (
	"context"
	"github.com/Kabanya/YAFDS/pkg/models"
	"github.com/Kabanya/YAFDS/pkg/service"
	"github.com/google/uuid"
)

type RestaurantUseCase interface {
	ListRestaurants(ctx context.Context) ([]models.Restaurant, error)
	GetMenu(ctx context.Context, restaurantID uuid.UUID) ([]models.RestaurantMenuItem, error)
}

type restaurantUseCase struct {
	svc service.RestaurantService
}

func NewRestaurantUseCase(svc service.RestaurantService) RestaurantUseCase {
	return &restaurantUseCase{svc: svc}
}

func (u *restaurantUseCase) ListRestaurants(ctx context.Context) ([]models.Restaurant, error) {
	return u.svc.ListRestaurants(ctx)
}

func (u *restaurantUseCase) GetMenu(ctx context.Context, restaurantID uuid.UUID) ([]models.RestaurantMenuItem, error) {
	return u.svc.GetMenu(ctx, restaurantID)
}
