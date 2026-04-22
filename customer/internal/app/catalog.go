package app

import (
	"net/http"

	"github.com/Kabanya/YAFDS/customer/internal/domain"
	"github.com/Kabanya/YAFDS/customer/internal/usecase"
	"github.com/Kabanya/YAFDS/pkg/utils"
	"github.com/google/uuid"
)

type CourierResponse struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	TransportType string    `json:"transport_type"`
	IsActive      bool      `json:"is_active"`
}

type RestaurantResponse struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Address string    `json:"address"`
	Status  bool      `json:"status"`
}

type RestaurantMenuItemResponse struct {
	ID           uuid.UUID `json:"id"`
	RestaurantID uuid.UUID `json:"restaurant_id"`
	Name         string    `json:"name"`
	Price        float64   `json:"price"`
	Quantity     int       `json:"quantity"`
	Description  string    `json:"description"`
}

func NewCourierCatalogHandler(catalogUC usecase.CatalogUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.WriteError(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		couriers, err := catalogUC.ListCouriers(r.Context())
		if err != nil {
			utils.WriteError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		utils.WriteJSON(w, courierResponsesFromDomain(couriers), http.StatusOK)
	}
}

func NewRestaurantCatalogHandler(catalogUC usecase.CatalogUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.WriteError(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		restaurants, err := catalogUC.ListRestaurants(r.Context())
		if err != nil {
			utils.WriteError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		utils.WriteJSON(w, restaurantResponsesFromDomain(restaurants), http.StatusOK)
	}
}

func NewRestaurantMenuCatalogHandler(catalogUC usecase.CatalogUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.WriteError(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		restaurantID, ok := utils.ParseRequiredUUIDOrBadRequest(w, r.URL.Query().Get("restaurant_id"), "restaurant_id")
		if !ok {
			return
		}

		menu, err := catalogUC.ListRestaurantMenu(r.Context(), restaurantID)
		if err != nil {
			utils.WriteError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		utils.WriteJSON(w, restaurantMenuResponsesFromDomain(menu), http.StatusOK)
	}
}

func courierResponsesFromDomain(items []domain.CourierCatalogItem) []CourierResponse {
	result := make([]CourierResponse, 0, len(items))
	for _, item := range items {
		result = append(result, CourierResponse{
			ID:            item.ID,
			Name:          item.Name,
			TransportType: item.TransportType,
			IsActive:      item.IsActive,
		})
	}
	return result
}

func restaurantResponsesFromDomain(items []domain.RestaurantCatalogItem) []RestaurantResponse {
	result := make([]RestaurantResponse, 0, len(items))
	for _, item := range items {
		result = append(result, RestaurantResponse{
			ID:      item.ID,
			Name:    item.Name,
			Address: item.Address,
			Status:  item.Status,
		})
	}
	return result
}

func restaurantMenuResponsesFromDomain(items []domain.RestaurantMenuCatalogItem) []RestaurantMenuItemResponse {
	result := make([]RestaurantMenuItemResponse, 0, len(items))
	for _, item := range items {
		result = append(result, RestaurantMenuItemResponse{
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
