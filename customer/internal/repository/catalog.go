package repository

import (
	"context"

	"github.com/google/uuid"
)

type CourierCatalogDTO struct {
	ID            uuid.UUID
	Name          string
	TransportType string
	IsActive      bool
}

type RestaurantCatalogDTO struct {
	ID      uuid.UUID
	Name    string
	Address string
	Status  bool
}

type RestaurantMenuCatalogDTO struct {
	ID           uuid.UUID
	RestaurantID uuid.UUID
	Name         string
	Price        float64
	Quantity     int
	Description  string
}

type CatalogRepository interface {
	ListCouriers(ctx context.Context) ([]CourierCatalogDTO, error)
	ListRestaurants(ctx context.Context) ([]RestaurantCatalogDTO, error)
	ListRestaurantMenu(ctx context.Context, restaurantID uuid.UUID) ([]RestaurantMenuCatalogDTO, error)
}
