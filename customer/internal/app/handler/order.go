package app

import (
	"encoding/json"
	"net/http"

	"github.com/Kabanya/YAFDS/pkg/app/client"
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

func NewOrderHandler(orderUC usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{
		orderUC: orderUC,
	}
}

func (h *OrderHandler) parseOrderID(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	orderIDStr := r.PathValue("order_id")
	if orderIDStr == "" {
		orderIDStr = r.URL.Query().Get("order_id")
	}
	if orderIDStr == "" {
		utils.WriteError(w, "order_id is required", http.StatusBadRequest)
		return uuid.Nil, false
	}

	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		utils.WriteError(w, "invalid order_id", http.StatusBadRequest)
		return uuid.Nil, false
	}
	return orderID, true
}

func (h *OrderHandler) OrdersHandler(repo repositoryModels.OrderRepo) http.HandlerFunc {
	listHandler := NewListHandler(repo)
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			listHandler(w, r)
		case http.MethodPost:
			h.CreateOrder(w, r)
		default:
			utils.WriteError(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func NewListHandler(repo repositoryModels.OrderRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		status := r.URL.Query().Get("status")
		customerIDStr := r.URL.Query().Get("customer_id")
		courierIDStr := r.URL.Query().Get("courier_id")

		var filter repositoryModels.Filter
		if status != "" {
			filter.Status = status
		}
		if customerIDStr != "" {
			if id, err := uuid.Parse(customerIDStr); err == nil {
				filter.CustomerID = &id
			}
		}
		if courierIDStr != "" {
			if id, err := uuid.Parse(courierIDStr); err == nil {
				filter.CourierID = &id
			}
		}

		orders, err := repo.ListOrders(r.Context(), filter)
		if err != nil {
			utils.WriteError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if orders == nil {
			orders = []models.Order{}
		}

		utils.WriteJSON(w, orders, http.StatusOK)
	}
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

	customerID, ok := utils.ParseUUIDOrBadRequest(w, req.CustomerID, "invalid customer_id")
	if !ok {
		return
	}

	courierID, ok := utils.ParseUUIDOrBadRequest(w, req.CourierID, "invalid courier_id")
	if !ok {
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

	customerID, ok := utils.ParseUUIDOrBadRequest(w, req.CustomerID, "invalid customer_id")
	if !ok {
		return
	}

	courierID, ok := utils.ParseUUIDOrBadRequest(w, req.CourierID, "invalid courier_id")
	if !ok {
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

	orderID, ok := h.parseOrderID(w, r)
	if !ok {
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

	orderID, ok := h.parseOrderID(w, r)
	if !ok {
		return
	}

	var req AcceptOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	customerID, ok := utils.ParseUUIDOrBadRequest(w, req.CustomerID, "invalid customer_id")
	if !ok {
		return
	}

	courierID, ok := utils.ParseUUIDOrBadRequest(w, req.CourierID, "invalid courier_id")
	if !ok {
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

	orderID, ok := h.parseOrderID(w, r)
	if !ok {
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

	orderID, ok := h.parseOrderID(w, r)
	if !ok {
		return
	}

	var req UpdateOrderStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err := h.orderUC.UpdateOrderStatus(r.Context(), orderID, req.Status)
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

	orderID, ok := h.parseOrderID(w, r)
	if !ok {
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

	customerID, ok := utils.ParseRequiredUUIDOrBadRequest(w, r.URL.Query().Get("customer_id"), "customer_id")
	if !ok {
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

	orderID, ok := h.parseOrderID(w, r)
	if !ok {
		return
	}

	var req AddItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	restaurantItemID, ok := utils.ParseUUIDOrBadRequest(w, req.RestaurantItemID, "invalid restaurant_item_id")
	if !ok {
		return
	}

	item := repositoryModels.OrderItemInput{
		RestaurantItemID: restaurantItemID,
		Price:            req.Price,
		Quantity:         req.Quantity,
	}

	err := h.orderUC.AddItemIntoOrder(r.Context(), orderID, item)
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

	orderID, ok := h.parseOrderID(w, r)
	if !ok {
		return
	}

	restaurantItemID, ok := utils.ParseRequiredUUIDOrBadRequest(w, r.URL.Query().Get("restaurant_item_id"), "restaurant_item_id")
	if !ok {
		return
	}

	err := h.orderUC.RemoveItemFromOrder(r.Context(), orderID, restaurantItemID)
	if err != nil {
		utils.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, map[string]string{"message": "item removed successfully"}, http.StatusOK)
}

// POST /orders/{order_id}/pay
func (h *OrderHandler) PayOrder(w http.ResponseWriter, r *http.Request, walletClient client.WalletClient) {
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

	orderID, ok := h.parseOrderID(w, r)
	if !ok {
		return
	}

	err := h.orderUC.PayOrder(r.Context(), orderID, walletClient)
	if err != nil {
		utils.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, map[string]string{"message": "order paid successfully"}, http.StatusOK)
}
