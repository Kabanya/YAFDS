package usecase

import (
	"context"

	"github.com/Kabanya/YAFDS/customer/internal/domain"
	"github.com/Kabanya/YAFDS/customer/internal/service"
	"github.com/google/uuid"
)

type CatalogUseCase interface {
	ListCouriers(ctx context.Context) ([]domain.CourierCatalogItem, error)
	ListRestaurants(ctx context.Context) ([]domain.RestaurantCatalogItem, error)
	ListRestaurantMenu(ctx context.Context, restaurantID uuid.UUID) ([]domain.RestaurantMenuCatalogItem, error)
}

type catalogUseCase struct {
	service service.CatalogService
}

func NewCatalogUseCase(service service.CatalogService) CatalogUseCase {
	return &catalogUseCase{service: service}
}

func (u *catalogUseCase) ListCouriers(ctx context.Context) ([]domain.CourierCatalogItem, error) {
	return u.service.ListCouriers(ctx)
}

func (u *catalogUseCase) ListRestaurants(ctx context.Context) ([]domain.RestaurantCatalogItem, error) {
	return u.service.ListRestaurants(ctx)
}

func (u *catalogUseCase) ListRestaurantMenu(ctx context.Context, restaurantID uuid.UUID) ([]domain.RestaurantMenuCatalogItem, error) {
	return u.service.ListRestaurantMenu(ctx, restaurantID)
}
