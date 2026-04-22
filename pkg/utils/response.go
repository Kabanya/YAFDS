package utils

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

var UuidNil = uuid.Nil

type ErrorResponse struct {
	ErrorMessage string `json:"error_message"`
}

func WriteJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, message string, statusCode int) {
	WriteJSON(w, ErrorResponse{ErrorMessage: message}, statusCode)
}
