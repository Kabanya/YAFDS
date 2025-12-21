package app

import (
	"customer/pkg"
	"customer/pkg/id"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"restaurant/internal/usecase"
	"restaurant/models"
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

func (h *Handler) writeJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) writeError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.ErrorResponce{ErrorMessage: message})
}

// Register user with password
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	logger, _ := pkg.Logger()
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
		h.writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.WalletAddress == "" {
		h.writeError(w, "wallet_address is required", http.StatusBadRequest)
		return
	}

	if req.Password == "" {
		h.writeError(w, "password is required", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		h.writeError(w, "name is required", http.StatusBadRequest)
		return
	}

	if req.Address == "" {
		h.writeError(w, "address is required", http.StatusBadRequest)
		return
	}

	// Derive deterministic ID from wallet to keep seeded data stable across runs
	userID := id.FromWallet(req.WalletAddress)

	// Register user with password
	err := h.userUseCase.Register(userID, req.Name, req.WalletAddress, req.Address, req.IsActive, req.Password)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, models.RegisterResponce{Id: userID}, http.StatusCreated)
	logger.Printf("User %s registered successfully", req.WalletAddress)
}

// Login user with password
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	logger, _ := pkg.Logger()
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
		h.writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.WalletAddress == "" {
		h.writeError(w, "wallet_address is required", http.StatusBadRequest)
		return
	}

	if req.Password == "" {
		h.writeError(w, "password is required", http.StatusBadRequest)
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
		h.writeError(w, message, statusCode)
		logger.Printf("Login failed for user: %s, error: %v", req.WalletAddress, err)
		return
	}

	h.writeJSON(w, loginResp, http.StatusOK)
	logger.Printf("User %s logged in successfully", req.WalletAddress)
}
