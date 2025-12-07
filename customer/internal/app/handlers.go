package app

import (
	"customer/internal/usecase"
	"customer/models"
	"customer/pkg"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
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

// via usecase
func (h *Handler) SaveUser(w http.ResponseWriter, r *http.Request) {
	logger, _ := pkg.Logger()
	logger.Println("SaveUser called")
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

	if req.Id == "" || req.Name == "" || req.WalletAddress == "" || req.Address == "" {
		h.writeError(w, "all fields are required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(req.Id)
	if err != nil {
		h.writeError(w, "invalid id format", http.StatusBadRequest)
		return
	}

	err = h.userUseCase.Save(id, req.Name, req.WalletAddress, req.Address)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, models.RegisterResponce{Id: id}, http.StatusCreated)
	logger.Printf("User %s saved successfully", req.WalletAddress)
}

// via usecase
func (h *Handler) LoadUser(w http.ResponseWriter, r *http.Request) {
	logger, _ := pkg.Logger()
	logger.Println("LoadUser called")
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

	user, err := h.userUseCase.Load(req.WalletAddress)
	if err != nil {
		h.writeError(w, "user not found", http.StatusNotFound)
		return
	}

	h.writeJSON(w, user, http.StatusOK)
	logger.Printf("User %s loaded successfully", req.WalletAddress)
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
