package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type Hasher interface {
	Hash(password string) (hash string, salt []byte, err error)
	Verify(password string, salt []byte, expected string) bool
}

type SessionManager interface {
	Create(ctx context.Context, userID uuid.UUID, ttl time.Duration) (token string, expiration time.Time, err error)
}

type Validator func(ctx context.Context, data RegisterInput) error

type Store interface {
	SaveWithPassword(ctx context.Context, data RegisterInput, passwordHash string, passwordSalt []byte) error
	LoadByWalletAddress(ctx context.Context, walletAddress string) (StoredUser, error)
}

type RegisterInput struct {
	ID            uuid.UUID
	Name          string
	WalletAddress string
	Address       string
	Password      string
}

type StoredUser struct {
	ID            uuid.UUID
	Name          string
	WalletAddress string
	Address       string
	PasswordHash  string
	PasswordSalt  []byte
}

type LoginResult struct {
	User       StoredUser
	Token      string
	Expiration time.Time
}

type ServiceConfig struct {
	Store      Store
	Hasher     Hasher
	Sessions   SessionManager
	Validator  Validator
	SessionTTL time.Duration
}

var NoopValidator Validator = func(context.Context, RegisterInput) error { return nil }
