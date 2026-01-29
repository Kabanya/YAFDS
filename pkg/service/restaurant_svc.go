package service

import (
	"context"
	"github.com/Kabanya/YAFDS/pkg/app/clients"
	"github.com/Kabanya/YAFDS/pkg/models"
	repositoryModels "github.com/Kabanya/YAFDS/pkg/repository/models"
	"github.com/google/uuid"
)

type RestaurantService interface {
	ListRestaurants(ctx context.Context) ([]models.Restaurant, error)
	GetMenu(ctx context.Context, restaurantID uuid.UUID) ([]models.RestaurantMenuItem, error)
}

type restaurantService struct {
	repo   repositoryModels.RestaurantRepo
	client clients.RestaurantClient
}

func NewRestaurantService(repo repositoryModels.RestaurantRepo, client clients.RestaurantClient) RestaurantService {
	return &restaurantService{repo: repo, client: client}
}

func (s *restaurantService) ListRestaurants(ctx context.Context) ([]models.Restaurant, error) {
	return s.repo.ListRestaurants(ctx)
}

func (s *restaurantService) GetMenu(ctx context.Context, restaurantID uuid.UUID) ([]models.RestaurantMenuItem, error) {
	// We need to map clients.RestaurantMenuItem to models.RestaurantMenuItem
	clientMenu, err := s.client.GetMenu(ctx, restaurantID)
	if err != nil {
		return nil, err
	}

	menu := make([]models.RestaurantMenuItem, len(clientMenu))
	for i, item := range clientMenu {
		menu[i] = models.RestaurantMenuItem{
			ID:           item.ID,
			RestaurantID: item.RestaurantID,
			Name:         item.Name,
			Price:        item.Price,
		}
	}
	return menu, nil
}
