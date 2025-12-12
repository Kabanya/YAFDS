package service

import (
	"customer/internal/repository"
	"customer/models"
	"customer/pkg"

	"github.com/google/uuid"
)

// абстракция чуть выше чем репо

type UserService interface {
	Register(uuid.UUID, string, string, string, string) error
	Login(walletAddress string, password string) (models.User, error)
}

type userService struct {
	repo repository.UserRepo
}

func NewUserService(repo repository.UserRepo) UserService {
	return &userService{repo: repo}
}

func (s *userService) Register(id uuid.UUID, name string, walletAddress string, address string, password string) error {
	logger, _ := pkg.Logger()

	if password == "" {
		return models.NewError("password is required")
	}

	// Hash password
	passwordHash, salt, err := pkg.HashPassword(password, pkg.ArgonParams{})
	if err != nil {
		logger.Printf("Failed to hash password: %v", err)
		return err
	}

	return s.repo.SaveWithPassword(id, name, walletAddress, address, passwordHash, salt)
}

func (s *userService) Login(walletAddress string, password string) (models.User, error) {
	logger, _ := pkg.Logger()

	// Load user from DB
	user, err := s.repo.LoadByWalletAddress(walletAddress)
	if err != nil {
		logger.Printf("Failed to load user: %v", err)
		return models.User{}, err
	}

	// Verify password
	if !pkg.VerifyPassword(password, pkg.ArgonParams{}, user.PasswordSalt, user.PasswordHash) {
		logger.Printf("Invalid password for user: %s", walletAddress)
		return models.User{}, models.NewError("invalid password")
	}

	logger.Printf("User logged in successfully: %s", walletAddress)
	// Never return password hash/salt to callers.
	user.PasswordHash = ""
	user.PasswordSalt = nil
	return user, nil
}
