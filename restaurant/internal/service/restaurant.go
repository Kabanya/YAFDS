package service

import (
	"context"

	models "github.com/Kabanya/YAFDS/restaurant/internal/domain"
	"github.com/google/uuid"
)

type RestaurantClient interface {
	GetMenu(ctx context.Context, restaurantID uuid.UUID) ([]models.RestaurantMenuItem, error)
}

type RestaurantService interface {
	ListRestaurants(ctx context.Context) ([]models.Restaurant, error)
	GetMenu(ctx context.Context, restaurantID uuid.UUID) ([]models.RestaurantMenuItem, error)
}

type restaurantService struct {
	repo   models.RestaurantRepo
	client RestaurantClient
}

func NewRestaurantService(repo models.RestaurantRepo, client RestaurantClient) RestaurantService {
	return &restaurantService{repo: repo, client: client}
}

func (s *restaurantService) ListRestaurants(ctx context.Context) ([]models.Restaurant, error) {
	return s.repo.ListRestaurants(ctx)
}

func (s *restaurantService) GetMenu(ctx context.Context, restaurantID uuid.UUID) ([]models.RestaurantMenuItem, error) {
	if s.client == nil {
		return []models.RestaurantMenuItem{}, nil
	}
	return s.client.GetMenu(ctx, restaurantID)
}
