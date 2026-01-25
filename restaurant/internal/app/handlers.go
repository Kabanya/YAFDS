package app

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"restaurant/internal/usecase"
	"restaurant/models"

	"github.com/Kabanya/YAFDS/pkg/id"
	pkgmodels "github.com/Kabanya/YAFDS/pkg/models"
	"github.com/Kabanya/YAFDS/pkg/utils"
)

const TransportType = "HTTP"

type Handler struct {
	userUseCase                usecase.UserUseCase
	restaurantMenuItemsUseCase usecase.RestaurantMenuItemsUseCase
	ordersUseCase              usecase.OrdersUseCase
}

func NewHandler(userUC usecase.UserUseCase, menuItemsUC usecase.RestaurantMenuItemsUseCase, ordersUC usecase.OrdersUseCase) *Handler {
	return &Handler{
		userUseCase:                userUC,
		restaurantMenuItemsUseCase: menuItemsUC,
		ordersUseCase:              ordersUC,
	}
}

// Health check endpoint
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, map[string]string{"status": "UP"}, http.StatusOK)
}

// Register user with password
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	logger, _ := utils.Logger()
	logger.Println("Register called")

	// CORS headers
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

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.WalletAddress == "" {
		utils.WriteError(w, "wallet_address is required", http.StatusBadRequest)
		return
	}

	if req.Password == "" {
		utils.WriteError(w, "password is required", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		utils.WriteError(w, "name is required", http.StatusBadRequest)
		return
	}

	if req.Address == "" {
		utils.WriteError(w, "address is required", http.StatusBadRequest)
		return
	}

	// Derive deterministic ID from wallet to keep seeded data stable across runs
	userID := id.FromWallet(req.WalletAddress)

	// Register user with password
	err := h.userUseCase.Register(userID, req.Name, req.WalletAddress, req.Address, req.IsActive, req.Password)
	if err != nil {
		utils.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, models.RegisterResponce{Id: userID}, http.StatusCreated)
	logger.Printf("User %s registered successfully", req.WalletAddress)
}

// Login user with password
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	logger, _ := utils.Logger()
	logger.Println("Login called")

	// CORS headers
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

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.WalletAddress == "" {
		utils.WriteError(w, "wallet_address is required", http.StatusBadRequest)
		return
	}

	if req.Password == "" {
		utils.WriteError(w, "password is required", http.StatusBadRequest)
		return
	}

	// Authenticate user and issue session token
	loginResp, err := h.userUseCase.Login(req.WalletAddress, req.Password)
	if err != nil {
		statusCode := http.StatusInternalServerError
		message := "internal server error"
		if errors.Is(err, models.ErrInvalidCredentials) || errors.Is(err, sql.ErrNoRows) {
			statusCode = http.StatusUnauthorized
			message = "invalid wallet address or password"
		}
		utils.WriteError(w, message, statusCode)
		logger.Printf("Login failed for user: %s, error: %v", req.WalletAddress, err)
		return
	}

	utils.WriteJSON(w, loginResp, http.StatusOK)
	logger.Printf("User %s logged in successfully", req.WalletAddress)
}

// ShowMenuItems returns menu items for a specific restaurant
func (h *Handler) ShowMenuItems(w http.ResponseWriter, r *http.Request) {
	logger, _ := utils.Logger()
	logger.Println("ShowMenuItems called")

	// CORS headers
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

	// Get restaurant_id from query parameter
	restaurantIDStr := r.URL.Query().Get("restaurant_id")
	if restaurantIDStr == "" {
		utils.WriteError(w, "restaurant_id is required", http.StatusBadRequest)
		return
	}

	restaurantID, err := utils.ParseUUID(restaurantIDStr)
	if err != nil {
		utils.WriteError(w, "invalid restaurant_id format", http.StatusBadRequest)
		return
	}

	menuItems, err := h.restaurantMenuItemsUseCase.ShowMenuItemsByRestaurantID(restaurantID)
	if err != nil {
		utils.WriteError(w, err.Error(), http.StatusInternalServerError)
		logger.Printf("Failed to get menu items for restaurant %s: %v", restaurantID, err)
		return
	}

	utils.WriteJSON(w, menuItems, http.StatusOK)
	logger.Printf("Successfully retrieved %d menu items for restaurant %s", len(menuItems), restaurantID)
}

// UploadMenuItem uploads a new menu item for a restaurant
func (h *Handler) UploadMenuItem(w http.ResponseWriter, r *http.Request) {
	logger, _ := utils.Logger()
	logger.Println("UploadMenuItem called")

	// CORS headers
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

	var menuItem pkgmodels.MenuItem
	if err := json.NewDecoder(r.Body).Decode(&menuItem); err != nil {
		utils.WriteError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if menuItem.RestaurantID == utils.UuidNil {
		utils.WriteError(w, "restaurant_id is required", http.StatusBadRequest)
		return
	}

	if menuItem.Name == "" {
		utils.WriteError(w, "name is required", http.StatusBadRequest)
		return
	}

	if menuItem.Price <= 0 {
		utils.WriteError(w, "price must be greater than 0", http.StatusBadRequest)
		return
	}

	// Generate order_item_id if not provided
	if menuItem.OrderItemID == utils.UuidNil {
		menuItem.OrderItemID = utils.NewUUID()
	}

	err := h.restaurantMenuItemsUseCase.UploadMenuItemsByRestaurantID(menuItem)
	if err != nil {
		utils.WriteError(w, err.Error(), http.StatusInternalServerError)
		logger.Printf("Failed to upload menu item: %v", err)
		return
	}

	utils.WriteJSON(w, map[string]interface{}{
		"message":       "menu item uploaded successfully",
		"order_item_id": menuItem.OrderItemID,
	}, http.StatusCreated)
	logger.Printf("Menu item %s uploaded successfully for restaurant %s", menuItem.Name, menuItem.RestaurantID)
}

// ListOrders returns orders for a specific restaurant
func (h *Handler) ListOrders(w http.ResponseWriter, r *http.Request) {
	logger, _ := utils.Logger()
	logger.Println("ListOrders called")

	// CORS headers
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

	restaurantIDStr := r.URL.Query().Get("restaurant_id")
	if restaurantIDStr == "" {
		utils.WriteError(w, "restaurant_id is required", http.StatusBadRequest)
		return
	}

	restaurantID, err := utils.ParseUUID(restaurantIDStr)
	if err != nil {
		utils.WriteError(w, "invalid restaurant_id format", http.StatusBadRequest)
		return
	}

	status := r.URL.Query().Get("status")

	orders, err := h.ordersUseCase.ListOrdersByRestaurantID(r.Context(), restaurantID, status)
	if err != nil {
		utils.WriteError(w, err.Error(), http.StatusInternalServerError)
		logger.Printf("Failed to list orders for restaurant %s: %v", restaurantID, err)
		return
	}

	utils.WriteJSON(w, orders, http.StatusOK)
	logger.Printf("Successfully retrieved %d orders for restaurant %s", len(orders), restaurantID)
}
