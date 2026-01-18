package service

import (
	"net/http"

	"customer/pkg/repository"
	"customer/pkg/utils"

	"github.com/google/uuid"
)

type Repository = repository.Repository
type Filter = repository.Filter

type createOrderRequest struct {
	CustomerID string `json:"customer_id"`
	CourierID  string `json:"courier_id"`
	Status     string `json:"status"`
}

type acceptOrderRequest struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

type createOrderResponce struct {
	OrderID string `json:"order_id"`
}

type createOrderErrorResponse struct {
	Error string `json:"error"`
}

func NewAcceptHandler(repo Repository) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		logger, _ := utils.Logger()
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method != http.MethodPost {
			utils.WriteError(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var filter Filter
		if v := r.URL.Query().Get("customer_id"); v != "" {
			id, err := uuid.Parse(v)
			if err != nil {
				utils.WriteError(w, "problem with customer_id (UUID)", http.StatusBadRequest)
				return
			}
			filter.CustomerID = &id
		}
		if v := r.URL.Query().Get("courier_id"); v != "" {
			id, err := uuid.Parse(v)
			if err != nil {
				utils.WriteError(w, "problem with courier_id (UUID)", http.StatusBadRequest)
				return
			}
			filter.CourierID = &id
		}
		if v := r.URL.Query().Get("status"); v != "" {
			filter.Status = v
		}

		orders, err := repo.List(r.Context(), filter)
		if err != nil {
			logger.Printf("orders: list failed: %v", err)
			utils.WriteError(w, "failed to fetch orders", http.StatusInternalServerError)
			return
		}

		utils.WriteJSON(w, orders, http.StatusOK)
	}
}
