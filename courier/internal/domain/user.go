package domain

import (
	"errors"

	"github.com/google/uuid"
)

type User struct {
	Id            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	WalletAddress string    `json:"wallet_address"`
	TransportType string    `json:"transport_type"`
	IsActive      bool      `json:"is_active"`
	Geolocation   string    `json:"geolocation"`
	PasswordHash  string    `json:"password_hash,omitempty"`
	PasswordSalt  []byte    `json:"password_salt,omitempty"`
}

type LoginResponse struct {
	Id            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	WalletAddress string    `json:"wallet_address"`
	TransportType string    `json:"transport_type"`
	Token         string    `json:"token"`
	Expiration    int64     `json:"expiration"`
}

var ErrInvalidCredentials = errors.New("invalid credentials")
