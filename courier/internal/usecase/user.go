package usecase

import (
	"courier/internal/service"
	"courier/models"

	"github.com/google/uuid"
)

// не очень умно копируем

type UserUseCase interface {
	Register(id uuid.UUID, name string, walletAddress string, transportType string, password string) error
	Login(walletAddress string, password string) (models.LoginResponse, error)
}

type userUseCase struct {
	service service.UserService
}

func NewUserUseCase(service service.UserService) UserUseCase {
	return &userUseCase{service: service}
}

func (u *userUseCase) Register(id uuid.UUID, name string, walletAddress string, transportType string, password string) error {
	return u.service.Register(id, name, walletAddress, transportType, password)
}

func (u *userUseCase) Login(walletAddress string, password string) (models.LoginResponse, error) {
	return u.service.Login(walletAddress, password)
}
