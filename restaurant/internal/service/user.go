package service

import (
	"context"
	"errors"
	"restaurant/internal/repository"
	"restaurant/models"
	"time"

	"github.com/Kabanya/YAFDS/pkg/auth"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// абстракция чуть выше чем репо

type UserService interface {
	Register(id uuid.UUID, name string, walletAddress string, address string, isActive bool, password string) error
	Login(walletAddress string, password string) (models.LoginResponse, error)
}

type userService struct {
	authService *auth.Service
}

func NewUserService(repo repository.UserRepo, redisClient *redis.Client, sessionTTL time.Duration) UserService {
	service, err := auth.NewService(auth.ServiceConfig{
		Store:      storeAdapter{repo: repo},
		Hasher:     auth.NewArgon2Hasher(auth.DefaultArgonParams),
		Sessions:   auth.NewRedisSessionManager(redisClient),
		Validator:  auth.NoopValidator,
		SessionTTL: sessionTTL,
	})
	if err != nil {
		panic(err)
	}
	return &userService{authService: service}
}

func (s *userService) Register(id uuid.UUID, name string, walletAddress string, address string, isActive bool, password string) error {
	return s.authService.Register(context.Background(), auth.RegisterInput{
		ID:            id,
		Name:          name,
		WalletAddress: walletAddress,
		Address:       address,
		IsActive:      isActive,
		Password:      password,
	})
}

func (s *userService) Login(walletAddress string, password string) (models.LoginResponse, error) {
	res, err := s.authService.Login(context.Background(), walletAddress, password)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return models.LoginResponse{}, models.ErrInvalidCredentials
		}
		return models.LoginResponse{}, err
	}
	return models.LoginResponse{
		Id:            res.User.ID,
		Name:          res.User.Name,
		WalletAddress: res.User.WalletAddress,
		Address:       res.User.Address,
		IsActive:      res.User.IsActive,
		Token:         res.Token,
		Expiration:    res.Expiration.Unix(),
	}, nil
}

type storeAdapter struct {
	repo repository.UserRepo
}

func (a storeAdapter) SaveWithPassword(ctx context.Context, data auth.RegisterInput, passwordHash string, passwordSalt []byte) error {
	return a.repo.SaveWithPassword(data.ID, data.Name, data.WalletAddress, data.Address, data.IsActive, passwordHash, passwordSalt)
}

func (a storeAdapter) LoadByWalletAddress(ctx context.Context, walletAddress string) (auth.StoredUser, error) {
	user, err := a.repo.LoadByWalletAddress(walletAddress)
	if err != nil {
		return auth.StoredUser{}, err
	}
	return auth.StoredUser{
		ID:            user.Id,
		Name:          user.Name,
		WalletAddress: user.WalletAddress,
		Address:       user.Address,
		IsActive:      user.IsActive,
		PasswordHash:  user.PasswordHash,
		PasswordSalt:  user.PasswordSalt,
	}, nil
}
