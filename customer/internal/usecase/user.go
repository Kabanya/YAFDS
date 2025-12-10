package usecase

import (
	"customer/internal/service"
	"customer/models"

	"github.com/google/uuid"
)

// не очень умно копируем

type UserUseCase interface {
	Save(uuid.UUID, string, string, string) error
	Load(walletAddress string) (models.User, error)
}

type userUseCase struct {
	service service.UserService
}

func NewUserUseCase(service service.UserService) UserUseCase {
	return &userUseCase{service: service}
}

func (u *userUseCase) Save(id uuid.UUID, name string, walletAddress string, address string) error {
	return u.service.Save(id, name, walletAddress, address)
}

func (u *userUseCase) Load(walletAddress string) (models.User, error) {
	return u.service.Load(walletAddress)
}
