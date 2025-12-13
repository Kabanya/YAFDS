package service

import (
	"context"
	"courier/internal/repository"
	"courier/models"
	"crypto/rand"
	"customer/pkg"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// абстракция чуть выше чем репо

type UserService interface {
	Register(uuid.UUID, string, string, string, string) error
	Login(walletAddress string, password string) (models.LoginResponse, error)
}

type userService struct {
	repo       repository.UserRepo
	redis      *redis.Client
	sessionTTL time.Duration
}

func NewUserService(repo repository.UserRepo, redisClient *redis.Client, sessionTTL time.Duration) UserService {
	return &userService{repo: repo, redis: redisClient, sessionTTL: sessionTTL}
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

func (s *userService) Login(walletAddress string, password string) (models.LoginResponse, error) {
	logger, _ := pkg.Logger()

	// Load user from DB
	user, err := s.repo.LoadByWalletAddress(walletAddress)
	if err != nil {
		logger.Printf("Failed to load user: %v", err)
		return models.LoginResponse{}, err
	}

	// Verify password
	if !pkg.VerifyPassword(password, pkg.DefaultParams, user.PasswordSalt, user.PasswordHash) {
		logger.Printf("Invalid password for user: %s", walletAddress)
		return models.LoginResponse{}, models.ErrInvalidCredentials
	}

	token, err := s.createSession(user.Id)
	if err != nil {
		logger.Printf("Failed to create session for user %s: %v", walletAddress, err)
		return models.LoginResponse{}, err
	}

	logger.Printf("User logged in successfully: %s", walletAddress)
	return models.LoginResponse{
		Id:            user.Id,
		Name:          user.Name,
		WalletAddress: user.WalletAddress,
		Address:       user.Address,
		Token:         token,
		Expiration:    time.Now().Add(s.sessionTTL).Unix(),
	}, nil
}

func (s *userService) createSession(userID uuid.UUID) (string, error) {
	logger, _ := pkg.Logger()
	logger.Printf("Creating session for user ID: %s", userID.String())
	ctx := context.Background()
	if s.redis == nil {
		return "", fmt.Errorf("redis client is not initialized")
	}

	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("failed to generate session token: %w", err)
	}

	token := base64.RawURLEncoding.EncodeToString(tokenBytes)
	key := fmt.Sprintf("session:%s", token)

	if err := s.redis.Set(ctx, key, userID.String(), s.sessionTTL).Err(); err != nil {
		return "", fmt.Errorf("failed to store session token: %w", err)
	}
	logger.Printf("Session created for user ID: %s with token: %s", userID.String(), token)
	return token, nil
}
