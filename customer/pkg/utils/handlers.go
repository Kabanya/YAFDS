package utils

import (
	"customer/models"
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, message string, statusCode int) {
	WriteJSON(w, models.ErrorResponce{ErrorMessage: message}, statusCode)
}
