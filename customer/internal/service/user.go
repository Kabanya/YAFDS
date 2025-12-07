package service

import (
	"customer/internal/repository"
	"customer/models"

	"github.com/google/uuid"
)

// абстракция чуть выше чем репо

type UserService interface {
	Save(uuid.UUID, string, string, string) error
	Load(walletAddress string) (models.User, error)
}

type userService struct {
	repo repository.UserRepo
}

func NewUserService(repo repository.UserRepo) UserService {
	return &userService{repo: repo}
}

func (s *userService) Save(id uuid.UUID, name string, walletAddress string, address string) error {
	return s.repo.Save(id, name, walletAddress, address)
}

func (s *userService) Load(walletAddress string) (models.User, error) {
	return s.repo.Load(walletAddress)
}
