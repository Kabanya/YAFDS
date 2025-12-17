package models

import (
	"errors"

	"github.com/google/uuid"
)

type RegisterRequest struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	WalletAddress string `json:"wallet_address"`
	Address       string `json:"address"`
	IsActive      bool   `json:"is_active"`
	Password      string `json:"password"`
}

type User struct { //моделька
	Id            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	WalletAddress string    `json:"wallet_address"`
	Address       string    `json:"address"`
	IsActive      bool      `json:"is_active"`
	PasswordHash  string    `json:"password_hash,omitempty"`
	PasswordSalt  []byte    `json:"password_salt,omitempty"`
}

type ErrorResponce struct {
	ErrorMessage string `json:"error_message"`
}

type RegisterResponce struct {
	Id uuid.UUID `json:"id"`
}

type LoginRequest struct {
	WalletAddress string `json:"wallet_address"`
	Password      string `json:"password"`
}

type LoginResponse struct {
	Id            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	WalletAddress string    `json:"wallet_address"`
	Address       string    `json:"address"`
	IsActive      bool      `json:"is_active"`
	Token         string    `json:"token"`
	Expiration    int64     `json:"expiration"`
}

var ErrInvalidCredentials = errors.New("invalid credentials")

func NewError(message string) error {
	return errors.New(message)
}
