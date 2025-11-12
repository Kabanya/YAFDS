package app

import (
	"customer/modules/user"
	logger "customer/pkg"
	"encoding/json"
	"net/http"
)

const TransportType = "HTTP"

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	logger.PrintLog("NewHandler of customer app is started")
	return &Handler{
		service: s,
	}
}

func (h *Handler) Hello(w http.ResponseWriter, r *http.Request) {
	logger.PrintLog("Hello of customer app is started")
	data := h.service.Hello()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req user.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.WalletAddress == "" || req.Address == "" {
		h.writeError(w, "all fields are required", http.StatusBadRequest)
		return
	}

	id, err := h.service.CreateUser(req.Name, req.WalletAddress, req.Address)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.writeJSON(w, user.RegisterResponce{Id: id}, http.StatusCreated)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req user.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.WalletAddress == "" {
		h.writeError(w, "wallet_address is required", http.StatusBadRequest)
		return
	}

	token, expiration, err := h.service.CreateToken(req.WalletAddress)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.writeJSON(w, user.LoginResponse{
		Token:      token,
		Expiration: expiration,
	}, http.StatusOK)
}

func (h *Handler) writeJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) writeError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(user.ErrorResponce{ErrorMessage: message})
}
