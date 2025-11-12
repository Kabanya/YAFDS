package user

import (
	"github.com/google/uuid"
)

type RegisterRequest struct {
	Name          string `json:"name"`
	WalletAddress string `json:"wallet_address"`
	Address       string `json:"address"`
}

type User struct {
	Id            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	WalletAddress string    `json:"wallet_address"`
	Address       string    `json:"address"`
}

type ErrorResponce struct {
	ErrorMessage string `json:"error_message"`
}

type RegisterResponce struct {
	Id uuid.UUID `json:"id"`
}

type LoginRequest struct {
	WalletAddress string `json:"wallet_address"`
}

type LoginResponse struct {
	Token      string `json:"token"`
	Expiration int64  `json:"expiration"`
}
