package app

import (
	"customer/modules/user"
	logger "customer/pkg"
	"errors"
	"sync"
	"time"

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
	GetUserByWalletAddress(walletAddress string) (*user.User, error)
	CreateToken(walletAddress string) (string, int64, error)
	ValidateToken(tokenValue string) (uuid.UUID, error)
}

type service struct {
	mu        sync.RWMutex
	users     map[uuid.UUID]*user.User
	tokens    map[string]*Token
	walletIdx map[string]uuid.UUID
}

func NewService() Service {
	logger.PrintLog("NewService is started!")
	return &service{
		users:     make(map[uuid.UUID]*user.User),
		tokens:    make(map[string]*Token),
		walletIdx: make(map[string]uuid.UUID),
	}
}

func (s *service) Hello() string {
	logger.PrintLog("Start of customer Hello!")
	return "Hello world from customer process"
}

func (s *service) CreateUser(name, walletAddress, address string) (uuid.UUID, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.walletIdx[walletAddress]; exists {
		return uuid.Nil, ErrUserAlreadyExists
	}

	id := uuid.New()
	newUser := &user.User{
		Id:            id,
		Name:          name,
		WalletAddress: walletAddress,
		Address:       address,
	}

	s.users[id] = newUser
	s.walletIdx[walletAddress] = id

	return id, nil
}

func (s *service) GetUserByWalletAddress(walletAddress string) (*user.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	userId, exists := s.walletIdx[walletAddress]
	if !exists {
		return nil, ErrUserNotFound
	}

	return s.users[userId], nil
}

func (s *service) CreateToken(walletAddress string) (string, int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	userId, exists := s.walletIdx[walletAddress]
	if !exists {
		return "", 0, ErrUserNotFound
	}

	tokenValue := uuid.New().String()
	expiration := time.Now().Add(24 * time.Hour).Unix()

	s.tokens[tokenValue] = &Token{
		Value:      tokenValue,
		UserId:     userId,
		Expiration: expiration,
	}

	return tokenValue, expiration, nil
}

func (s *service) ValidateToken(tokenValue string) (uuid.UUID, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	token, exists := s.tokens[tokenValue]
	if !exists {
		return uuid.Nil, ErrTokenNotFound
	}

	if time.Now().Unix() > token.Expiration {
		return uuid.Nil, ErrTokenNotFound
	}

	return token.UserId, nil
}

var _ Service = (*service)(nil)
