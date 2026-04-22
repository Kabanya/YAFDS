package usecase

import (
	"context"

	domain "github.com/Kabanya/YAFDS/restaurant/internal/domain"
	"github.com/Kabanya/YAFDS/restaurant/internal/service"
	"github.com/google/uuid"
)

type RestaurantUseCase interface {
	ListRestaurants(ctx context.Context) ([]domain.Restaurant, error)
	GetMenu(ctx context.Context, restaurantID uuid.UUID) ([]domain.RestaurantMenuItem, error)
}

type restaurantUseCase struct {
	svc service.RestaurantService
}

func NewRestaurantUseCase(svc service.RestaurantService) RestaurantUseCase {
	return &restaurantUseCase{svc: svc}
}

func (u *restaurantUseCase) ListRestaurants(ctx context.Context) ([]domain.Restaurant, error) {
	return u.svc.ListRestaurants(ctx)
}

func (u *restaurantUseCase) GetMenu(ctx context.Context, restaurantID uuid.UUID) ([]domain.RestaurantMenuItem, error) {
	return u.svc.GetMenu(ctx, restaurantID)
}
