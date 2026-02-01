package usecase

import (
	"customer/internal/service"
	"customer/models"

	"github.com/google/uuid"
)

// не очень умно копируем

type UserUseCase interface {
	Register(uuid.UUID, string, string, string, string) error
	Login(walletAddress string, password string) (models.LoginResponse, error)
}

type userUseCase struct {
	service service.UserService
}

func NewUserUseCase(service service.UserService) UserUseCase {
	return &userUseCase{service: service}
}

func (u *userUseCase) Register(id uuid.UUID, name string, walletAddress string, address string, password string) error {
	return u.service.Register(id, name, walletAddress, address, password)
}

func (u *userUseCase) Login(walletAddress string, password string) (models.LoginResponse, error) {
	return u.service.Login(walletAddress, password)
}
