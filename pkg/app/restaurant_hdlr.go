package app

import (
	"net/http"

	"github.com/Kabanya/YAFDS/pkg/common/utils"
	"github.com/Kabanya/YAFDS/pkg/usecase"
	"github.com/google/uuid"
)

func NewRestaurantsHandler(resUC usecase.RestaurantUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		restaurants, err := resUC.ListRestaurants(r.Context())
		if err != nil {
			utils.WriteError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		utils.WriteJSON(w, restaurants, http.StatusOK)
	}
}

func NewRestaurantMenuHandler(resUC usecase.RestaurantUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		restaurantIDStr := r.URL.Query().Get("restaurant_id")
		if restaurantIDStr == "" {
			utils.WriteError(w, "restaurant_id is required", http.StatusBadRequest)
			return
		}

		restaurantID, err := uuid.Parse(restaurantIDStr)
		if err != nil {
			utils.WriteError(w, "invalid restaurant_id", http.StatusBadRequest)
			return
		}

		menu, err := resUC.GetMenu(r.Context(), restaurantID)
		if err != nil {
			utils.WriteError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		utils.WriteJSON(w, menu, http.StatusOK)
	}
}
