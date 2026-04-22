package domain

import "github.com/google/uuid"

type CourierCatalogItem struct {
	ID            uuid.UUID
	Name          string
	TransportType string
	IsActive      bool
}

type RestaurantCatalogItem struct {
	ID      uuid.UUID
	Name    string
	Address string
	Status  bool
}

type RestaurantMenuCatalogItem struct {
	ID           uuid.UUID
	RestaurantID uuid.UUID
	Name         string
	Price        float64
	Quantity     int
	Description  string
}
