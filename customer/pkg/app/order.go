package app

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"customer/pkg/clients"
	"customer/pkg/repository"
	"customer/pkg/utils"

	"github.com/google/uuid"
)

// Type aliases from repository
type Repository = repository.Repository
type Filter = repository.Filter
type Order = repository.Order

// Error aliases from repository
var (
	ErrCustomerNotFound = repository.ErrCustomerNotFound
	ErrCourierNotFound  = repository.ErrCourierNotFound
)

type errorResponse struct {
	Error string `json:"error"`
}

type courierResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type restaurantResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type createRequest struct {
	CustomerID   string                   `json:"customer_id"`
	CourierID    string                   `json:"courier_id"`
	RestaurantID string                   `json:"restaurant_id"`
	Status       string                   `json:"status"`
	Items        []createOrderItemRequest `json:"items"`
}

type createOrderItemRequest struct {
	RestaurantItemID string `json:"restaurant_item_id"`
	Quantity         int    `json:"quantity"`
}

type acceptOrderItemRequest struct {
	RestaurantItemID string  `json:"restaurant_item_id"`
	Price            float64 `json:"price"`
	Quantity         int     `json:"quantity"`
}

type acceptOrderRequest struct {
	CustomerID string                   `json:"customer_id"`
	CourierID  string                   `json:"courier_id"`
	Items      []acceptOrderItemRequest `json:"items"`
}

type addOrderItemRequest struct {
	RestaurantID     string `json:"restaurant_id"`
	RestaurantItemID string `json:"restaurant_item_id"`
	Quantity         int    `json:"quantity"`
}

type menuItemResponse struct {
	OrderItemID  uuid.UUID `json:"order_item_id"`
	RestaurantID uuid.UUID `json:"restaurant_id"`
	Name         string    `json:"name"`
	Price        float64   `json:"price"`
	Description  string    `json:"description"`
}

type RestaurantMenuClient interface {
	GetMenuItems(ctx context.Context, restaurantID uuid.UUID) ([]clients.RestaurantMenuItem, error)
}

const itemNotAvailableError = "ITEM_NOT_AVAILABLE"

func NewHandler(repo Repository, menuClient RestaurantMenuClient) http.HandlerFunc {
	create := NewCreateHandler(repo, menuClient)
	list := NewListHandler(repo)

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			create(w, r)
		case http.MethodGet:
			list(w, r)
		case http.MethodOptions:
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.WriteHeader(http.StatusOK)
		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			_ = json.NewEncoder(w).Encode(errorResponse{Error: "method not allowed"})
		}
	}
}

func NewCreateHandler(repo Repository, menuClient RestaurantMenuClient) http.HandlerFunc {
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
			writeError(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req createRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, "invalid request body", http.StatusBadRequest)
			return
		}
		if menuClient == nil {
			writeError(w, "menu service unavailable", http.StatusInternalServerError)
			return
		}

		customerID, err := uuid.Parse(req.CustomerID)
		if err != nil {
			writeError(w, "customer_id must be UUID", http.StatusBadRequest)
			return
		}
		courierID, err := uuid.Parse(req.CourierID)
		if err != nil {
			writeError(w, "courier_id must be UUID", http.StatusBadRequest)
			return
		}
		restaurantID, err := uuid.Parse(req.RestaurantID)
		if err != nil {
			writeError(w, "restaurant_id must be UUID", http.StatusBadRequest)
			return
		}
		if len(req.Items) == 0 {
			writeError(w, "items must not be empty", http.StatusBadRequest)
			return
		}

		menuItems, err := menuClient.GetMenuItems(r.Context(), restaurantID)
		if err != nil {
			logger.Printf("orders: fetch menu items failed: %v", err)
			writeError(w, "failed to fetch restaurant menu", http.StatusBadGateway)
			return
		}
		menuByID := make(map[uuid.UUID]clients.RestaurantMenuItem, len(menuItems))
		for _, item := range menuItems {
			menuByID[item.OrderItemID] = item
		}

		items := make([]repository.OrderItemInput, 0, len(req.Items))
		for i, item := range req.Items {
			itemID, err := uuid.Parse(item.RestaurantItemID)
			if err != nil {
				writeError(w, "items["+strconv.Itoa(i)+"].restaurant_item_id must be UUID", http.StatusBadRequest)
				return
			}
			menuItem, ok := menuByID[itemID]
			if !ok {
				writeError(w, itemNotAvailableError, http.StatusConflict)
				return
			}
			if item.Quantity <= 0 {
				writeError(w, "items["+strconv.Itoa(i)+"].quantity must be positive", http.StatusBadRequest)
				return
			}
			if menuItem.Quantity <= 0 || item.Quantity > menuItem.Quantity {
				writeError(w, itemNotAvailableError, http.StatusConflict)
				return
			}
			items = append(items, repository.OrderItemInput{
				RestaurantItemID: menuItem.OrderItemID,
				Price:            menuItem.Price,
				Quantity:         item.Quantity,
			})
		}

		created, err := repo.CreateWithItems(r.Context(), Order{
			CustomerID: customerID,
			CourierID:  courierID,
			Status:     req.Status,
		}, items)
		if err != nil {
			logger.Printf("orders: create failed: %v", err)
			switch {
			case errors.Is(err, ErrCustomerNotFound):
				writeError(w, "customer_id not found", http.StatusBadRequest)
			case errors.Is(err, ErrCourierNotFound):
				writeError(w, "courier_id not found", http.StatusBadRequest)
			default:
				writeError(w, "failed to create order", http.StatusInternalServerError)
			}
			return
		}

		writeJSON(w, created, http.StatusCreated)
	}
}

func NewRestaurantMenuHandler(menuClient RestaurantMenuClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger, _ := utils.Logger()
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method != http.MethodGet {
			writeError(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if menuClient == nil {
			writeError(w, "menu service unavailable", http.StatusInternalServerError)
			return
		}

		restaurantIDStr := r.URL.Query().Get("restaurant_id")
		if restaurantIDStr == "" {
			writeError(w, "restaurant_id is required", http.StatusBadRequest)
			return
		}
		restaurantID, err := uuid.Parse(restaurantIDStr)
		if err != nil {
			writeError(w, "restaurant_id must be UUID", http.StatusBadRequest)
			return
		}

		items, err := menuClient.GetMenuItems(r.Context(), restaurantID)
		if err != nil {
			logger.Printf("menu: fetch restaurant items failed: %v", err)
			writeError(w, "failed to fetch restaurant menu", http.StatusBadGateway)
			return
		}

		response := make([]menuItemResponse, 0, len(items))
		for _, item := range items {
			response = append(response, menuItemResponse{
				OrderItemID:  item.OrderItemID,
				RestaurantID: item.RestaurantID,
				Name:         item.Name,
				Price:        item.Price,
				Description:  item.Description,
			})
		}

		writeJSON(w, response, http.StatusOK)
	}
}

func NewListHandler(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger, _ := utils.Logger()
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method != http.MethodGet {
			writeError(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var filter Filter
		if v := r.URL.Query().Get("customer_id"); v != "" {
			id, err := uuid.Parse(v)
			if err != nil {
				writeError(w, "problem with customer_id (UUID)", http.StatusBadRequest)
				return
			}
			filter.CustomerID = &id
		}
		if v := r.URL.Query().Get("courier_id"); v != "" {
			id, err := uuid.Parse(v)
			if err != nil {
				writeError(w, "problem with courier_id (UUID)", http.StatusBadRequest)
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
			writeError(w, "failed to fetch orders", http.StatusInternalServerError)
			return
		}

		writeJSON(w, orders, http.StatusOK)
	}
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
			writeError(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/orders/")
		path = strings.Trim(path, "/")
		parts := strings.Split(path, "/")
		if len(parts) != 2 || parts[1] != "accept" {
			writeError(w, "not found", http.StatusNotFound)
			return
		}

		orderID, err := uuid.Parse(parts[0])
		if err != nil {
			writeError(w, "order_id must be UUID", http.StatusBadRequest)
			return
		}

		var req acceptOrderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, "invalid request body", http.StatusBadRequest)
			return
		}

		customerID, err := uuid.Parse(req.CustomerID)
		if err != nil {
			writeError(w, "customer_id must be UUID", http.StatusBadRequest)
			return
		}
		courierID, err := uuid.Parse(req.CourierID)
		if err != nil {
			writeError(w, "courier_id must be UUID", http.StatusBadRequest)
			return
		}
		if len(req.Items) == 0 {
			writeError(w, "items must not be empty", http.StatusBadRequest)
			return
		}

		items := make([]repository.OrderItemInput, 0, len(req.Items))
		for i, item := range req.Items {
			restaurantItemID, err := uuid.Parse(item.RestaurantItemID)
			if err != nil {
				writeError(w, "items["+strconv.Itoa(i)+"].restaurant_item_id must be UUID", http.StatusBadRequest)
				return
			}
			if item.Quantity <= 0 {
				writeError(w, "items["+strconv.Itoa(i)+"].quantity must be positive", http.StatusBadRequest)
				return
			}
			if item.Price <= 0 {
				writeError(w, "items["+strconv.Itoa(i)+"].price must be positive", http.StatusBadRequest)
				return
			}
			items = append(items, repository.OrderItemInput{
				RestaurantItemID: restaurantItemID,
				Price:            item.Price,
				Quantity:         item.Quantity,
			})
		}

		accepted, err := repo.Accept(r.Context(), repository.AcceptInput{
			OrderID:    orderID,
			CustomerID: customerID,
			CourierID:  courierID,
			Items:      items,
		})
		if err != nil {
			logger.Printf("orders: accept failed: %v", err)
			switch {
			case errors.Is(err, ErrCustomerNotFound):
				writeError(w, "customer_id not found", http.StatusBadRequest)
			case errors.Is(err, ErrCourierNotFound):
				writeError(w, "courier_id not found", http.StatusBadRequest)
			default:
				writeError(w, "failed to accept order", http.StatusInternalServerError)
			}
			return
		}

		writeJSON(w, accepted, http.StatusOK)
	}
}

func NewOrderActionHandler(repo Repository, menuClient RestaurantMenuClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger, _ := utils.Logger()
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		path := strings.TrimPrefix(r.URL.Path, "/orders/")
		path = strings.Trim(path, "/")
		parts := strings.Split(path, "/")
		if len(parts) != 2 {
			writeError(w, "not found", http.StatusNotFound)
			return
		}

		orderID, err := uuid.Parse(parts[0])
		if err != nil {
			writeError(w, "order_id must be UUID", http.StatusBadRequest)
			return
		}

		switch parts[1] {
		case "accept":
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
			if r.Method != http.MethodPost {
				writeError(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}

			var req acceptOrderRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				writeError(w, "invalid request body", http.StatusBadRequest)
				return
			}

			customerID, err := uuid.Parse(req.CustomerID)
			if err != nil {
				writeError(w, "customer_id must be UUID", http.StatusBadRequest)
				return
			}
			courierID, err := uuid.Parse(req.CourierID)
			if err != nil {
				writeError(w, "courier_id must be UUID", http.StatusBadRequest)
				return
			}
			if len(req.Items) == 0 {
				writeError(w, "items must not be empty", http.StatusBadRequest)
				return
			}

			items := make([]repository.OrderItemInput, 0, len(req.Items))
			for i, item := range req.Items {
				restaurantItemID, err := uuid.Parse(item.RestaurantItemID)
				if err != nil {
					writeError(w, "items["+strconv.Itoa(i)+"].restaurant_item_id must be UUID", http.StatusBadRequest)
					return
				}
				if item.Quantity <= 0 {
					writeError(w, "items["+strconv.Itoa(i)+"].quantity must be positive", http.StatusBadRequest)
					return
				}
				if item.Price <= 0 {
					writeError(w, "items["+strconv.Itoa(i)+"].price must be positive", http.StatusBadRequest)
					return
				}
				items = append(items, repository.OrderItemInput{
					RestaurantItemID: restaurantItemID,
					Price:            item.Price,
					Quantity:         item.Quantity,
				})
			}

			accepted, err := repo.Accept(r.Context(), repository.AcceptInput{
				OrderID:    orderID,
				CustomerID: customerID,
				CourierID:  courierID,
				Items:      items,
			})
			if err != nil {
				logger.Printf("orders: accept failed: %v", err)
				switch {
				case errors.Is(err, ErrCustomerNotFound):
					writeError(w, "customer_id not found", http.StatusBadRequest)
				case errors.Is(err, ErrCourierNotFound):
					writeError(w, "courier_id not found", http.StatusBadRequest)
				default:
					writeError(w, "failed to accept order", http.StatusInternalServerError)
				}
				return
			}

			writeJSON(w, accepted, http.StatusOK)
		case "items":
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
			if r.Method != http.MethodPost {
				writeError(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			if menuClient == nil {
				writeError(w, "menu service unavailable", http.StatusInternalServerError)
				return
			}

			var req addOrderItemRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				writeError(w, "invalid request body", http.StatusBadRequest)
				return
			}
			restaurantID, err := uuid.Parse(req.RestaurantID)
			if err != nil {
				writeError(w, "restaurant_id must be UUID", http.StatusBadRequest)
				return
			}
			restaurantItemID, err := uuid.Parse(req.RestaurantItemID)
			if err != nil {
				writeError(w, "restaurant_item_id must be UUID", http.StatusBadRequest)
				return
			}
			if req.Quantity <= 0 {
				writeError(w, "quantity must be positive", http.StatusBadRequest)
				return
			}

			menuItems, err := menuClient.GetMenuItems(r.Context(), restaurantID)
			if err != nil {
				logger.Printf("orders: fetch menu items failed: %v", err)
				writeError(w, "failed to fetch restaurant menu", http.StatusBadGateway)
				return
			}
			menuByID := make(map[uuid.UUID]clients.RestaurantMenuItem, len(menuItems))
			for _, item := range menuItems {
				menuByID[item.OrderItemID] = item
			}
			menuItem, ok := menuByID[restaurantItemID]
			if !ok || menuItem.Quantity <= 0 || req.Quantity > menuItem.Quantity {
				writeError(w, itemNotAvailableError, http.StatusConflict)
				return
			}

			if err := repo.AddItem(r.Context(), orderID, repository.OrderItemInput{
				RestaurantItemID: menuItem.OrderItemID,
				Price:            menuItem.Price,
				Quantity:         req.Quantity,
			}); err != nil {
				logger.Printf("orders: add item failed: %v", err)
				switch {
				case errors.Is(err, repository.ErrOrderNotFound):
					writeError(w, "order_id not found", http.StatusNotFound)
				default:
					writeError(w, "failed to add order item", http.StatusInternalServerError)
				}
				return
			}

			writeJSON(w, map[string]any{
				"order_id":           orderID,
				"restaurant_item_id": menuItem.OrderItemID,
				"quantity":           req.Quantity,
				"price":              menuItem.Price,
			}, http.StatusCreated)
		default:
			writeError(w, "not found", http.StatusNotFound)
		}
	}
}

func NewCouriersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger, _ := utils.Logger()

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method != http.MethodGet {
			writeError(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		rows, err := db.QueryContext(r.Context(), "SELECT emp_id, name FROM COURIERS WHERE is_active = TRUE")
		if err != nil {
			logger.Printf("orders: list couriers failed: %v", err)
			writeError(w, "failed to fetch couriers", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var couriers []courierResponse
		for rows.Next() {
			var c courierResponse
			if err := rows.Scan(&c.ID, &c.Name); err != nil {
				logger.Printf("orders: scan couriers failed: %v", err)
				writeError(w, "failed to fetch couriers", http.StatusInternalServerError)
				return
			}
			couriers = append(couriers, c)
		}

		if err := rows.Err(); err != nil {
			logger.Printf("orders: iterate couriers failed: %v", err)
			writeError(w, "failed to fetch couriers", http.StatusInternalServerError)
			return
		}

		writeJSON(w, couriers, http.StatusOK)
	}
}

func NewRestaurantsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger, _ := utils.Logger()

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method != http.MethodGet {
			writeError(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		rows, err := db.QueryContext(r.Context(), "SELECT emp_id, name FROM RESTAURANTS WHERE status = TRUE")
		if err != nil {
			logger.Printf("orders: list restaurants failed: %v", err)
			writeError(w, "failed to fetch restaurants", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var restaurants []restaurantResponse
		for rows.Next() {
			var res restaurantResponse
			if err := rows.Scan(&res.ID, &res.Name); err != nil {
				logger.Printf("orders: scan restaurants failed: %v", err)
				writeError(w, "failed to fetch restaurants", http.StatusInternalServerError)
				return
			}
			restaurants = append(restaurants, res)
		}

		if err := rows.Err(); err != nil {
			logger.Printf("orders: iterate restaurants failed: %v", err)
			writeError(w, "failed to fetch restaurants", http.StatusInternalServerError)
			return
		}

		writeJSON(w, restaurants, http.StatusOK)
	}
}

func writeJSON(w http.ResponseWriter, data any, status int) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, message string, status int) {
	writeJSON(w, errorResponse{Error: message}, status)
}
