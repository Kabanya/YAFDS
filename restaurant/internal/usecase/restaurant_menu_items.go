package usecase

import (
	"restaurant/internal/service"

	"github.com/Kabanya/YAFDS/pkg/models"

	"github.com/google/uuid"
)

type RestaurantMenuItemsUseCase interface {
	ShowMenuItemsByRestaurantID(restaurantID uuid.UUID) ([]models.MenuItem, error)
	UploadMenuItemsByRestaurantID(menuItem models.MenuItem) error
}

type restaurantMenuItemsUseCase struct {
	service service.RestaurantMenuItemsService
}

func NewRestaurantMenuItemsUseCase(service service.RestaurantMenuItemsService) RestaurantMenuItemsUseCase {
	return &restaurantMenuItemsUseCase{
		service: service,
	}
}

func (u *restaurantMenuItemsUseCase) ShowMenuItemsByRestaurantID(restaurantID uuid.UUID) ([]models.MenuItem, error) {
	return u.service.ShowMenuItemsByRestaurantID(restaurantID)
}

func (u *restaurantMenuItemsUseCase) UploadMenuItemsByRestaurantID(menuItem models.MenuItem) error {
	return u.service.UploadMenuItemsByRestaurantID(menuItem)
}
