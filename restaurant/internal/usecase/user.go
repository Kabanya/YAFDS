package usecase

import (
	"restaurant/internal/service"
	"restaurant/models"

	"github.com/google/uuid"
)

// не очень умно копируем

type UserUseCase interface {
	Register(id uuid.UUID, name string, walletAddress string, address string, isActive bool, password string) error
	Login(walletAddress string, password string) (models.LoginResponse, error)
}

type userUseCase struct {
	service service.UserService
}

func NewUserUseCase(service service.UserService) UserUseCase {
	return &userUseCase{service: service}
}

func (u *userUseCase) Register(id uuid.UUID, name string, walletAddress string, address string, isActive bool, password string) error {
	return u.service.Register(id, name, walletAddress, address, isActive, password)
}

func (u *userUseCase) Login(walletAddress string, password string) (models.LoginResponse, error) {
	return u.service.Login(walletAddress, password)
}
