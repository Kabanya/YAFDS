package app

import (
	"courier/internal/usecase"
	"courier/models"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Kabanya/YAFDS/pkg/common/id"
	"github.com/Kabanya/YAFDS/pkg/common/utils"
)

const TransportType = "HTTP"

type Handler struct {
	userUseCase usecase.UserUseCase
}

func NewHandler(userUC usecase.UserUseCase) *Handler {
	return &Handler{
		userUseCase: userUC,
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

	transportType := req.TransportType
	if transportType == "" {
		transportType = "bicycle" // default
	}

	// Derive deterministic ID from wallet to keep seeded data stable across runs
	userID := id.FromWallet(req.WalletAddress)

	// After stress testing, need to add limit to registrations with same data

	// Register user with password
	err := h.userUseCase.Register(userID, req.Name, req.WalletAddress, transportType, req.Password)
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
