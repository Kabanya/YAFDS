package utils

import (
	"net/http"

	"github.com/google/uuid"
)

func NewUUID() uuid.UUID {
	return uuid.New()
}

func ParseUUID(idStr string) (uuid.UUID, error) {
	return uuid.Parse(idStr)
}

// ParseUUIDOrBadRequest parses a UUID string and writes a BadRequest error
// to the provided ResponseWriter on failure. Returns the parsed UUID and
// a boolean indicating success.[]
func ParseUUIDOrBadRequest(w http.ResponseWriter, idStr string, errorMsg string) (uuid.UUID, bool) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		WriteError(w, errorMsg, http.StatusBadRequest)
		return uuid.Nil, false
	}
	return id, true
}

// ParseRequiredUUIDOrBadRequest ensures the string is present and a valid UUID,
// writing a BadRequest if it's missing or invalid.
func ParseRequiredUUIDOrBadRequest(w http.ResponseWriter, idStr string, name string) (uuid.UUID, bool) {
	if idStr == "" {
		WriteError(w, name+" is required", http.StatusBadRequest)
		return uuid.Nil, false
	}
	return ParseUUIDOrBadRequest(w, idStr, "invalid "+name)
}
