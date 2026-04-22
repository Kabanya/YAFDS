package service

import (
	"context"

	"github.com/Kabanya/YAFDS/customer/internal/domain"
	customerrepo "github.com/Kabanya/YAFDS/customer/internal/repository"
	"github.com/google/uuid"
)

type CatalogService interface {
	ListCouriers(ctx context.Context) ([]domain.CourierCatalogItem, error)
	ListRestaurants(ctx context.Context) ([]domain.RestaurantCatalogItem, error)
	ListRestaurantMenu(ctx context.Context, restaurantID uuid.UUID) ([]domain.RestaurantMenuCatalogItem, error)
}

type catalogService struct {
	repo customerrepo.CatalogRepository
}

func NewCatalogService(repo customerrepo.CatalogRepository) CatalogService {
	return &catalogService{repo: repo}
}

func (s *catalogService) ListCouriers(ctx context.Context) ([]domain.CourierCatalogItem, error) {
	couriers, err := s.repo.ListCouriers(ctx)
	if err != nil {
		return nil, err
	}
	return courierCatalogDTOsToDomain(couriers), nil
}

func (s *catalogService) ListRestaurants(ctx context.Context) ([]domain.RestaurantCatalogItem, error) {
	restaurants, err := s.repo.ListRestaurants(ctx)
	if err != nil {
		return nil, err
	}
	return restaurantCatalogDTOsToDomain(restaurants), nil
}

func (s *catalogService) ListRestaurantMenu(ctx context.Context, restaurantID uuid.UUID) ([]domain.RestaurantMenuCatalogItem, error) {
	menu, err := s.repo.ListRestaurantMenu(ctx, restaurantID)
	if err != nil {
		return nil, err
	}
	return restaurantMenuCatalogDTOsToDomain(menu), nil
}

func courierCatalogDTOsToDomain(items []customerrepo.CourierCatalogDTO) []domain.CourierCatalogItem {
	result := make([]domain.CourierCatalogItem, 0, len(items))
	for _, item := range items {
		result = append(result, domain.CourierCatalogItem{
			ID:            item.ID,
			Name:          item.Name,
			TransportType: item.TransportType,
			IsActive:      item.IsActive,
		})
	}
	return result
}

func restaurantCatalogDTOsToDomain(items []customerrepo.RestaurantCatalogDTO) []domain.RestaurantCatalogItem {
	result := make([]domain.RestaurantCatalogItem, 0, len(items))
	for _, item := range items {
		result = append(result, domain.RestaurantCatalogItem{
			ID:      item.ID,
			Name:    item.Name,
			Address: item.Address,
			Status:  item.Status,
		})
	}
	return result
}

func restaurantMenuCatalogDTOsToDomain(items []customerrepo.RestaurantMenuCatalogDTO) []domain.RestaurantMenuCatalogItem {
	result := make([]domain.RestaurantMenuCatalogItem, 0, len(items))
	for _, item := range items {
		result = append(result, domain.RestaurantMenuCatalogItem{
			ID:           item.ID,
			RestaurantID: item.RestaurantID,
			Name:         item.Name,
			Price:        item.Price,
			Quantity:     item.Quantity,
			Description:  item.Description,
		})
	}
	return result
}
