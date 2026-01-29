package utils

import "github.com/google/uuid"

func NewUUID() uuid.UUID {
	return uuid.New()
}

func ParseUUID(idStr string) (uuid.UUID, error) {
	return uuid.Parse(idStr)
}
