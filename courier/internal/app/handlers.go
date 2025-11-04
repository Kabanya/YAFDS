package app

import (
	"encoding/json"
	"net/http"
)

const TransportType = "HTTP"

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) Hello(w http.ResponseWriter, r *http.Request) {
	data := h.service.Hello()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
