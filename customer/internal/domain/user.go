package domain

import "github.com/google/uuid"

type RegisterRequest struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	WalletAddress string `json:"wallet_address"`
	Address       string `json:"address"`
	IsActive      bool   `json:"is_active"`
	Password      string `json:"password"`
}

type RegisterResponse struct {
	ID uuid.UUID `json:"id"`
}

type LoginRequest struct {
	WalletAddress string `json:"wallet_address"`
	Password      string `json:"password"`
}

type LoginResponse struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	WalletAddress string    `json:"wallet_address"`
	Address       string    `json:"address"`
	IsActive      bool      `json:"is_active"`
	Token         string    `json:"token"`
	Expiration    int64     `json:"expiration"`
}

type User struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	WalletAddress string    `json:"wallet_address"`
	Address       string    `json:"address"`
	IsActive      bool      `json:"is_active"`
	PasswordHash  string    `json:"password_hash,omitempty"`
	PasswordSalt  []byte    `json:"password_salt,omitempty"`
}
