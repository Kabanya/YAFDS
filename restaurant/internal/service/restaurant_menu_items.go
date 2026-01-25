package service

import (
	"restaurant/internal/repository"

	"github.com/Kabanya/YAFDS/pkg/models"

	"github.com/google/uuid"
)

type RestaurantMenuItemsService interface {
	ShowMenuItemsByRestaurantID(restaurantID uuid.UUID) ([]models.MenuItem, error)
	UploadMenuItemsByRestaurantID(menuItem models.MenuItem) error
}

type restaurantMenuItemsService struct {
	repo repository.RestaurantMenuItemsRepo
}

func NewRestaurantMenuItemsService(repo repository.RestaurantMenuItemsRepo) RestaurantMenuItemsService {
	return &restaurantMenuItemsService{
		repo: repo,
	}
}

func (s *restaurantMenuItemsService) ShowMenuItemsByRestaurantID(restaurantID uuid.UUID) ([]models.MenuItem, error) {
	return s.repo.ShowMenuItemsByRestaurantID(restaurantID)
}

func (s *restaurantMenuItemsService) UploadMenuItemsByRestaurantID(menuItem models.MenuItem) error {
	return s.repo.UploadMenuItemsByRestaurantID(menuItem)
}
