package models

import (
	"context"
	pkgModels "github.com/Kabanya/YAFDS/pkg/models"
)

type RestaurantRepo interface {
	ListRestaurants(ctx context.Context) ([]pkgModels.Restaurant, error)
}
