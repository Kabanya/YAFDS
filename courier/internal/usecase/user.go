package usecase

import (
	"github.com/Kabanya/YAFDS/courier/internal/domain"
	"github.com/Kabanya/YAFDS/courier/internal/service"

	"github.com/google/uuid"
)

// не очень умно копируем

type UserUseCase interface {
	Register(id uuid.UUID, name string, walletAddress string, transportType string, password string) error
	Login(walletAddress string, password string) (domain.LoginResponse, error)
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

func (u *userUseCase) Login(walletAddress string, password string) (domain.LoginResponse, error) {
	return u.service.Login(walletAddress, password)
}
