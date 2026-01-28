package app

import (
	"encoding/json"
	"net/http"

	"github.com/Kabanya/YAFDS/pkg/common/utils"
	"github.com/Kabanya/YAFDS/pkg/models"
	repositoryModels "github.com/Kabanya/YAFDS/pkg/repository/models"
	"github.com/Kabanya/YAFDS/pkg/usecase"
	"github.com/google/uuid"
)

// на это уровне могут быть
// клиенты / ручки /

// ПИШЕМ хендлер создания ордера
// всё что он делает -- вызывает usecase
//
// usecase нет напрямую доступа к репозиторию, поэтому он вызывает сервис

// Type aliases from repository
// type Repository = repository.Repository // антипаттерн

// usecase -- сетка зависимостей

// Error aliases from repository

type OrderHandler struct {
	orderUC usecase.OrderUseCase
}

func NewOrderHandler(orderUC usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{
		orderUC: orderUC,
	}
}

// Request/Response structs
type CreateOrderRequest struct {
	CustomerID string `json:"customer_id"`
	CourierID  string `json:"courier_id"`
}

type CreateOrderWithItemsRequest struct {
	CustomerID string                            `json:"customer_id"`
	CourierID  string                            `json:"courier_id"`
	Items      []repositoryModels.OrderItemInput `json:"items"`
}

type GetOrderRequest struct {
	OrderID string `json:"order_id"`
}

type AcceptOrderRequest struct {
	CustomerID string                            `json:"customer_id"`
	CourierID  string                            `json:"courier_id"`
	Items      []repositoryModels.OrderItemInput `json:"items"`
	Status     models.OrderStatus                `json:"status"`
}

type UpdateOrderStatusRequest struct {
	Status models.OrderStatus `json:"status"`
}

type AddItemRequest struct {
	RestaurantItemID string  `json:"restaurant_item_id"`
	Price            float64 `json:"price"`
	Quantity         int     `json:"quantity"`
}

type RemoveItemRequest struct {
	RestaurantItemID string `json:"restaurant_item_id"`
}

// POST /orders
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
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

	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	customerID, err := uuid.Parse(req.CustomerID)
	if err != nil {
		utils.WriteError(w, "invalid customer_id", http.StatusBadRequest)
		return
	}

	courierID, err := uuid.Parse(req.CourierID)
	if err != nil {
		utils.WriteError(w, "invalid courier_id", http.StatusBadRequest)
		return
	}

	order, err := h.orderUC.CreateOrder(r.Context(), customerID, courierID)
	if err != nil {
		utils.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, order, http.StatusCreated)
}

// CreateOrderWithItems creates a new order with items
// POST /orders/with-items
func (h *OrderHandler) CreateOrderWithItems(w http.ResponseWriter, r *http.Request) {
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

	var req CreateOrderWithItemsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	customerID, err := uuid.Parse(req.CustomerID)
	if err != nil {
		utils.WriteError(w, "invalid customer_id", http.StatusBadRequest)
		return
	}

	courierID, err := uuid.Parse(req.CourierID)
	if err != nil {
		utils.WriteError(w, "invalid courier_id", http.StatusBadRequest)
		return
	}

	order, err := h.orderUC.CreateOrderWithItems(r.Context(), customerID, courierID, req.Items)
	if err != nil {
		utils.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, order, http.StatusCreated)
}

// GET /orders/{order_id}
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		utils.WriteError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	orderIDStr := r.URL.Query().Get("order_id")
	if orderIDStr == "" {
		utils.WriteError(w, "order_id is required", http.StatusBadRequest)
		return
	}

	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		utils.WriteError(w, "invalid order_id", http.StatusBadRequest)
		return
	}

	order, err := h.orderUC.GetOrder(r.Context(), orderID)
	if err != nil {
		utils.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, order, http.StatusOK)
}

// POST /orders/{order_id}/accept
func (h *OrderHandler) AcceptOrder(w http.ResponseWriter, r *http.Request) {
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

	orderIDStr := r.URL.Query().Get("order_id")
	if orderIDStr == "" {
		utils.WriteError(w, "order_id is required", http.StatusBadRequest)
		return
	}

	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		utils.WriteError(w, "invalid order_id", http.StatusBadRequest)
		return
	}

	var req AcceptOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	customerID, err := uuid.Parse(req.CustomerID)
	if err != nil {
		utils.WriteError(w, "invalid customer_id", http.StatusBadRequest)
		return
	}

	courierID, err := uuid.Parse(req.CourierID)
	if err != nil {
		utils.WriteError(w, "invalid courier_id", http.StatusBadRequest)
		return
	}

	result, err := h.orderUC.AcceptOrder(r.Context(), orderID, customerID, courierID, req.Items, req.Status)
	if err != nil {
		utils.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, result, http.StatusOK)
}

// GET /orders/{order_id}/status
func (h *OrderHandler) GetOrderStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		utils.WriteError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	orderIDStr := r.URL.Query().Get("order_id")
	if orderIDStr == "" {
		utils.WriteError(w, "order_id is required", http.StatusBadRequest)
		return
	}

	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		utils.WriteError(w, "invalid order_id", http.StatusBadRequest)
		return
	}

	status, err := h.orderUC.GetOrderStatus(r.Context(), orderID)
	if err != nil {
		utils.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, map[string]interface{}{"status": status}, http.StatusOK)
}

// PUT /orders/{order_id}/status
func (h *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPut {
		utils.WriteError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	orderIDStr := r.URL.Query().Get("order_id")
	if orderIDStr == "" {
		utils.WriteError(w, "order_id is required", http.StatusBadRequest)
		return
	}

	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		utils.WriteError(w, "invalid order_id", http.StatusBadRequest)
		return
	}

	var req UpdateOrderStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = h.orderUC.UpdateOrderStatus(r.Context(), orderID, req.Status)
	if err != nil {
		utils.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, map[string]string{"message": "status updated successfully"}, http.StatusOK)
}

// GET /orders/{order_id}/total
func (h *OrderHandler) CalculateOrderTotal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		utils.WriteError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	orderIDStr := r.URL.Query().Get("order_id")
	if orderIDStr == "" {
		utils.WriteError(w, "order_id is required", http.StatusBadRequest)
		return
	}

	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		utils.WriteError(w, "invalid order_id", http.StatusBadRequest)
		return
	}

	total, err := h.orderUC.CalculateOrderTotal(r.Context(), orderID)
	if err != nil {
		utils.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, map[string]interface{}{"total": total}, http.StatusOK)
}

// GET /customers/{customer_id}/wallet
func (h *OrderHandler) GetCustomerWalletAddress(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		utils.WriteError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	customerIDStr := r.URL.Query().Get("customer_id")
	if customerIDStr == "" {
		utils.WriteError(w, "customer_id is required", http.StatusBadRequest)
		return
	}

	customerID, err := uuid.Parse(customerIDStr)
	if err != nil {
		utils.WriteError(w, "invalid customer_id", http.StatusBadRequest)
		return
	}

	walletAddress, err := h.orderUC.GetCustomerWalletAddress(r.Context(), customerID)
	if err != nil {
		utils.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, map[string]string{"wallet_address": walletAddress}, http.StatusOK)
}

// POST /orders/{order_id}/items
func (h *OrderHandler) AddItemIntoOrder(w http.ResponseWriter, r *http.Request) {
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

	orderIDStr := r.URL.Query().Get("order_id")
	if orderIDStr == "" {
		utils.WriteError(w, "order_id is required", http.StatusBadRequest)
		return
	}

	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		utils.WriteError(w, "invalid order_id", http.StatusBadRequest)
		return
	}

	var req AddItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	restaurantItemID, err := uuid.Parse(req.RestaurantItemID)
	if err != nil {
		utils.WriteError(w, "invalid restaurant_item_id", http.StatusBadRequest)
		return
	}

	item := repositoryModels.OrderItemInput{
		RestaurantItemID: restaurantItemID,
		Price:            req.Price,
		Quantity:         req.Quantity,
	}

	err = h.orderUC.AddItemIntoOrder(r.Context(), orderID, item)
	if err != nil {
		utils.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, map[string]string{"message": "item added successfully"}, http.StatusOK)
}

// DELETE /orders/{order_id}/items
func (h *OrderHandler) RemoveItemFromOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodDelete {
		utils.WriteError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	orderIDStr := r.URL.Query().Get("order_id")
	if orderIDStr == "" {
		utils.WriteError(w, "order_id is required", http.StatusBadRequest)
		return
	}

	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		utils.WriteError(w, "invalid order_id", http.StatusBadRequest)
		return
	}

	restaurantItemIDStr := r.URL.Query().Get("restaurant_item_id")
	if restaurantItemIDStr == "" {
		utils.WriteError(w, "restaurant_item_id is required", http.StatusBadRequest)
		return
	}

	restaurantItemID, err := uuid.Parse(restaurantItemIDStr)
	if err != nil {
		utils.WriteError(w, "invalid restaurant_item_id", http.StatusBadRequest)
		return
	}

	err = h.orderUC.RemoveItemFromOrder(r.Context(), orderID, restaurantItemID)
	if err != nil {
		utils.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, map[string]string{"message": "item removed successfully"}, http.StatusOK)
}
