// // принимает интерфейсы , возвращаем указатели на структуры сервиса
// // dependency infection
// // не содержит объявление интерфейсов
package service

import (
	"customer/models"
	"errors"

	// "os/user" // Remove this import, as you are using your own models.User
	"sync"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrTokenNotFound     = errors.New("token not found")
)

type Token struct {
	Value      string
	UserId     uuid.UUID
	Expiration int64
}

type Service interface {
	Hello() string
	CreateUser(name, walletAddress, address string) (uuid.UUID, error)
	GetUserByWalletAddress(walletAddress string) (*models.User, error)
	CreateToken(walletAddress string) (string, int64, error)
	ValidateToken(tokenValue string) (uuid.UUID, error)
}

type service struct {
	mu        sync.RWMutex
	users     map[uuid.UUID]*models.User
	tokens    map[string]*Token
	walletIdx map[string]uuid.UUID
}

// принимает интерфейсы , возвращаем указатели на структуры сервиса
func NewService() /*Service - не тру, потому что интерфейс*/ *service {
	return &service{
		users:     make(map[uuid.UUID]*models.User),
		tokens:    make(map[string]*Token),
		walletIdx: make(map[string]uuid.UUID),
	}
}

func (s *service) Hello() string {
	return "Hello world from customer process"
}

func (s *service) CreateUser(name, walletAddress, address string) (uuid.UUID, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.walletIdx[walletAddress]; exists {
		return uuid.Nil, ErrUserAlreadyExists
	}

	id := uuid.New()
	newUser := &models.User{
		Id:            id,
		Name:          name,
		WalletAddress: walletAddress,
		Address:       address,
	}

	s.users[id] = newUser
	s.walletIdx[walletAddress] = id

	return id, nil
}

func (s *service) GetUserByWalletAddress(walletAddress string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	userId, exists := s.walletIdx[walletAddress]
	if !exists {
		return nil, ErrUserNotFound
	}

	return s.users[userId], nil
}
