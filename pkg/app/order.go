package app

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/Kabanya/YAFDS/pkg/common/utils"
	"github.com/Kabanya/YAFDS/pkg/models"
	"github.com/Kabanya/YAFDS/pkg/repository"
	pkgRepoModels "github.com/Kabanya/YAFDS/pkg/repository/models"
	"github.com/Kabanya/YAFDS/pkg/usecase"

	"github.com/google/uuid"
)

// type Filter = repository.Filter
// type Order = repository.Order

func NewOrderHandler(repo pkgRepoModels.OrderRepo, menuClient RestaurantMenuClient) http.HandlerFunc {
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

func NewCreateHandler(repo pkgRepoModels.OrderRepo, menuClient RestaurantMenuClient) http.HandlerFunc {
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

		var req createRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteError(w, "invalid request body", http.StatusBadRequest)
			return
		}
		if menuClient == nil {
			utils.WriteError(w, "menu service unavailable", http.StatusInternalServerError)
			return
		}

		customerID, err := uuid.Parse(req.CustomerID)
		if err != nil {
			utils.WriteError(w, "customer_id must be UUID", http.StatusBadRequest)
			return
		}
		courierID, err := uuid.Parse(req.CourierID)
		if err != nil {
			utils.WriteError(w, "courier_id must be UUID", http.StatusBadRequest)
			return
		}
		restaurantID, err := uuid.Parse(req.RestaurantID)
		if err != nil {
			utils.WriteError(w, "restaurant_id must be UUID", http.StatusBadRequest)
			return
		}
		if len(req.Items) == 0 {
			utils.WriteError(w, "items must not be empty", http.StatusBadRequest)
			return
		}

		menuItems, err := menuClient.GetMenuItems(r.Context(), restaurantID)
		if err != nil {
			logger.Printf("orders: fetch menu items failed: %v", err)
			utils.WriteError(w, "failed to fetch restaurant menu", http.StatusBadGateway)
			return
		}
		menuByID := make(map[uuid.UUID]models.MenuItem, len(menuItems))
		for _, item := range menuItems {
			menuByID[item.OrderItemID] = item
		}

		items := make([]pkgRepoModels.OrderItemInput, 0, len(req.Items))
		for i, item := range req.Items {
			itemID, err := uuid.Parse(item.RestaurantItemID)
			if err != nil {
				utils.WriteError(w, "items["+strconv.Itoa(i)+"].restaurant_item_id must be UUID", http.StatusBadRequest)
				return
			}
			menuItem, ok := menuByID[itemID]
			if !ok {
				utils.WriteError(w, itemNotAvailableError, http.StatusConflict)
				return
			}
			if item.Quantity <= 0 {
				utils.WriteError(w, "items["+strconv.Itoa(i)+"].quantity must be positive", http.StatusBadRequest)
				return
			}
			if menuItem.Quantity <= 0 || item.Quantity > menuItem.Quantity {
				utils.WriteError(w, itemNotAvailableError, http.StatusConflict)
				return
			}
			items = append(items, pkgRepoModels.OrderItemInput{
				RestaurantItemID: menuItem.OrderItemID,
				Price:            menuItem.Price,
				Quantity:         item.Quantity,
			})
		}

		created, err := repo.CreateOrderWithItems(r.Context(), models.Order{
			CustomerID: customerID,
			CourierID:  courierID,
			Status:     models.OrderStatusCustomerCreated,
		}, items)

		if err != nil {
			logger.Printf("orders: create failed: %v", err)
			switch {
			case errors.Is(err, ErrCustomerNotFound):
				utils.WriteError(w, "customer_id not found", http.StatusBadRequest)
			case errors.Is(err, ErrCourierNotFound):
				utils.WriteError(w, "courier_id not found", http.StatusBadRequest)
			default:
				utils.WriteError(w, "failed to create order", http.StatusInternalServerError)
			}
			return
		}

		utils.WriteJSON(w, created, http.StatusCreated)
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
			utils.WriteError(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if menuClient == nil {
			utils.WriteError(w, "menu service unavailable", http.StatusInternalServerError)
			return
		}

		restaurantIDStr := r.URL.Query().Get("restaurant_id")
		if restaurantIDStr == "" {
			utils.WriteError(w, "restaurant_id is required", http.StatusBadRequest)
			return
		}
		restaurantID, err := uuid.Parse(restaurantIDStr)
		if err != nil {
			utils.WriteError(w, "restaurant_id must be UUID", http.StatusBadRequest)
			return
		}

		items, err := menuClient.GetMenuItems(r.Context(), restaurantID)
		if err != nil {
			logger.Printf("menu: fetch restaurant items failed: %v", err)
			utils.WriteError(w, "failed to fetch restaurant menu", http.StatusBadGateway)
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

		utils.WriteJSON(w, response, http.StatusOK)
	}
}

func NewListHandler(repo pkgRepoModels.OrderRepo) http.HandlerFunc {
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
			utils.WriteError(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var filter pkgRepoModels.Filter
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

		orders, err := repo.ListOrders(r.Context(), filter)
		if err != nil {
			logger.Printf("orders: list failed: %v", err)
			utils.WriteError(w, "failed to fetch orders", http.StatusInternalServerError)
			return
		}

		utils.WriteJSON(w, orders, http.StatusOK)
	}
}

func NewAcceptHandler(repo pkgRepoModels.OrderRepo) http.HandlerFunc {
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

		path := strings.TrimPrefix(r.URL.Path, "/orders/")
		path = strings.Trim(path, "/")
		parts := strings.Split(path, "/")
		if len(parts) != 2 || parts[1] != "accept" {
			utils.WriteError(w, "not found", http.StatusNotFound)
			return
		}

		orderID, err := uuid.Parse(parts[0])
		if err != nil {
			utils.WriteError(w, "order_id must be UUID", http.StatusBadRequest)
			return
		}

		var req acceptOrderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteError(w, "invalid request body", http.StatusBadRequest)
			return
		}

		customerID, err := uuid.Parse(req.CustomerID)
		if err != nil {
			utils.WriteError(w, "customer_id must be UUID", http.StatusBadRequest)
			return
		}
		courierID, err := uuid.Parse(req.CourierID)
		if err != nil {
			utils.WriteError(w, "courier_id must be UUID", http.StatusBadRequest)
			return
		}
		if len(req.Items) == 0 {
			utils.WriteError(w, "items must not be empty", http.StatusBadRequest)
			return
		}

		items := make([]pkgRepoModels.OrderItemInput, 0, len(req.Items))
		for i, item := range req.Items {
			restaurantItemID, err := uuid.Parse(item.RestaurantItemID)
			if err != nil {
				utils.WriteError(w, "items["+strconv.Itoa(i)+"].restaurant_item_id must be UUID", http.StatusBadRequest)
				return
			}
			if item.Quantity <= 0 {
				utils.WriteError(w, "items["+strconv.Itoa(i)+"].quantity must be positive", http.StatusBadRequest)
				return
			}
			if item.Price <= 0 {
				utils.WriteError(w, "items["+strconv.Itoa(i)+"].price must be positive", http.StatusBadRequest)
				return
			}
			items = append(items, pkgRepoModels.OrderItemInput{
				RestaurantItemID: restaurantItemID,
				Price:            item.Price,
				Quantity:         item.Quantity,
			})
		}

		accepted, err := repo.AcceptOrder(r.Context(), pkgRepoModels.AcceptInput{
			OrderID:    orderID,
			CustomerID: customerID,
			CourierID:  courierID,
			Items:      items,
		})
		if err != nil {
			logger.Printf("orders: accept failed: %v", err)
			switch {
			case errors.Is(err, ErrCustomerNotFound):
				utils.WriteError(w, "customer_id not found", http.StatusBadRequest)
			case errors.Is(err, ErrCourierNotFound):
				utils.WriteError(w, "courier_id not found", http.StatusBadRequest)
			default:
				utils.WriteError(w, "failed to accept order", http.StatusInternalServerError)
			}
			return
		}

		utils.WriteJSON(w, accepted, http.StatusOK)
	}
}

func NewOrderActionHandler(repo pkgRepoModels.OrderRepo, menuClient RestaurantMenuClient, orderUC usecase.OrderUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger, _ := utils.Logger()
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			w.WriteHeader(http.StatusOK)
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/orders/")
		path = strings.Trim(path, "/")
		parts := strings.Split(path, "/")
		if len(parts) != 2 {
			utils.WriteError(w, "not found", http.StatusNotFound)
			return
		}

		orderID, err := uuid.Parse(parts[0])
		if err != nil {
			utils.WriteError(w, "order_id must be UUID", http.StatusBadRequest)
			return
		}

		switch parts[1] {
		case "pay":
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			if r.Method != http.MethodPost {
				utils.WriteError(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			if orderUC == nil {
				utils.WriteError(w, "order usecase unavailable", http.StatusInternalServerError)
				return
			}

			var req payOrderRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				utils.WriteError(w, "invalid request body", http.StatusBadRequest)
				return
			}

			customerID, err := uuid.Parse(req.CustomerID)
			if err != nil {
				utils.WriteError(w, "customer_id must be UUID", http.StatusBadRequest)
				return
			}

			newStatus, err := orderUC.Pay(r.Context(), orderID, customerID)
			if err != nil {
				logger.Printf("orders: pay failed: %v", err)
				if errors.Is(err, usecase.ErrInsufficientFunds) {
					utils.WriteJSON(w, map[string]string{
						"order_id": orderID.String(),
						"status":   string(newStatus),
						"error":    err.Error(),
					}, http.StatusPaymentRequired)
					return
				}
				utils.WriteError(w, err.Error(), http.StatusInternalServerError)
				return
			}

			utils.WriteJSON(w, map[string]string{
				"order_id": orderID.String(),
				"status":   string(newStatus),
			}, http.StatusOK)

		case "accept":
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
			if r.Method != http.MethodPost {
				utils.WriteError(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			if menuClient == nil {
				utils.WriteError(w, "menu service unavailable", http.StatusInternalServerError)
				return
			}

			var req acceptOrderRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				utils.WriteError(w, "invalid request body", http.StatusBadRequest)
				return
			}

			customerID, err := uuid.Parse(req.CustomerID)
			if err != nil {
				utils.WriteError(w, "customer_id must be UUID", http.StatusBadRequest)
				return
			}
			courierID, err := uuid.Parse(req.CourierID)
			if err != nil {
				utils.WriteError(w, "courier_id must be UUID", http.StatusBadRequest)
				return
			}
			restaurantID, err := uuid.Parse(req.RestaurantID)
			if err != nil {
				utils.WriteError(w, "restaurant_id must be UUID", http.StatusBadRequest)
				return
			}
			if len(req.Items) == 0 {
				utils.WriteError(w, "items must not be empty", http.StatusBadRequest)
				return
			}

			menuItems, err := menuClient.GetMenuItems(r.Context(), restaurantID)
			if err != nil {
				logger.Printf("orders: fetch menu items failed: %v", err)
				utils.WriteError(w, "failed to fetch restaurant menu", http.StatusBadGateway)
				return
			}
			menuByID := make(map[uuid.UUID]models.MenuItem, len(menuItems))
			for _, item := range menuItems {
				menuByID[item.OrderItemID] = item
			}

			items := make([]pkgRepoModels.OrderItemInput, 0, len(req.Items))
			status := models.OrderStatusKitchenAccepted
			for i, item := range req.Items {
				restaurantItemID, err := uuid.Parse(item.RestaurantItemID)
				if err != nil {
					utils.WriteError(w, "items["+strconv.Itoa(i)+"].restaurant_item_id must be UUID", http.StatusBadRequest)
					return
				}
				if item.Quantity <= 0 {
					utils.WriteError(w, "items["+strconv.Itoa(i)+"].quantity must be positive", http.StatusBadRequest)
					return
				}
				if item.Price <= 0 {
					utils.WriteError(w, "items["+strconv.Itoa(i)+"].price must be positive", http.StatusBadRequest)
					return
				}
				menuItem, ok := menuByID[restaurantItemID]
				if !ok || menuItem.Quantity <= 0 || item.Quantity > menuItem.Quantity {
					status = models.OrderStatusKitchenDenied
				}
				items = append(items, pkgRepoModels.OrderItemInput{
					RestaurantItemID: restaurantItemID,
					Price:            item.Price,
					Quantity:         item.Quantity,
				})
			}

			accepted, err := repo.AcceptOrder(r.Context(), pkgRepoModels.AcceptInput{
				OrderID:    orderID,
				CustomerID: customerID,
				CourierID:  courierID,
				Items:      items,
				Status:     status,
			})
			if err != nil {
				logger.Printf("orders: accept failed: %v", err)
				switch {
				case errors.Is(err, ErrCustomerNotFound):
					utils.WriteError(w, "customer_id not found", http.StatusBadRequest)
				case errors.Is(err, ErrCourierNotFound):
					utils.WriteError(w, "courier_id not found", http.StatusBadRequest)
				default:
					utils.WriteError(w, "failed to accept order", http.StatusInternalServerError)
				}
				return
			}

			utils.WriteJSON(w, accepted, http.StatusOK)
		case "items":
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
			if r.Method != http.MethodPost {
				utils.WriteError(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			if menuClient == nil {
				utils.WriteError(w, "menu service unavailable", http.StatusInternalServerError)
				return
			}

			var req addOrderItemRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				utils.WriteError(w, "invalid request body", http.StatusBadRequest)
				return
			}
			restaurantID, err := uuid.Parse(req.RestaurantID)
			if err != nil {
				utils.WriteError(w, "restaurant_id must be UUID", http.StatusBadRequest)
				return
			}
			restaurantItemID, err := uuid.Parse(req.RestaurantItemID)
			if err != nil {
				utils.WriteError(w, "restaurant_item_id must be UUID", http.StatusBadRequest)
				return
			}
			if req.Quantity <= 0 {
				utils.WriteError(w, "quantity must be positive", http.StatusBadRequest)
				return
			}

			menuItems, err := menuClient.GetMenuItems(r.Context(), restaurantID)
			if err != nil {
				logger.Printf("orders: fetch menu items failed: %v", err)
				utils.WriteError(w, "failed to fetch restaurant menu", http.StatusBadGateway)
				return
			}
			menuByID := make(map[uuid.UUID]models.MenuItem, len(menuItems))
			for _, item := range menuItems {
				menuByID[item.OrderItemID] = item
			}
			menuItem, ok := menuByID[restaurantItemID]
			if !ok || menuItem.Quantity <= 0 || req.Quantity > menuItem.Quantity {
				utils.WriteError(w, itemNotAvailableError, http.StatusConflict)
				return
			}

			if err := repo.AddItemIntoOrder(r.Context(), orderID, pkgRepoModels.OrderItemInput{
				RestaurantItemID: menuItem.OrderItemID,
				Price:            menuItem.Price,
				Quantity:         req.Quantity,
			}); err != nil {
				logger.Printf("orders: add item failed: %v", err)
				switch {
				case errors.Is(err, repository.ErrOrderNotFound):
					utils.WriteError(w, "order_id not found", http.StatusNotFound)
				default:
					utils.WriteError(w, "failed to add order item", http.StatusInternalServerError)
				}
				return
			}

			utils.WriteJSON(w, map[string]any{
				"order_id":           orderID,
				"restaurant_item_id": menuItem.OrderItemID,
				"quantity":           req.Quantity,
				"price":              menuItem.Price,
			}, http.StatusCreated)
		default:
			utils.WriteError(w, "not found", http.StatusNotFound)
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
			utils.WriteError(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		rows, err := db.QueryContext(r.Context(), "SELECT emp_id, name FROM COURIERS WHERE is_active = TRUE")
		if err != nil {
			logger.Printf("orders: list couriers failed: %v", err)
			utils.WriteError(w, "failed to fetch couriers", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var couriers []courierResponse
		for rows.Next() {
			var c courierResponse
			if err := rows.Scan(&c.ID, &c.Name); err != nil {
				logger.Printf("orders: scan couriers failed: %v", err)
				utils.WriteError(w, "failed to fetch couriers", http.StatusInternalServerError)
				return
			}
			couriers = append(couriers, c)
		}

		if err := rows.Err(); err != nil {
			logger.Printf("orders: iterate couriers failed: %v", err)
			utils.WriteError(w, "failed to fetch couriers", http.StatusInternalServerError)
			return
		}

		utils.WriteJSON(w, couriers, http.StatusOK)
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
			utils.WriteError(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		rows, err := db.QueryContext(r.Context(), "SELECT emp_id, name FROM RESTAURANTS WHERE status = TRUE")
		if err != nil {
			logger.Printf("orders: list restaurants failed: %v", err)
			utils.WriteError(w, "failed to fetch restaurants", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var restaurants []restaurantResponse
		for rows.Next() {
			var res restaurantResponse
			if err := rows.Scan(&res.ID, &res.Name); err != nil {
				logger.Printf("orders: scan restaurants failed: %v", err)
				utils.WriteError(w, "failed to fetch restaurants", http.StatusInternalServerError)
				return
			}
			restaurants = append(restaurants, res)
		}

		if err := rows.Err(); err != nil {
			logger.Printf("orders: iterate restaurants failed: %v", err)
			utils.WriteError(w, "failed to fetch restaurants", http.StatusInternalServerError)
			return
		}

		utils.WriteJSON(w, restaurants, http.StatusOK)
	}
}
